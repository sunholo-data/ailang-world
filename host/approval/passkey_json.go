package approval

// Row 201 A2 (design §2.7): the response files the ceremony page saves: the subset of L3's
// PublicKeyCredential.toJSON() the laws read. Unknown members are ignored; no envelope value is
// trusted unchecked (rawId is checked against authData or selects the key; the byte fields are
// covered by the enrolment code or the signature).

import (
	"encoding/base64"
	"encoding/json"
	"errors"
)

// MaxPasskeyResponseBytes bounds a response file (E1, checked by the CLI's reader).
const MaxPasskeyResponseBytes = 16 << 10

var ErrPasskeyResponse = errors.New("approval: not a passkey response file (id, rawId, type public-key, response)")

// RegistrationResponse is a parsed registration file.
type RegistrationResponse struct{ RawID, ClientDataJSON, AttestationObject []byte }

// AssertionResponse is a parsed assertion file. UserHandle is nil when absent or null.
type AssertionResponse struct {
	RawID, ClientDataJSON, AuthenticatorData, Signature, UserHandle []byte
	HasUserHandle                                                   bool
}

type passkeyEnvelope struct {
	ID, RawID, Type *string
	Response        *struct {
		ClientDataJSON, AttestationObject, AuthenticatorData, Signature, UserHandle *string
	}
}

// decodeEnvelope parses the envelope (E2, E3) and returns rawId's bytes.
func decodeEnvelope(b []byte) (passkeyEnvelope, []byte, error) {
	var e passkeyEnvelope
	if err := json.Unmarshal(b, &e); err != nil || e.ID == nil || e.RawID == nil || e.Type == nil || e.Response == nil { // E2
		return passkeyEnvelope{}, nil, ErrPasskeyResponse
	}
	if *e.Type != "public-key" || *e.ID != *e.RawID { // E3
		return passkeyEnvelope{}, nil, ErrPasskeyResponse
	}
	raw, err := b64Field(e.RawID)
	return e, raw, err
}

// b64Field decodes a required member as strict, unpadded base64url (E2 when missing, E4).
func b64Field(s *string) ([]byte, error) {
	if s == nil { // E2
		return nil, ErrPasskeyResponse
	}
	b, err := base64.RawURLEncoding.Strict().DecodeString(*s)
	if err != nil { // E4
		return nil, ErrPasskeyResponse
	}
	return b, nil
}

// ParseRegistrationResponse reads a registration file.
func ParseRegistrationResponse(b []byte) (RegistrationResponse, error) {
	e, raw, err := decodeEnvelope(b)
	if err != nil {
		return RegistrationResponse{}, err
	}
	r := RegistrationResponse{RawID: raw}
	if r.ClientDataJSON, err = b64Field(e.Response.ClientDataJSON); err != nil {
		return RegistrationResponse{}, err
	}
	if r.AttestationObject, err = b64Field(e.Response.AttestationObject); err != nil {
		return RegistrationResponse{}, err
	}
	return r, nil
}

// ParseAssertionResponse reads an assertion file.
func ParseAssertionResponse(b []byte) (AssertionResponse, error) {
	e, raw, err := decodeEnvelope(b)
	if err != nil {
		return AssertionResponse{}, err
	}
	r := AssertionResponse{RawID: raw}
	for _, f := range []struct {
		src *string
		dst *[]byte
	}{{e.Response.ClientDataJSON, &r.ClientDataJSON}, {e.Response.AuthenticatorData, &r.AuthenticatorData}, {e.Response.Signature, &r.Signature}} {
		if *f.dst, err = b64Field(f.src); err != nil {
			return AssertionResponse{}, err
		}
	}
	if e.Response.UserHandle != nil {
		if r.UserHandle, err = b64Field(e.Response.UserHandle); err != nil {
			return AssertionResponse{}, err
		}
		r.HasUserHandle = true
	}
	return r, nil
}

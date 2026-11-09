package approval

// fixtureMessageHex is SignedMessage(fixtureDigest, fixtureChallenge(), fixtureExpiry):
// the 90 bytes every SSHSIG fixture signs. testdata/sshsig/generate.sh reads this line.
const fixtureMessageHex = "776f726c642d617070726f76616c2f7631000ace927ef952f541b1e454b066f6d6b1941fb5f3d710c784c22dd80a53d771eba0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf000000006955b900"

// fixtureExpiry is 2026-01-01T00:00:00Z.
const fixtureExpiry int64 = 1767225600

func fixtureChallenge() [32]byte {
	var c [32]byte
	for i := range c {
		c[i] = byte(0xa0 + i)
	}
	return c
}

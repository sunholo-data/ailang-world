package daemon

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const (
	episodeFilesSchema          = "world/surface-episode-files/v1"
	surfaceEpisodeFileIntents   = 128 // per invocation namespace
	surfaceEpisodeFileLookups   = 1024
	surfaceEpisodeFilePaths     = 128 // per activity class
	surfaceEpisodeFilePathBytes = 512
)

type surfaceEpisodeFiles struct {
	Schema      string                 `json:"schema"`
	Episode     string                 `json:"episode"`
	SourceState string                 `json:"sourceState"`
	Status      string                 `json:"status"`
	Reads       []surfaceFileUse       `json:"reads"`
	Writes      []surfaceFileUse       `json:"writes"`
	Produced    []surfaceFileUse       `json:"produced"`
	Omitted     *surfaceFileUseOmitted `json:"omitted,omitempty"`
	DeltaGaps   int                    `json:"deltaGaps,omitempty"`
}

type surfaceFileUse struct {
	Path       string `json:"path"`
	Op         string `json:"op,omitempty"`
	Count      int    `json:"count"`
	FirstEntry int64  `json:"firstEntry"`
	lastEntry  int64
}

type surfaceFileUseOmitted struct {
	Reads    int `json:"reads"`
	Writes   int `json:"writes"`
	Produced int `json:"produced"`
	Unlisted int `json:"unlisted"`
}

var (
	errSurfaceEpisodeAddress = errors.New("expected a brane and state qualified by the same repository")
	errEpisodeFileBudget     = errors.New("episode file evidence read limit")
)

func selectSurfaceEpisode(env surface.FollowIndex, brane, at string, now time.Time) (surfacepage.FollowStack, time.Time, error) {
	repo, kind, _, idErr := surface.ParseQualifiedID(at)
	braneRepo, braneKind, _, braneErr := surface.ParseQualifiedID(brane)
	if idErr != nil || braneErr != nil || repo != braneRepo || braneKind != "brane" || (kind != "commit" && kind != "episode") {
		return surfacepage.FollowStack{}, time.Time{}, errSurfaceEpisodeAddress
	}
	st, err := surfacepage.BuildFollowStack(env, brane, at, now)
	if err != nil {
		return surfacepage.FollowStack{}, time.Time{}, err
	}
	st = trajectorySelection(env, st)
	if st.Front.Kind == surface.KindStep {
		through, err := time.Parse(time.RFC3339, st.Front.At)
		return st, through, err
	}
	return st, time.Time{}, nil
}

func (d *Daemon) handleSurfaceEpisodeFiles(w http.ResponseWriter, r *http.Request) {
	refuse := func(status int, reason string) {
		d.writeSurfaceJSON(w, r, status, episodeFilesSchema, surfaceReason{Schema: episodeFilesSchema, Reason: reason})
	}
	if len(r.URL.RawQuery) > 2048 {
		refuse(http.StatusBadRequest, "expected one brane and state")
		return
	}
	q, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(q) != 2 || len(q["brane"]) != 1 || len(q["at"]) != 1 {
		refuse(http.StatusBadRequest, "expected one brane and state")
		return
	}
	if d.follow == nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	if d.follow.fileReads.Add(1) > 2 {
		d.follow.fileReads.Add(-1)
		refuse(http.StatusServiceUnavailable, "episode file activity is busy")
		return
	}
	defer d.follow.fileReads.Add(-1)
	env := d.followSnapshot().Envelope
	selection, through, err := selectSurfaceEpisode(env, q.Get("brane"), q.Get("at"), time.Now())
	if errors.Is(err, errSurfaceEpisodeAddress) {
		refuse(http.StatusBadRequest, errSurfaceEpisodeAddress.Error())
		return
	}
	if err != nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	start := d.readFollowEpisodeStart(ctx, env, selection)
	if start == nil && ctx.Err() != nil {
		refuse(http.StatusServiceUnavailable, "episode file activity could not be read")
		return
	}
	if start == nil {
		refuse(http.StatusNotFound, "no verified episode for this state")
		return
	}
	files := buildEpisodeFileActivity(ctx, d.intents, d.reads, start.Episode, through)
	files.SourceState = selection.Front.ID
	d.writeSurfaceJSON(w, r, http.StatusOK, episodeFilesSchema, files)
}

type episodeFileReader struct {
	objectWorldReader
	left int
}

func (r *episodeFileReader) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
	if r.left == 0 {
		return store.Object{}, false, errEpisodeFileBudget
	}
	r.left--
	return r.objectWorldReader.GetObject(ctx, ref)
}

func (r *episodeFileReader) GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error) {
	if r.left == 0 {
		return store.World{}, false, errEpisodeFileBudget
	}
	r.left--
	return r.objectWorldReader.GetWorld(ctx, ref)
}

// buildEpisodeFileActivity uses the CLI report's classification and sorting
// without returning any raw invocation payloads, paths outside the worktree,
// or outcome reasons. A bounded prefix is explicitly partial.
func buildEpisodeFileActivity(ctx context.Context, intents intentReader, reads objectWorldReader, episode string, through time.Time) surfaceEpisodeFiles {
	out := surfaceEpisodeFiles{Schema: episodeFilesSchema, Episode: episode, Status: "unavailable",
		Reads: []surfaceFileUse{}, Writes: []surfaceFileUse{}, Produced: []surfaceFileUse{}}
	if intents == nil || reads == nil || !surface.EpisodeGrammar.MatchString(episode) {
		return out
	}
	b := newReportBuilder(episode)
	b.filesOnly = true
	reader := &episodeFileReader{objectWorldReader: reads, left: surfaceEpisodeFileLookups}
	partial, checked := false, false
	for _, via := range reportSurfaces {
		page, err := intents.InvocationIntents(ctx, string(via)+":"+episode+":", 0, surfaceEpisodeFileIntents+1)
		if err != nil {
			partial = true
			break
		}
		if len(page) > surfaceEpisodeFileIntents {
			page = page[:surfaceEpisodeFileIntents]
			partial = true
		}
		failed := false
		for _, it := range page {
			if it.Outcome == nil || it.Outcome.Status != "committed" ||
				!through.IsZero() && time.Unix(it.Intent.LogicalTime, 0).After(through) {
				continue
			}
			if err := b.add(ctx, reader, string(via), it); err != nil {
				failed = true
				if reader.left == 0 || ctx.Err() != nil {
					break
				}
				continue
			}
			checked = true
		}
		if failed {
			partial = true
			if reader.left == 0 || ctx.Err() != nil {
				break
			}
		} else {
			checked = true
		}
	}
	if !checked && len(b.reads)+len(b.written)+len(b.produced) == 0 {
		return out
	}

	b.finish()
	readsReport, writesReport := b.rep.Reads, b.rep.Writes
	omitted := surfaceFileUseOmitted{Reads: readsReport.Truncated, Writes: writesReport.WrittenTruncated,
		Produced: writesReport.ProducedTruncated, Unlisted: readsReport.Unlisted + writesReport.Unlisted}
	appendUse := func(to *[]surfaceFileUse, dropped *int, use surfaceFileUse) {
		if !safeEpisodeFilePath(use.Path) {
			omitted.Unlisted++
			return
		}
		if len(*to) == surfaceEpisodeFilePaths {
			(*dropped)++
			return
		}
		*to = append(*to, use)
	}
	for _, use := range readsReport.Files {
		appendUse(&out.Reads, &omitted.Reads, surfaceFileUse{Path: use.Path, Count: use.Count, FirstEntry: use.FirstEntry})
	}
	for _, use := range writesReport.Written {
		appendUse(&out.Writes, &omitted.Writes, surfaceFileUse{Path: use.Path, Op: use.Op, Count: use.Count, FirstEntry: use.FirstEntry, lastEntry: use.lastEntry})
	}
	for _, use := range writesReport.Produced {
		appendUse(&out.Produced, &omitted.Produced, surfaceFileUse{Path: use.Path, Count: use.Count, FirstEntry: use.FirstEntry})
	}
	out.DeltaGaps = len(writesReport.DeltasTruncated) + len(writesReport.DeltaErrors) + writesReport.DeltaTruncatedOmitted + writesReport.DeltaErrorsOmitted
	out.Status = "complete"
	if omitted != (surfaceFileUseOmitted{}) {
		out.Omitted = &omitted
	}
	if partial || b.filesIncomplete || out.DeltaGaps > 0 || out.Omitted != nil {
		out.Status = "partial"
	}
	return out
}

func safeEpisodeFilePath(path string) bool {
	return len(path) <= surfaceEpisodeFilePathBytes && utf8.ValidString(path) && filepath.IsLocal(path) &&
		path != "." && filepath.Clean(path) == path && !strings.ContainsAny(path, "\\:") &&
		strings.IndexFunc(path, unicode.IsControl) == -1 &&
		surface.SanitizeOpenText(path, surfaceEpisodeFilePathBytes) == path
}

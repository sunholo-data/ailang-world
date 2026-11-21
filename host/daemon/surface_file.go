package daemon

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacefollow"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const surfaceFileSchema = "world/surface-file/v1"

func (d *Daemon) handleSurfaceFile(w http.ResponseWriter, r *http.Request) {
	refuse := func(status int, reason string) {
		d.writeSurfaceJSON(w, r, status, surfaceFileSchema, surfaceReason{Schema: surfaceFileSchema, Reason: reason})
	}
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		refuse(http.StatusMethodNotAllowed, "use GET for a committed file preview")
		return
	}
	if len(r.URL.RawQuery) > 4*surfacefollow.MaxPreviewPath {
		refuse(http.StatusBadRequest, "expected one brane, state and file path")
		return
	}
	q, err := url.ParseQuery(r.URL.RawQuery)
	path := q.Get("path")
	source := q.Get("source")
	sourceCount := len(q["source"])
	if err != nil || len(q) != 3+sourceCount || sourceCount > 1 ||
		(sourceCount == 1 && source != "episode-read" && source != "episode-write") ||
		len(q["brane"]) != 1 || len(q["at"]) != 1 || len(q["path"]) != 1 ||
		len(path) == 0 || len(path) > surfacefollow.MaxPreviewPath || strings.ContainsRune(path, 0) {
		refuse(http.StatusBadRequest, "expected one brane, state and file path")
		return
	}
	if sourceCount == 1 {
		d.handleSurfaceEpisodeFile(w, r, q, path, source)
		return
	}
	repo, kind, id, idErr := surface.ParseQualifiedID(q.Get("at"))
	braneRepo, braneKind, lane, braneErr := surface.ParseQualifiedID(q.Get("brane"))
	if idErr != nil || braneErr != nil || repo != braneRepo || kind != "commit" || braneKind != "brane" {
		refuse(http.StatusBadRequest, "expected a commit and brane qualified by the same repository")
		return
	}
	if d.follow == nil {
		refuse(http.StatusNotFound, "no committed file preview for this selection")
		return
	}
	merge, change, ok := findLandedOutput(d.followSnapshot().Envelope, repo, lane, id, path)
	if !ok {
		refuse(http.StatusNotFound, "no committed file preview for this selection")
		return
	}
	for _, source := range d.follow.config.Repos {
		if source.ID != repo {
			continue
		}
		// Bound concurrent Git preview work independently of follow probes.
		if d.follow.fileReads.Add(1) > 2 {
			d.follow.fileReads.Add(-1)
			refuse(http.StatusServiceUnavailable, "committed file preview is busy")
			return
		}
		defer d.follow.fileReads.Add(-1)
		file, err := surfacefollow.ReadFilePreview(r.Context(), source, merge, change)
		if err != nil {
			refuse(http.StatusServiceUnavailable, "committed file could not be read")
			return
		}
		d.writeSurfaceJSON(w, r, http.StatusOK, surfaceFileSchema, file)
		return
	}
	refuse(http.StatusNotFound, "no committed file preview for this selection")
}

// findLandedOutput authorizes an exact path on the selected repo, node and
// lane. A step's source is its linked shared merge in the same partition.
func findLandedOutput(env surface.FollowIndex, repo, lane, id, path string) (string, surface.Change, bool) {
	for _, partition := range env.Repos {
		if partition.ID != repo {
			continue
		}
		for _, state := range partition.Index.States {
			if state.ID != id {
				continue
			}
			merge := state.ID
			stateLane := surfacepage.SharedKey
			if state.Kind == surface.KindStep {
				merge = state.Merge
				if state.Brane != "" {
					stateLane = state.Brane
				}
			}
			if lane != stateLane {
				break
			}
			for _, change := range partition.Index.LandedFiles[merge].Changes {
				if change.Path == path {
					return merge, change, true
				}
			}
			break
		}
		break
	}
	return "", surface.Change{}, false
}

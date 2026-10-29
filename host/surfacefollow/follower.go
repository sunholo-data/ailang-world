package surfacefollow

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/sunholo-data/ailang-world/host/surface"
)

// Status is a point-in-time value. Partition is the last accepted projection;
// Head is the latest observed local head and can differ while rebuilding.
// Sink must return promptly and must treat delivered data as immutable.
type Status struct {
	ID          string
	Head        string
	BuiltAt     string
	CheckedAt   string
	LastFailure *surface.FollowFailure
	RetryAt     time.Time
	Partition   *surface.FollowPartition
}

type Sink func(Status)

// Follower is run once, with a caller-owned lifetime. Restart with a new Follower
// to reset backoff and the rolling call census.
type Follower struct {
	Config    Config
	Collector *Collector
	Sink      Sink
}

type repoState struct {
	status   Status
	pending  *job
	active   *job
	failures int
	retryAt  time.Time
}

type job struct {
	repo   int
	head   string
	queued time.Time
	cancel context.CancelFunc
}

type event struct {
	repo      int
	head      string
	err       error
	pollDone  bool
	job       *job
	partition surface.FollowPartition
}

func retryDelay(n int) time.Duration {
	switch n {
	case 0, 1:
		return 15 * time.Second
	case 2:
		return 30 * time.Second
	default:
		return 60 * time.Second
	}
}

// Run probes on a fixed start-based cadence. A single event loop owns all
// scheduling/status state; only git calls and builds run outside it.
func (f *Follower) Run(ctx context.Context) error {
	c := f.Collector
	if c == nil {
		c = NewCollector(nil, nil)
	}
	runCtx, stop := context.WithCancel(ctx)
	var wg sync.WaitGroup
	defer func() {
		stop()
		wg.Wait()
	}()
	states := make([]repoState, len(f.Config.Repos))
	for i, r := range f.Config.Repos {
		states[i].status.ID = r.ID
	}
	events := make(chan event, 16)
	send := func(e event) {
		select {
		case events <- e:
		case <-runCtx.Done():
		}
	}
	emit := func(s *repoState) {
		if f.Sink != nil {
			out := s.status
			if out.Partition != nil {
				p := *out.Partition
				p.CheckedAt, p.LastFailure = out.CheckedAt, out.LastFailure
				out.Partition = &p
			}
			f.Sink(out)
		}
	}
	markFailure := func(s *repoState, err error) {
		s.status.LastFailure = classify(err)
		if strings.HasPrefix(err.Error(), "candidate") {
			s.status.LastFailure = failure("candidate")
		}
		if s.status.LastFailure.Class != "budget" {
			s.failures++
		}
		delay := retryDelay(s.failures)
		if s.status.LastFailure.Class == "budget" {
			delay = PollInterval
		}
		s.retryAt = c.Clock.Now().Add(delay)
		s.status.RetryAt = s.retryAt
		emit(s)
	}
	polling, pollAgain := false, false
	startPoll := func() {
		if polling {
			pollAgain = true
			return
		}
		polling = true
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i, r := range f.Config.Repos {
				head, err := c.Probe(runCtx, r)
				send(event{repo: i, head: head, err: err})
				if runCtx.Err() != nil {
					return
				}
			}
			send(event{pollDone: true})
		}()
	}
	running := 0
	dispatch := func() {
		// Finish probes before starting new history reads; the reserved
		// budget also protects subsequent polls from busy workers.
		if polling {
			return
		}
		for i := range states {
			s := &states[i]
			if running == 2 || s.active != nil || s.pending == nil {
				continue
			}
			j := s.pending
			s.pending, s.active = nil, j
			jobCtx, cancel := context.WithCancel(runCtx)
			j.cancel = cancel
			running++
			wg.Add(1)
			go func(j *job, parent context.Context) {
				defer wg.Done()
				defer j.cancel()
				remaining := JobTimeout - c.Clock.Now().Sub(j.queued)
				jobCtx, closeScope := bounded(parent, c.Clock, remaining)
				defer closeScope()
				var p surface.FollowPartition
				var err error
				if remaining <= 0 {
					err = context.DeadlineExceeded
				} else {
					p, err = c.BuildPinned(jobCtx, f.Config.Repos[j.repo], j.head)
				}
				if jobCtx.Err() != nil {
					err = context.DeadlineExceeded
				}
				send(event{repo: j.repo, job: j, partition: p, err: err})
			}(j, jobCtx)
		}
	}
	next := c.Clock.Now().Add(PollInterval)
	timer := c.Clock.NewTimer(PollInterval)
	defer func() { timer.Stop() }()
	startPoll()
	for {
		select {
		case <-runCtx.Done():
			return ctx.Err()
		case <-timer.C():
			now := c.Clock.Now()
			for !next.After(now) {
				next = next.Add(PollInterval)
			}
			timer = c.Clock.NewTimer(next.Sub(now))
			startPoll()
		case e := <-events:
			if runCtx.Err() != nil {
				return ctx.Err()
			}
			if e.pollDone {
				polling = false
				if pollAgain {
					pollAgain = false
					startPoll()
				}
				dispatch()
				continue
			}
			s := &states[e.repo]
			now := c.Clock.Now()
			if e.job == nil {
				if e.err != nil {
					markFailure(s, e.err)
					continue
				}
				s.status.CheckedAt = now.UTC().Format(time.RFC3339)
				changed := s.status.Head != e.head
				s.status.Head = e.head
				if changed {
					s.failures, s.retryAt = 0, time.Time{}
					s.status.RetryAt = time.Time{}
					s.status.LastFailure = nil
					if s.active != nil {
						s.active.cancel()
					}
					s.pending = nil
				}
				accepted := s.status.Partition != nil && s.status.Partition.Head == e.head
				if accepted {
					s.status.LastFailure = nil
					s.status.RetryAt = time.Time{}
					s.failures, s.retryAt = 0, time.Time{}
				} else if (changed || !now.Before(s.retryAt)) && s.pending == nil &&
					(s.active == nil || s.active.head != e.head) {
					s.pending = &job{repo: e.repo, head: e.head, queued: now}
				}
				emit(s)
				continue
			}
			running--
			s.active = nil
			if e.job.head != s.status.Head {
				// The pending job already names the newest observation.
				// Re-probe as well to reconcile movement since that poll.
				startPoll()
				dispatch()
				continue
			}
			if e.err == nil {
				candidate := surface.FollowIndex{Schema: surface.FollowSchema, Generation: "candidate"}
				for i := range states {
					p := states[i].status.Partition
					if i == e.repo {
						p = &e.partition
					}
					if p != nil {
						candidate.Repos = append(candidate.Repos, *p)
					}
				}
				if err := surface.ValidateFollow(candidate); err != nil {
					e.err = fmt.Errorf("candidate aggregate validation refused")
				} else if b, err := json.Marshal(candidate); err != nil || len(b) > surface.MaxBytes {
					e.err = fmt.Errorf("candidate aggregate bytes refused")
				}
			}
			if e.err != nil {
				markFailure(s, e.err)
			} else {
				s.status.Partition = &e.partition
				s.status.BuiltAt = e.partition.BuiltAt
				s.status.RetryAt = time.Time{}
				s.status.LastFailure = nil
				s.failures, s.retryAt = 0, time.Time{}
				emit(s)
			}
			dispatch()
		}
	}
}

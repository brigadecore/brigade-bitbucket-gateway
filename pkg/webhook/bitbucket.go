package webhook

import (
	"log"

	"github.com/brigadecore/brigade/pkg/brigade"
	"github.com/brigadecore/brigade/pkg/storage"
)

const (
	brigadeJSFile      = "brigade.js"
	hubSignatureHeader = "X-Hub-Signature"
)

type bitbucketHandler struct {
	store storage.Store
}

type fileGetter func(commit, path string, proj *brigade.Project) ([]byte, error)

// NewBitbucketHandler creates a bitbucket handler.
func NewBitbucketHandler(s storage.Store) *bitbucketHandler {
	return &bitbucketHandler{
		store: s,
	}
}

func (s *bitbucketHandler) HandleEvent(repo string, eventType string, rev brigade.Revision, payload []byte, secret string) {

	proj, err := s.store.GetProject(repo)

	if err != nil {
		log.Printf("Project %q not found. No secret loaded. %s", repo, err)
		return
	}

	if proj.SharedSecret == "" {
		log.Printf("No secret is configured for this repo.")
		return
	}

	if proj.SharedSecret != secret {
		log.Printf("Secret mismatch for this repo.")
		return
	}

	if proj.Name != repo {
		log.Printf("!!!WARNING!!! Expected project secret to have name %q, got %q", repo, proj.Name)
	}

	s.build(eventType, rev, payload, proj)

	log.Printf("Build Creation Complete")
}

func truncAt(str string, max int) string {
	if len(str) > max {
		short := str[0 : max-3]
		return short + "..."
	}
	return str
}

func (s *bitbucketHandler) build(eventType string, rev brigade.Revision, payload []byte, proj *brigade.Project) error {

	b := &brigade.Build{
		ProjectID: proj.ID,
		Type:      eventType,
		Provider:  "bitbucket",
		Revision:  &rev,
		Payload:   payload,
	}

	return s.store.CreateBuild(b)
}

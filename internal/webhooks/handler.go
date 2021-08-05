package webhooks

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/webhooks/v6/bitbucket"
	"github.com/pkg/errors"
)

// handler is an implementation of the http.Handler interface that can handle
// webhooks (events) from Bitbucket by delegating to a transport-agnostic
// Service interface.
type handler struct {
	service Service
	hook    *bitbucket.Webhook
}

// handler is an implementation of the http.Handler interface that can handle
// webhooks (events) from Bitbucket by delegating to a transport-agnostic
// Service interface.
func NewHandler(service Service) (http.Handler, error) {
	hook, err := bitbucket.New()
	if err != nil {
		return nil, errors.Wrap(err, "error creating handler")
	}
	return &handler{
		service: service,
		hook:    hook,
	}, nil
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")

	payload, err := h.hook.Parse(
		r,
		bitbucket.IssueCommentCreatedEvent,
		bitbucket.IssueCreatedEvent,
		bitbucket.IssueUpdatedEvent,
		bitbucket.PullRequestApprovedEvent,
		bitbucket.PullRequestCommentCreatedEvent,
		bitbucket.PullRequestCommentDeletedEvent,
		bitbucket.PullRequestCommentUpdatedEvent,
		bitbucket.PullRequestCreatedEvent,
		bitbucket.PullRequestDeclinedEvent,
		bitbucket.PullRequestMergedEvent,
		bitbucket.PullRequestUnapprovedEvent,
		bitbucket.PullRequestUpdatedEvent,
		bitbucket.RepoCommitCommentCreatedEvent,
		bitbucket.RepoCommitStatusCreatedEvent,
		bitbucket.RepoCommitStatusUpdatedEvent,
		bitbucket.RepoForkEvent,
		bitbucket.RepoPushEvent,
		bitbucket.RepoUpdatedEvent,
	)
	if err != nil {
		if err == bitbucket.ErrEventNotFound {
			w.WriteHeader(http.StatusNotImplemented)
		} else {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte("{}")) // nolint: errcheck
		return
	}

	events, err := h.service.Handle(r.Context(), payload)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{}")) // nolint: errcheck
		return
	}

	responseObj := struct {
		EventIDs []string `json:"eventIDs"`
	}{
		EventIDs: make([]string, len(events.Items)),
	}
	for i, event := range events.Items {
		responseObj.EventIDs[i] = event.ID
	}

	responseJSON, err := json.Marshal(responseObj)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{}")) // nolint: errcheck
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseJSON) // nolint: errcheck
}

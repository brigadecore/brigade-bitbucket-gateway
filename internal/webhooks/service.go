package webhooks

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/brigadecore/brigade/sdk/v3"
	"github.com/go-playground/webhooks/v6/bitbucket"
	"github.com/pkg/errors"
)

// Service is an interface for components that can handle webhooks (events) from
// Bitbucket. Implementations of this interface are transport-agnostic.
type Service interface {
	// Handle handles a Bitbucket webhook (event).
	Handle(
		ctx context.Context,
		payload interface{},
	) (sdk.EventList, error)
}

type service struct {
	eventsClient sdk.EventsClient
}

// NewService returns an implementation of the Service interface for handling
// webhooks (events) from Bitbucket.
func NewService(eventsClient sdk.EventsClient) Service {
	return &service{
		eventsClient: eventsClient,
	}
}

// nolint: gocyclo
func (s *service) Handle(
	ctx context.Context,
	payload interface{},
) (sdk.EventList, error) {
	var events sdk.EventList

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return events, errors.Wrap(err, "error marshaling event payload")
	}
	event := sdk.Event{
		Source:  "brigade.sh/bitbucket",
		Payload: string(payloadBytes),
	}

	switch p := payload.(type) {

	// nolint: lll
	// issue:comment_created
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-created
	//
	// A user comments on an issue associated with a repository.
	case bitbucket.IssueCommentCreatedPayload:
		event.Type = string(bitbucket.IssueCommentCreatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}

	// nolint: lll
	// issue:created
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Created
	//
	// A user creates an issue for a repository.
	case bitbucket.IssueCreatedPayload:
		event.Type = string(bitbucket.IssueCreatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}

	// nolint: lll
	// issue:updated
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Updated.1
	//
	// A user updated an issue for a repository.
	case bitbucket.IssueUpdatedPayload:
		event.Type = string(bitbucket.IssueUpdatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}

	// nolint: lll
	// pullrequest:approved
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Approved
	//
	// A user approves a pull request for a repository.
	case bitbucket.PullRequestApprovedPayload:
		event.Type = string(bitbucket.PullRequestApprovedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:comment_created
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-created.1
	//
	// A user comments on a pull request.
	case bitbucket.PullRequestCommentCreatedPayload:
		event.Type = string(bitbucket.PullRequestCommentCreatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:comment_deleted
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-deleted
	//
	// A user deletes a comment on a pull request.
	case bitbucket.PullRequestCommentDeletedPayload:
		event.Type = string(bitbucket.PullRequestCommentDeletedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:comment_updated
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-updated
	//
	// A user updates a comment on a pull request.
	case bitbucket.PullRequestCommentUpdatedPayload:
		event.Type = string(bitbucket.PullRequestCommentUpdatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:created
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Created.1
	//
	// A user creates a pull request for a repository.
	case bitbucket.PullRequestCreatedPayload:
		event.Type = string(bitbucket.PullRequestCreatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:rejected
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Declined
	//
	// A user declines a pull request for a repository.
	case bitbucket.PullRequestDeclinedPayload:
		event.Type = string(bitbucket.PullRequestDeclinedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:fulfilled
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Merged
	//
	// A user merges a pull request for a repository.
	case bitbucket.PullRequestMergedPayload:
		event.Type = string(bitbucket.PullRequestMergedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:unapproved
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Approval-removed
	//
	// A user removes an approval from a pull request for a repository.
	case bitbucket.PullRequestUnapprovedPayload:
		event.Type = string(bitbucket.PullRequestUnapprovedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// pullrequest:updated
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Updated.2
	//
	// A user updates a pull request for a repository.
	case bitbucket.PullRequestUpdatedPayload:
		event.Type = string(bitbucket.PullRequestUpdatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.PullRequest.Source.Commit.Hash,
			Ref:    p.PullRequest.Source.Branch.Name,
		}

	// nolint: lll
	// repo:commit_comment_created
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#hardBreak
	//
	// A user comments on a commit in a repository.
	case bitbucket.RepoCommitCommentCreatedPayload:
		event.Type = string(bitbucket.RepoCommitCommentCreatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.Commit.Hash,
		}

	// nolint: lll
	// repo:commit_status_created
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Build-status-created
	//
	// A build system, CI tool, or another vendor recognizes that a user recently
	// pushed a commit and updates the commit with its status.
	case bitbucket.RepoCommitStatusCreatedPayload:
		event.Type = string(bitbucket.RepoCommitStatusCreatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		url := fmt.Sprintf("%v", p.CommitStatus.Links.Commit)
		urls := strings.Split(url, "/")
		event.Git = &sdk.GitDetails{
			Commit: urls[len(urls)-1],
		}

	// nolint: lll
	// repo:commit_status_updated
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Build-status-updated
	//
	// A build system, CI tool, or another vendor recognizes that a commit has a
	// new status and updates the commit with its status.
	case bitbucket.RepoCommitStatusUpdatedPayload:
		event.Type = string(bitbucket.RepoCommitStatusUpdatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		url := fmt.Sprintf("%v", p.CommitStatus.Links.Commit)
		urls := strings.Split(url, "/")
		event.Git = &sdk.GitDetails{
			Commit: urls[len(urls)-1],
		}

	// nolint: lll
	// repo:fork
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Fork
	//
	// A user forks a repository.
	case bitbucket.RepoForkPayload:
		event.Type = string(bitbucket.RepoForkEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}

	// nolint: lll
	// repo:push
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Push
	//
	// A user pushes 1 or more commits to a repository.
	case bitbucket.RepoPushPayload:
		event.Type = string(bitbucket.RepoPushEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}
		event.Git = &sdk.GitDetails{
			Commit: p.Push.Changes[0].New.Target.Hash,
			Ref:    p.Push.Changes[0].New.Name,
		}

	// nolint: lll
	// repo:updated
	// From https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Updated
	//
	// A user updates the  Name ,  Description ,  Website , or  Language  fields
	// under the  Repository details  page of the repository settings.
	case bitbucket.RepoUpdatedPayload:
		event.Type = string(bitbucket.RepoUpdatedEvent)
		event.Qualifiers = map[string]string{
			"repo": p.Repository.FullName,
		}

	default:
		return events, nil
	}

	events, err = s.eventsClient.Create(ctx, event, nil)
	return events, errors.Wrap(err, "error emitting event(s) into Brigade")
}

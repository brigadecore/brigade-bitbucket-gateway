# Event Reference

This section exists primarily for reference purposes and documents all Bitbucket
webhooks which can be handled by this gateway and their corresponding Brigade
events that may be emitted into the Brigade event bus.

The transformation of a webhook into an event is relatively straightforward and
subject to a few very simple rules:

1. With every webhook handled by this gateway being indicative of activity
   involving some specific repository, the name of the affected repository is
   copied from the webhook's JSON payload and promoted to the `repo` qualifier
   on the corresponding event. This permits projects to subscribe to events
   relating only to specific repositories. Read more about qualifiers
   [here](https://docs.brigade.sh/topics/project-developers/events/#qualifiers).

1. For any webhook that is indicative of activity involving not only a specific
   repository, but also some specific ref (branch or tag) or commit (identified
   by SHA), this gateway copies those details from the webhook's JSON payload
   and promotes them to the corresponding event's `git.ref` and/or `git.commit`
   fields. By doing so, Brigade is enabled to locate specific code referenced by
   the webhook/event. The importance of this cannot be understated, as it is
   what permits Brigade to be used for implementing CI/CD pipelines.

1. For _all_ webhooks, without exception, the entire JSON payload, without any
   modification, becomes the corresponding event's `payload`. The event
   `payload` field is a string field, however, so script authors wishing to
   access the payload will need to parse the payload themselves with a
   `JSON.parse()` call or similar.

The following table summarizes all Bitbucket webhooks that can be handled by
this gateway and the corresponding event(s) that are emitted into Brigade's
event bus.

| Webhook | Scope | Event Type(s) Emitted |
|---------|-------|-----------------------|
[`issue:comment_created`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-created) | specific repository | `issue:comment_created` |
[`issue:created`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Created) | specific repository | `issue:created` |
[`issue:updated`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Updated.1) | specific repository | `issue:updated` |
[`pullrequest:approved`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Approved) | specific commit | `pullrequest:approved` |
[`pullrequest:comment_created`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-created.1) | specific commit | `pullrequest:comment_created` |
[`pullrequest:comment_deleted`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-deleted) | specific commit | `pullrequest:comment_deleted` |
[`pullrequest:comment_updated`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Comment-updated) | specific commit | `pullrequest:comment_updated` |
[`pullrequest:created`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Created.1) | specific commit | `pullrequest:created` |
[`pullrequest:fulfilled`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Merged) | specific commit | `pullrequest:fulfilled` |
[`pullrequest:rejected`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Declined) | specific commit | `pullrequest:rejected` |
[`pullrequest:unapproved`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Approval-removed) | specific commit | `pullrequest:unapproved` |
[`pullrequest:updated`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Updated.2) | specific commit | `pullrequest:updated` |
[`repo:commit_comment_created`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#hardBreak) | specific commit | `repo:commit_comment_created` |
[`repo:commit_status_created`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Build-status-created) | specific commit | `repo:commit_status_created` |
[`repo:commit_status_updated`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Build-status-updated) | specific commit | `repo:commit_status_updated` |
[`repo:fork`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Fork) | specific repository | `repo:fork` |
[`repo:push`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Push) | specific commit | `repo:push` |
[`repo:updated`](https://support.atlassian.com/bitbucket-cloud/docs/event-payloads/#Updated) | specific repository | `repo:updated` |

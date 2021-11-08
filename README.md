# Brigade Bitbucket Gateway

![build](https://badgr.brigade2.io/v1/github/checks/brigadecore/brigade-bitbucket-gateway/badge.svg?branch=v2&appID=99005)
[![codecov](https://codecov.io/gh/brigadecore/brigade-bitbucket-gateway/branch/v2/graph/badge.svg?token=RJyZsepTmV)](https://codecov.io/gh/brigadecore/brigade-bitbucket-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/brigadecore/brigade-bitbucket-gateway)](https://goreportcard.com/report/github.com/brigadecore/brigade-bitbucket-gateway)
[![slack](https://img.shields.io/badge/slack-brigade-brightgreen.svg?logo=slack)](https://kubernetes.slack.com/messages/C87MF1RFD)

<img width="100" align="left" src="logo.png">

This is a work-in-progress
[Brigade 2](https://github.com/brigadecore/brigade/tree/v2)
compatible gateway that receives events 
([webhooks](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html))
from Bitbucket and propagates them into Brigade 2's event bus.

<br clear="left"/>

## Installation

Prerequisites:

* A Bitbucket account

* A Kubernetes cluster:
    * For which you have the `admin` cluster role
    * That is already running Brigade 2
    * Capable of provisioning a _public IP address_ for a service of type
      `LoadBalancer`. (This means you won't have much luck running the gateway
      locally in the likes of kind or minikube unless you're able and willing to
      mess with port forwarding settings on your router, which we won't be
      covering here.)

* `kubectl`, `helm` (commands below require Helm 3.7.0+), and `brig` (the
  Brigade 2 CLI)

### 1. Create a Service Account for the Gateway

__Note:__ To proceed beyond this point, you'll need to be logged into Brigade 2
as the "root" user (not recommended) or (preferably) as a user with the `ADMIN`
role. Further discussion of this is beyond the scope of this documentation.
Please refer to Brigade's own documentation.

Using Brigade 2's `brig` CLI, create a service account for the gateway to use:

```console
$ brig service-account create \
    --id brigade-bitbucket-gateway \
    --description brigade-bitbucket-gateway
```

Make note of the __token__ returned. This value will be used in another step.
_It is your only opportunity to access this value, as Brigade does not save it._

Authorize this service account to create new events:

```console
$ brig role grant EVENT_CREATOR \
    --service-account brigade-bitbucket-gateway \
    --source brigade.sh/bitbucket
```

__Note:__ The `--source brigade.sh/bitbucket` option specifies that this service
account can be used _only_ to create events having a value of
`brigade.sh/bitbucket` in the event's `source` field. _This is a security
measure that prevents the gateway from using this token for impersonating other
gateways._

### 2. Install the Bitbucket Gateway

For now, we're using the [GitHub Container Registry](https://ghcr.io) (which is
an [OCI registry](https://helm.sh/docs/topics/registries/)) to host our Helm
chart. Helm 3.7 has _experimental_ support for OCI registries. In the event that
the Helm 3.7 dependency proves troublesome for users, or in the event that this
experimental feature goes away, or isn't working like we'd hope, we will revisit
this choice before going GA.

First, be sure you are using
[Helm 3.7.0](https://github.com/helm/helm/releases/tag/v3.7.0) or greater and
enable experimental OCI support:

```console
$ export HELM_EXPERIMENTAL_OCI=1
```

As this chart requires custom configuration as described above to function
properly, we'll need to create a chart values file with said config.

Use the following command to extract the full set of configuration options into
a file you can modify:

```console
$ helm inspect values oci://ghcr.io/brigadecore/brigade-bitbucket-gateway \
    --version v2.0.0-alpha.3 > ~/brigade-bitbucket-gateway-values.yaml
```

Edit `~/brigade-bitbucket-gateway-values.yaml`, making the following changes:

* `host`: Set this to the host name where you'd like the gateway to be
  accessible.

* `brigade.apiAddress`: Address of the Brigade API server, beginning with
  `https://`

* `brigade.apiToken`: Service account token from step 2

Save your changes to `~/brigade-bitbucket-gateway-values.yaml` and use the
following command to install the gateway using the above customizations:

```console
$ helm install brigade-bitbucket-gateway \
    oci://ghcr.io/brigadecore/brigade-bitbucket-gateway \
    --version v2.0.0-alpha.3 \
    --create-namespace \
    --namespace brigade-bitbucket-gateway \
    --values ~/brigade-bitbucket-gateway-values.yaml
```

### 3. (RECOMMENDED) Create a DNS Entry

If you installed the gateway without enabling support for an ingress controller,
this command should help you find the gateway's public IP address:

```console
$ kubectl get svc brigade-bitbucket-gateway \
    --namespace brigade-bitbucket-gateway \
    --output jsonpath='{.status.loadBalancer.ingress[0].ip}'
```

If you overrode defaults and enabled support for an ingress controller, you
probably know what you're doing well enough to track down the correct IP without
our help. 😉

With this public IP in hand, edit your name servers and add an `A` record
pointing your domain to the public IP.

### 4. Create Webhooks

In your browser, go to your Bitbucket repository for which you'd like to send
webhooks to this gateway. From the menu on the left, select __Repository
settings__ and then __Webhooks__. On this page, click __Add webhooks__.

* In the __Title__ field, add a name for your webhook. It must be unique to this
repository.

* In the __URL__ field, use a value of the form
  `https://<DNS hostname or publicIP>/events`. Note that Bitbucket will not
  permit URLs that it cannot reach.

* Check the __Active__ checkbox.

* If you're using a self-signed certificate (which you are unless you made
  additional configuration changes when deploying the gateway), check the
  __Skip certificate verification__ checkbox.

* Check any/all triggers for which you'd like a webhook sent to this gateway.

* Click __Save__

__Note:__ Those who are also familiar with the Brigade GitHub Gateway might be
perplexed at the lack of anything along the lines of a "shared secret" when
configuring webhooks in Bitbucket. How then do webhook requests authenticate
themselves to your gateway? In short, they do not. In practice, however, this
does _not_ mean that just anyone can send webhooks (directly) to your gateway. 

This gateway is pre-configured (see Helm chart configuration options) with a
list of allowed IPs / IP ranges for inbound requests. This list reflects the IPs
utilized by Bitbucket for outbound requests. This effectively prevents anyone
except Bitbucket from (successfully) sending webhooks to your gateway.

This strategy does not, however, prevent any random Bitbucket user (who happens
to know the address of your gateway) from configuring their own repositories to
send webhooks your way. However, this matters very little, because Brigade 2
operates on a subscription model and if none of your own Brigade projects
subscribe to events originating from the third-party repository in question,
nothing happens.

### 5. Add a Brigade Project

You can create any number of Brigade projects (or modify an existing one) to
listen for events that were sent from BitBucket to your gateway and, in turn,
emitted into Brigade's event bus. You can subscribe to all event types emitted
by the gateway, or just specific ones.

In the example project definition below, we subscribe to all events emitted by
the gateway, provided they've originated from the fictitious
`example-org/example` repository (see the `repo` qualifier).

```yaml
apiVersion: brigade.sh/v2-beta
kind: Project
metadata:
  id: bitbucket-demo
description: A project that demonstrates integration with Bitbucket
spec:
  eventSubscriptions:
  - source: brigade.sh/bitbucket
    types:
    - *
    qualifiers:
      repo: example-org/example
  workerTemplate:
    defaultConfigFiles:
      brigade.js: |-
        const { events } = require("@brigadecore/brigadier");

        events.on("brigade.sh/bitbucket", "issue:comment_created", () => {
          console.log("Someone created a new issue in the example-org/example repository!");
        });

        events.process();
```

In the alternative example below, we subscribe _only_ to `issue:comment_created`
events:

```yaml
apiVersion: brigade.sh/v2-beta
kind: Project
metadata:
  id: bitbucket-demo
description: A project that demonstrates integration with Bitbucket
spec:
  eventSubscriptions:
  - source: brigade.sh/bitbucket
    types:
    - issue:comment_created
    qualifiers:
      repo: example-org/example
  workerTemplate:
    defaultConfigFiles:
      brigade.js: |-
        const { events } = require("@brigadecore/brigadier");

        events.on("brigade.sh/bitbucket", "issue:comment_created", () => {
          console.log("Someone created a new issue in the example-org/example repository!");
        });

        events.process();
```

Assuming this file were named `project.yaml`, you can create the project like
so:

```console
$ brig project create --file project.yaml
```

Open an issue in the Bitbucket repo for which you configured webhooks to send an
event (webhook) to your gateway. The gateway, in turn, will emit the event into
Brigade's event bus. Brigade should initialize a worker (containerized event
handler) for every project that has subscribed to the event, and the worker
should execute the `brigade.js` script that was embedded in the project
definition.

List the events for the `bitbucket-demo` project to confirm this:

```console
$ brig event list --project bitbucket-demo
```

Full coverage of `brig` commands is beyond the scope of this documentation, but
at this point, additional `brig` commands can be applied to monitor the event's
status and view logs produced in the course of handling the event.

## Events Received and Emitted by this Gateway

Events received by this gateway from Bitbucket are, in turn, emitted into
Brigade's event bus.

Events received from BitBucket vary in _scope of specificity_. All events
handled by this gateway are _at least_ indicative of activity involving some
specific repository in some way -- for instance, a BitBucket user having forked
a repository. Some events, however, are more specific than this, being
indicative of activity involving not only a specific repository, but also some
specific branch, tag, or commit -- for instance, a new pull request has been
opened or a new tag has been pushed. In such cases (and only in such cases),
this gateway includes git reference or commit information in the event that is
emitted into Brigade's event bus. By doing so, Brigade (which has built in _git_
support) is enabled to locate specific code affected by the event.

The following table summarizes all Bitbucket event types that can be received by
this gateway and the corresponding event types that are emitted into Brigade's
event bus.

| Bitbucket Event Type | Scope | Event Type(s) Emitted |
|----------------------|-------|-----------------------|
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

## Examples Projects

See `examples/` for complete Brigade projects that demonstrate various
scenarios.

## Contributing

The Brigade project accepts contributions via GitHub pull requests. The
[Contributing](CONTRIBUTING.md) document outlines the process to help get your
contribution accepted.

## Support & Feedback

We have a slack channel!
[Kubernetes/#brigade](https://kubernetes.slack.com/messages/C87MF1RFD) Feel free
to join for any support questions or feedback, we are happy to help. To report
an issue or to request a feature open an issue
[here](https://github.com/brigadecore/brigade-bitbucket-gateway/issues)

## Code of Conduct

Participation in the Brigade project is governed by the
[CNCF Code of Conduct](https://github.com/cncf/foundation/blob/master/code-of-conduct.md).

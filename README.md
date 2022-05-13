# Brigade Bitbucket Gateway

![build](https://badgr.brigade2.io/v1/github/checks/brigadecore/brigade-bitbucket-gateway/badge.svg?branch=v2&appID=99005)
[![codecov](https://codecov.io/gh/brigadecore/brigade-bitbucket-gateway/branch/v2/graph/badge.svg?token=RJyZsepTmV)](https://codecov.io/gh/brigadecore/brigade-bitbucket-gateway)
[![Go Report Card](https://goreportcard.com/badge/github.com/brigadecore/brigade-bitbucket-gateway)](https://goreportcard.com/report/github.com/brigadecore/brigade-bitbucket-gateway)
[![slack](https://img.shields.io/badge/slack-brigade-brightgreen.svg?logo=slack)](https://kubernetes.slack.com/messages/C87MF1RFD)

<img width="100" align="left" src="logo.png">

The Brigade Bitbucket Gateway receives events 
([webhooks](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html))
from Bitbucket, transforms them into Brigade
[events](https://docs.brigade.sh/topics/project-developers/events/), and emits
them into Brigade's event bus.

<br clear="left"/>

> ⚠️&nbsp;&nbsp;If you are familiar with the Brigade GitHub Gateway, be advised
> that this gateway is not as full-featured as that one. At this time, it merely
> offers parity with its predecessor -- the Brigade v1.x-compatible Bitbucket
> Gateway -- meaning it only handles simple webhooks and does not offer
> any integration with
> [Bitbucket Cloud Apps](https://support.atlassian.com/bitbucket-cloud/docs/bitbucket-cloud-apps-overview/).
>
> A more full-featured successor to this gateway is something the Brigade
> maintainers do wish to explore in collaboration with the Brigade community.
> Please reach out on the
> [Kubernetes/#brigade](https://kubernetes.slack.com/messages/C87MF1RFD) Slack
> channel if you are interested in working on that.

## Creating Webhooks

After [installation](docs/INSTALLATION.md), browse to any of your Bitbucket
repositories for which you'd like to send webhooks to this gateway. From the
menu on the left, select __Repository settings__ and then __Webhooks__. On this
page, click __Add webhooks__.

* In the __Title__ field, add a name for your webhook. It must be unique to this
repository.

* In the __URL__ field, use a value of the form
  `https://<DNS hostname or publicIP>/events`. Note that Bitbucket will not
  permit URLs that it cannot reach.

  > ⚠️&nbsp;&nbsp;Instructions for finding the public IP are in the
  > [installation docs](docs/INSTALLATION.md).

* Check the __Active__ checkbox.

* If you're using a self-signed certificate (again, refer to the
  [installation docs](docs/INSTALLATION.md)), check the __Skip certificate
  verification__ checkbox.

* Check any/all triggers for which you'd like a webhook sent to this gateway.

* Click __Save__

> ⚠️&nbsp;&nbsp;Those who are also familiar with the Brigade GitHub Gateway
> might be perplexed at the lack of anything along the lines of a "shared
> secret" when configuring webhooks in Bitbucket. How then do webhook requests
> authenticate themselves to your gateway? In short, they do not. In practice,
> however, this does _not_ mean that just anyone can send webhooks (directly) to
> your gateway. 
>
> This gateway is pre-configured (see Helm chart configuration options) with a
> list of allowed IPs / IP ranges for inbound requests. This list reflects the
> IPs utilized by Bitbucket for outbound requests. This effectively prevents
> anyone except Bitbucket from (successfully) sending webhooks to your gateway.
>
> This strategy does not, however, prevent any random Bitbucket user (who
> happens to know the address of your gateway) from configuring their own
> repositories to send webhooks your way. However, this matters very little,
> because Brigade 2 operates on a subscription model and if none of your own
> Brigade projects subscribe to events originating from the third-party
> repository in question, nothing happens.

## Subscribing

Now subscribe any number of Brigade
[projects](https://docs.brigade.sh/topics/project-developers/projects/)
to events emitted by this gateway -- all of which have a value of
`brigade.sh/bitbucket` in their `source` field. You can subscribe to all event
types emitted by the gateway, or just specific ones.

In the example project definition below, we subscribe to
`issue:created` events, provided they've originated from the fictitious
`example-org/example` repository (see the `repo` 
[qualifier](https://docs.brigade.sh/topics/project-developers/events/#qualifiers)).
You should adjust this value to match a repository for which you are sending
webhooks to your new gateway (see
[installation instructions](docs/INSTALLATION.md)).

```yaml
apiVersion: brigade.sh/v2
kind: Project
metadata:
  id: bitbucket-demo
description: A project that demonstrates integration with Bitbucket
spec:
  eventSubscriptions:
  - source: brigade.sh/bitbucket
    types:
    - issue:created
    qualifiers:
      repo: example-org/example
  workerTemplate:
    defaultConfigFiles:
      brigade.js: |-
        const { events } = require("@brigadecore/brigadier");

        events.on("brigade.sh/bitbucket", "issue:created", () => {
          console.log("Someone created a new issue in the example-org/example repository!");
        });

        events.process();
```

Assuming this file were named `project.yaml`, you can create the project like
so:

```console
$ brig project create --file project.yaml
```

Creating a new issue in the corresponding repo should now send a webhook from
Bitbucket to your gateway. The gateway, in turn, will emit an event into
Brigade's event bus. Brigade should initialize a worker (containerized event
handler) for every project that has subscribed to the event, and the worker
should execute the `brigade.js` script that was embedded in the project
definition.

List the events for the `bitbucket-demo` project to confirm this:

```console
$ brig event list --project bitbucket-demo
```

Full coverage of `brig` commands is beyond the scope of this documentation, but
at this point,
[additional `brig` commands](https://docs.brigade.sh/topics/project-developers/brig/)
can be applied to monitor the event's status and view logs produced in the
course of handling the event.

## Further Reading

* [Installation](docs/INSTALLATION.md): Check this out if you're an operator who
  wants to integrate Bitbucket with your Brigade installation.
* [Event Reference](docs/EVENT_REFERENCE.md): Check this out if you're a script
  author or contributor who requires detailed information about all the webhooks
  handled by this gateway and the corresponding events the gateway emits into
  Brigade.

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

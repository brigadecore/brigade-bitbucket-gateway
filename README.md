# Brigade Bitbucket Gateway

Send [Bitbucket events](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html) into a [Brigade](https://github.com/brigadecore/brigade) pipeline. 

This is a Brigade gateway that listens to bitbucket webhooks event stream and triggers events inside of Brigade.

## Prerequisites

1. Have a running [Kubernetes](https://kubernetes.io/docs/setup/) environment
2. Setup [Helm](https://github.com/kubernetes/helm)
3. Setup [Brigade](https://github.com/brigadecore/brigade) core

## Install

### From File
Clone Brigade Bitbucket Gateway and change directory
```bash
$ git clone https://github.com/lukepatrick/brigade-bitbucket-gateway
$ cd brigade-bitbucket-gateway
```
Helm install brigade-bitbucket-gateway
> note name and namespace (something important about brigade core)
```bash
$ helm install --name bb-gw ./charts/brigade-bitbucket-gateway
```

### From Repo
Add this project as a helm repo

```bash
$ helm repo add bb-gw https://lukepatrick.github.io/brigade-bitbucket-gateway
$ helm install -n bb-gw bb-gw/brigade-bitbucket-gateway
```

## Building from Source
You must have the Go toolchain, make, and dep installed. For Docker support, you will need to have Docker installed as well. 
See more at [Brigade Developers Guide](https://github.com/brigadecore/brigade/blob/master/docs/topics/developers.md) 
From there:

```bash
$ make build
```
To build a Docker image
```bash
$ make docker-build
```

## Compatibility

| BitBucket Gateway | Brigade Core |
|-------------------|--------------|
| v0.10.0+          | v0.10.0+     |
| v0.1.0            | v0.9.0-      |


## BitBucket Integration
The Default URL for the BitBucket Gateway is at `:7448/events/bitbucket/`. In your BitBucket project, go to Settings -> Webhooks. Depending on how you set up 
your Kubernetes and the BitBucket Gateway will determine your externally accessable host/IP/Port. Out of the box the gateway sets up as LoadBalancer; use the host/Cluster IP and check the BitBucket Gateway Kubernetes Service for the external port (something like 30001).

Enter that IP/Port and URL at the Webhook URL. 

Rather than supplying a Shared Secret like GitHub/GitLab, **you must extract the `X-Hook-UUID` from the BitBucket Webhook created**. Store this value as the Brigade Project *values.yaml* `sharedSecret` property.

Check the boxes for the Trigger events to publish from the BitBucket instance. SSL is optional.

## [Scripting Guide](docs/scripting.md)
tl;dr: Bitbucket Gateway produces 16 events:
`push`,
`repo:commit_comment_created`,
`repo:commit_status_created`,
`repo:commit_status_updated`,
`issue:created`,
`issue:updated`,
`issue:comment_created`,
`pullrequest:created`,
`pullrequest:updated`,
`pullrequest:approved`,
`pullrequest:unapproved`,
`pullrequest:fulfilled`,
`pullrequest:rejected`,
`pullrequest:comment_created`,
`pullrequest:comment_updated`,
`pullrequest:comment_deleted`

and two not supported: `repo:fork`, `repo:updated`


# Contributing

This project welcomes contributions and suggestions.

# License

MIT
# Brigade bitbucket Gateway

Send [Bitbucket events](https://confluence.atlassian.com/bitbucket/manage-webhooks-735643732.html) into a [Brigade](https://github.com/Azure/brigade) pipeline. 

This is a Brigade gateway that listens to bitbucket webhooks event stream and triggers events inside of Brigade.

## Prerequisites

1. Have a running [Kubernetes](https://kubernetes.io/docs/setup/) environment
2. Setup [Helm](https://github.com/kubernetes/helm)
3. Setup [Brigade](https://github.com/Azure/brigade) core

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
$ helm install --name brigade-bb ./charts/brigade-bitbucket-gateway
```

### From Repo
Add this project as a helm repo

```bash
$ helm repo add glgw https://lukepatrick.github.io/brigade-bitbucket-gateway
$ helm install -n brig-bb glgw/brigade-bitbucket-gateway
```

## Building from Source
You must have the Go toolchain, make, and dep installed. For Docker support, you will need to have Docker installed as well. 
See more at [Brigade Developers Guide](https://github.com/Azure/brigade/blob/master/docs/topics/developers.md) 
From there:

```bash
$ make build
```
To build a Docker image
```bash
$ make docker-build
```

## [Scripting Guide](docs/scripting.md)
tl;dr: Bitbucket Gateway produces x events: `push`.


# Contributing

This project welcomes contributions and suggestions.

# License

MIT
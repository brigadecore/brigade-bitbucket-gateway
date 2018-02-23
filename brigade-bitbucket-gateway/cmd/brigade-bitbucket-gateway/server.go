package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	whbitbucket "github.com/lukepatrick/brigade-bitbucket-gateway/pkg/webhook"

	"k8s.io/api/core/v1"

	"github.com/Azure/brigade/pkg/storage/kube"

	"gopkg.in/go-playground/webhooks.v3"
	"gopkg.in/go-playground/webhooks.v3/bitbucket"
)

var (
	kubeconfig  string
	master      string
	namespace   string
	gatewayPort string
)

const (
	path = "/events/bitbucket"
)

func init() {
	flag.StringVar(&kubeconfig, "kubeconfig", "", "absolute path to the kubeconfig file")
	flag.StringVar(&master, "master", "", "master url")
	flag.StringVar(&namespace, "namespace", defaultNamespace(), "kubernetes namespace")
	flag.StringVar(&gatewayPort, "gateway-port", defaultGatewayPort(), "TCP port to use for brigade-bitbucket-gateway")
}

func main() {
	flag.Parse()

	hook := bitbucket.New(&bitbucket.Config{UUID: ""})
	hook.RegisterEvents(HandleMultiple,
		bitbucket.RepoPushEvent,
		bitbucket.RepoForkEvent,
		bitbucket.RepoUpdatedEvent,
		bitbucket.RepoCommitCommentCreatedEvent,
		bitbucket.RepoCommitStatusCreatedEvent,
		bitbucket.RepoCommitStatusUpdatedEvent,
		bitbucket.IssueCreatedEvent,
		bitbucket.IssueUpdatedEvent,
		bitbucket.IssueCommentCreatedEvent,
		bitbucket.PullRequestCreatedEvent,
		bitbucket.PullRequestUpdatedEvent,
		bitbucket.PullRequestApprovedEvent,
		bitbucket.PullRequestUnapprovedEvent,
		bitbucket.PullRequestMergedEvent,
		bitbucket.PullRequestDeclinedEvent,
		bitbucket.PullRequestCommentCreatedEvent,
		bitbucket.PullRequestCommentUpdatedEvent,
		bitbucket.PullRequestCommentDeletedEvent) // Add as many as needed

	err := webhooks.Run(hook, ":"+gatewayPort, path)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}

func defaultNamespace() string {
	if ns, ok := os.LookupEnv("BRIGADE_NAMESPACE"); ok {
		return ns
	}
	return v1.NamespaceDefault
}

func defaultGatewayPort() string {
	if port, ok := os.LookupEnv("BRIGADE_GATEWAY_PORT"); ok {
		return port
	}
	return "7748"
}

// HandleMultiple handles multiple bitbucket events
func HandleMultiple(payload interface{}, header webhooks.Header) {
	log.Println("HandleMultiple Payload..")

	clientset, err := kube.GetClient(master, kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	store := kube.New(clientset, namespace)
	store.GetProjects()

	bbhandler := whbitbucket.NewBitbucketHandler(store)

	var repo, commit, secret string
	secret = strings.Join(header["X-Hook-UUID"], "")

	switch payload.(type) {
	case bitbucket.RepoPushPayload:
		log.Println("case bitbucket.RepoPushPayload")
		release := payload.(bitbucket.RepoPushPayload)

		repo = release.Repository.FullName
		commit = release.Push.Changes[0].Commits[0].Hash

		bbhandler.HandleEvent(repo, "push", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.RepoForkPayload:
		log.Println("case bitbucket.RepoForkPayload")
		release := payload.(bitbucket.RepoForkPayload)

		repo = release.Repository.FullName
		commit = "master"

		bbhandler.HandleEvent(repo, "repo:fork", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.RepoUpdatedPayload:
		log.Println("case bitbucket.RepoUpdatedPayload")
		release := payload.(bitbucket.RepoUpdatedPayload)

		repo = release.Repository.FullName
		commit = "master"

		bbhandler.HandleEvent(repo, "repo:updated", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.RepoCommitCommentCreatedPayload:
		log.Println("case bitbucket.RepoCommitCommentCreatedPayload")
		release := payload.(bitbucket.RepoCommitCommentCreatedPayload)

		repo = release.Repository.FullName
		commit = release.Commit.Hash

		bbhandler.HandleEvent(repo, "repo:commit_comment_created", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.RepoCommitStatusCreatedPayload:
		log.Println("case bitbucket.RepoCommitStatusCreatedPayload")
		release := payload.(bitbucket.RepoCommitStatusCreatedPayload)

		repo = release.Repository.FullName

		url := fmt.Sprintf("%v", release.CommitStatus.Links.Commit)
		urls := strings.Split(url, "/")

		commit = urls[len(urls)-1]

		bbhandler.HandleEvent(repo, "repo:commit_status_created", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.RepoCommitStatusUpdatedPayload:
		log.Println("case bitbucket.RepoCommitStatusUpdatedPayload")
		release := payload.(bitbucket.RepoCommitStatusUpdatedPayload)

		repo = release.Repository.FullName

		url := fmt.Sprintf("%v", release.CommitStatus.Links.Commit)
		urls := strings.Split(url, "/")

		commit = urls[len(urls)-1]

		bbhandler.HandleEvent(repo, "repo:commit_status_updated", commit, []byte(fmt.Sprintf("%v", release)), secret)

	default:
		log.Printf("Unsupported event")
		return
	}

}

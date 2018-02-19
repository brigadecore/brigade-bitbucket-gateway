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
	log.Println("Handling Payload..")

	clientset, err := kube.GetClient(master, kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	store := kube.New(clientset, namespace)
	store.GetProjects()

	glhandler := whbitbucket.NewbitbucketHandler(store)

	var repo, commit, secret string
	secret = strings.Join(header["X-bitbucket-Token"], "")

	switch payload.(type) {
	case bitbucket.PushEventPayload:
		log.Println("case bitbucket.PushEventPayload")
		release := payload.(bitbucket.PushEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.CheckoutSHA

		glhandler.HandleEvent(repo, "push", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.TagEventPayload:
		log.Println("case bitbucket.TagEventPayload")
		release := payload.(bitbucket.TagEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.CheckoutSHA

		glhandler.HandleEvent(repo, "tag", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.IssueEventPayload:
		log.Println("case bitbucket.IssueEventPayload")
		release := payload.(bitbucket.IssueEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "issue", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.ConfidentialIssueEventPayload:
		log.Println("case bitbucket.ConfidentialIssueEventPayload")
		release := payload.(bitbucket.ConfidentialIssueEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "issue", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.CommentEventPayload:
		log.Println("case bitbucket.CommentEventPayload")
		release := payload.(bitbucket.CommentEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.Commit.ID

		glhandler.HandleEvent(repo, "comment", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.MergeRequestEventPayload:
		log.Println("case bitbucket.MergeRequestEventPayload")
		//release := payload.(bitbucket.MergeRequestEventPayload)

		log.Println("do nothing")
		// do nothing, waiting on
		// https://github.com/go-playground/webhooks/pull/24
		// repo = release.
		// commit = release.Commit.ID

		// //repo string, eventType string, commit string, payload []byte, secret)
		// glhandler.HandleEvent(repo, "mergerequest", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.WikiPageEventPayload:
		log.Println("case bitbucket.WikiPageEventPayload")
		release := payload.(bitbucket.WikiPageEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.Project.DefaultBranch

		glhandler.HandleEvent(repo, "wikipage", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.PipelineEventPayload:
		log.Println("case bitbucket.PipelineEventPayload")
		release := payload.(bitbucket.PipelineEventPayload)

		repo = release.Project.PathWithNamespace
		commit = release.Commit.ID

		glhandler.HandleEvent(repo, "pipeline", commit, []byte(fmt.Sprintf("%v", release)), secret)

	case bitbucket.BuildEventPayload:
		log.Println("case bitbucket.BuildEventPayload")
		release := payload.(bitbucket.BuildEventPayload)

		repo = release.ProjectName
		commit = release.SHA

		glhandler.HandleEvent(repo, "build", commit, []byte(fmt.Sprintf("%v", release)), secret)

	default:
		log.Printf("Unsupported event")
		return
	}

}

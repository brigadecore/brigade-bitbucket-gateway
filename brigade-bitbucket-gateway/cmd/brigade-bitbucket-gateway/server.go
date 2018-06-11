package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	whbitbucket "github.com/lukepatrick/brigade-bitbucket-gateway/pkg/webhook"

	"k8s.io/api/core/v1"

	"github.com/Azure/brigade/pkg/brigade"
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

	var repo, secret string
	var rev brigade.Revision
	secret = strings.Join(header["X-Hook-Uuid"], "")

	switch payload.(type) {
	case bitbucket.RepoPushPayload:
		log.Println("case bitbucket.RepoPushPayload")
		release := payload.(bitbucket.RepoPushPayload)

		repo = release.Repository.FullName
		rev.Commit = release.Push.Changes[0].New.Target.Hash
		rev.Ref = release.Push.Changes[0].New.Name

		bbhandler.HandleEvent(repo, "push", rev, getJSON(release), secret)

	case bitbucket.RepoForkPayload:
		log.Println("case bitbucket.RepoForkPayload")
		// No Commit or Ref in a Fork payload, skipping as a supported event
		log.Println("skipping forked event, no Commit or Ref in a Fork payload")
		//release := payload.(bitbucket.RepoForkPayload)
		//repo = release.Repository.FullName
		//rev.Ref = "master"
		//bbhandler.HandleEvent(repo, "repo:fork", rev, getJSON(release), secret)

	case bitbucket.RepoUpdatedPayload:
		log.Println("case bitbucket.RepoUpdatedPayload")
		// No Commit or Ref in a repo:updated payload, skipping as a supported event
		log.Println("skipping repo:updated event, no Commit or Ref in a repo:updated payload")
		// release := payload.(bitbucket.RepoUpdatedPayload)

		// repo = release.Repository.FullName
		// rev.Ref = "master"

		// bbhandler.HandleEvent(repo, "repo:updated", rev, getJSON(release), secret)

	case bitbucket.RepoCommitCommentCreatedPayload:
		log.Println("case bitbucket.RepoCommitCommentCreatedPayload")
		release := payload.(bitbucket.RepoCommitCommentCreatedPayload)

		repo = release.Repository.FullName
		rev.Commit = release.Commit.Hash
		rev.Ref = ""

		bbhandler.HandleEvent(repo, "repo:commit_comment_created", rev, getJSON(release), secret)

	case bitbucket.RepoCommitStatusCreatedPayload:
		log.Println("case bitbucket.RepoCommitStatusCreatedPayload")
		release := payload.(bitbucket.RepoCommitStatusCreatedPayload)

		repo = release.Repository.FullName

		url := fmt.Sprintf("%v", release.CommitStatus.Links.Commit)
		urls := strings.Split(url, "/")

		rev.Commit = urls[len(urls)-1]
		rev.Ref = ""

		bbhandler.HandleEvent(repo, "repo:commit_status_created", rev, getJSON(release), secret)

	case bitbucket.RepoCommitStatusUpdatedPayload:
		log.Println("case bitbucket.RepoCommitStatusUpdatedPayload")
		release := payload.(bitbucket.RepoCommitStatusUpdatedPayload)

		repo = release.Repository.FullName

		url := fmt.Sprintf("%v", release.CommitStatus.Links.Commit)
		urls := strings.Split(url, "/")

		rev.Commit = urls[len(urls)-1]
		rev.Ref = ""

		bbhandler.HandleEvent(repo, "repo:commit_status_updated", rev, getJSON(release), secret)

	case bitbucket.IssueCreatedPayload:
		log.Println("case bitbucket.IssueCreatedPayload")
		release := payload.(bitbucket.IssueCreatedPayload)

		repo = release.Repository.FullName
		rev.Ref = "master"
		rev.Commit = ""

		bbhandler.HandleEvent(repo, "issue:created", rev, getJSON(release), secret)

	case bitbucket.IssueUpdatedPayload:
		log.Println("case bitbucket.IssueUpdatedPayload")
		release := payload.(bitbucket.IssueUpdatedPayload)

		repo = release.Repository.FullName
		rev.Ref = "master"
		rev.Commit = ""

		bbhandler.HandleEvent(repo, "issue:updated", rev, getJSON(release), secret)

	case bitbucket.IssueCommentCreatedPayload:
		log.Println("case bitbucket.IssueCommentCreatedPayload")
		release := payload.(bitbucket.IssueCommentCreatedPayload)

		repo = release.Repository.FullName
		rev.Ref = "master"
		rev.Commit = ""

		bbhandler.HandleEvent(repo, "issue:comment_created", rev, getJSON(release), secret)

	case bitbucket.PullRequestCreatedPayload:
		log.Println("case bitbucket.PullRequestCreatedPayload")
		release := payload.(bitbucket.PullRequestCreatedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:created", rev, getJSON(release), secret)

	case bitbucket.PullRequestUpdatedPayload:
		log.Println("case bitbucket.PullRequestUpdatedPayload")
		release := payload.(bitbucket.PullRequestUpdatedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:updated", rev, getJSON(release), secret)

	case bitbucket.PullRequestApprovedPayload:
		log.Println("case bitbucket.PullRequestApprovedPayload")
		release := payload.(bitbucket.PullRequestApprovedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:approved", rev, getJSON(release), secret)

	case bitbucket.PullRequestUnapprovedPayload:
		log.Println("case bitbucket.PullRequestUnapprovedPayload")
		release := payload.(bitbucket.PullRequestUnapprovedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:unapproved", rev, getJSON(release), secret)

	case bitbucket.PullRequestMergedPayload:
		log.Println("case bitbucket.PullRequestMergedPayload")
		release := payload.(bitbucket.PullRequestMergedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.MergeCommit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:fulfilled", rev, getJSON(release), secret)

	case bitbucket.PullRequestDeclinedPayload:
		log.Println("case bitbucket.PullRequestDeclinedPayload")
		release := payload.(bitbucket.PullRequestDeclinedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:rejected", rev, getJSON(release), secret)

	case bitbucket.PullRequestCommentCreatedPayload:
		log.Println("case bitbucket.PullRequestCommentCreatedPayload")
		release := payload.(bitbucket.PullRequestCommentCreatedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:comment_created", rev, getJSON(release), secret)

	case bitbucket.PullRequestCommentUpdatedPayload:
		log.Println("case bitbucket.PullRequestCommentUpdatedPayload")
		release := payload.(bitbucket.PullRequestCommentUpdatedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:comment_updated", rev, getJSON(release), secret)

	case bitbucket.PullRequestCommentDeletedPayload:
		log.Println("case bitbucket.PullRequestCommentDeletedPayload")
		release := payload.(bitbucket.PullRequestCommentDeletedPayload)

		repo = release.Repository.FullName
		rev.Ref = release.PullRequest.Destination.Branch.Name
		rev.Commit = release.PullRequest.Destination.Commit.Hash

		bbhandler.HandleEvent(repo, "pullrequest:comment_deleted", rev, getJSON(release), secret)

	default:
		log.Printf("Unsupported event")
		return
	}

}

func getJSON(release interface{}) []byte {
	p, err := json.Marshal(release)

	if err != nil {
		log.Fatal(err)
	}
	return []byte(p)
}

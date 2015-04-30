package githubservice

import (
	"strings"

	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

type GithubService struct {
	PersonalAccessToken string
}

func New(personalAccessToken string) *GithubService {
	g := GithubService{
		PersonalAccessToken: personalAccessToken,
	}
	return &g
}

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func (g *GithubService) loadIssuesForAssignee(owner string, assignee string) ([]github.Issue, error) {
	tokenSource := &TokenSource{
		AccessToken: g.PersonalAccessToken,
	}
	oauthClient := oauth2.NewClient(context.TODO(), tokenSource)
	client := github.NewClient(oauthClient)

	var all []github.Issue
	var e error
	opt := &github.SearchOptions{

		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		issueSearchResults, resp, err := client.Search.Issues("user:"+owner+" assignee:"+assignee, opt)

		if err != nil {
			e = err
			break
		}

		all = append(all, issueSearchResults.Issues...)

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	return all, e
}

func (g *GithubService) loadIssuesForRepo(owner string, repo string, assigned string) ([]github.Issue, error) {
	tokenSource := &TokenSource{
		AccessToken: g.PersonalAccessToken,
	}
	oauthClient := oauth2.NewClient(context.TODO(), tokenSource)
	client := github.NewClient(oauthClient)

	var allIssues []github.Issue
	var e error
	opt := &github.IssueListByRepoOptions{
		Assignee:    assigned,
		ListOptions: github.ListOptions{PerPage: 500},
	}

	for {
		issues, resp, err := client.Issues.ListByRepo(owner, repo, opt)

		if err != nil {
			e = err
			break
		}

		allIssues = append(allIssues, issues...)

		if resp.NextPage == 0 {
			break
		}

		opt.ListOptions.Page = resp.NextPage
	}

	return allIssues, e
}

func (g *GithubService) makeIssueList(owner string, repo string, assigned string, lambda func(github.Issue) bool) ([]github.Issue, error) {

	issues, err := g.loadIssuesForRepo(owner, repo, assigned)

	if err != nil {
		return nil, err
	}

	var sprintIssues []github.Issue

	for _, issue := range issues {

		// Backlog items are ones that don't match any of the other categories
		if issue.Labels != nil {
			if lambda(issue) {
				sprintIssues = append(sprintIssues, issue)
			}
		}
	}

	return sprintIssues, err
}

func (g *GithubService) AssignedTo(owner string, repo string, login string) ([]github.Issue, error) {
	if repo == "*" {
		return g.loadIssuesForAssignee(owner, login)

	} else {
		return g.makeIssueList(owner, repo, login, g.any)
	}
}

func (g *GithubService) Sprint(owner string, repo string) ([]github.Issue, error) {
	return g.makeIssueList(owner, repo, "", g.isSprintItem)
}

func (g *GithubService) InProgress(owner string, repo string) ([]github.Issue, error) {
	return g.makeIssueList(owner, repo, "", g.isInProgress)
}

func (g *GithubService) ReadyForQA(owner string, repo string) ([]github.Issue, error) {
	return g.makeIssueList(owner, repo, "", g.isReadyForQA)
}

func (g *GithubService) QAPass(owner string, repo string) ([]github.Issue, error) {
	return g.makeIssueList(owner, repo, "", g.isQAPass)
}

func (g *GithubService) Backlog(owner string, repo string) ([]github.Issue, error) {
	return g.makeIssueList(owner, repo, "", g.isBacklogItem)
}

func (g *GithubService) ReadyForReview(owner string, repo string) ([]github.Issue, error) {
	return g.makeIssueList(owner, repo, "", g.isReadyForReview)
}

func (g *GithubService) getLabelString(labels []github.Label) string {
	var retval string
	for _, label := range labels {
		retval += strings.ToLower(*label.Name) + " "
	}
	return retval
}

func (g *GithubService) any(issue github.Issue) bool {
	return true
}

func (g *GithubService) isBacklogItem(issue github.Issue) bool {

	if g.isSprintItem(issue) || g.isInProgress(issue) || g.isReadyForQA(issue) || g.isQAPass(issue) || g.isDone(issue) {
		return false
	} else {
		return true
	}
}

func (g *GithubService) isProductBacklogItem(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "product") && strings.Contains(label, "backlog")
}

func (g *GithubService) isSprintItem(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "sprint")
}

func (g *GithubService) isInProgress(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "in") && strings.Contains(label, "progress")
}

func (g *GithubService) isReadyForQA(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "ready") && strings.Contains(label, "for") && strings.Contains(label, "qa")
}

func (g *GithubService) isReadyForReview(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "ready") && strings.Contains(label, "for") && strings.Contains(label, "review")
}

func (g *GithubService) isQAPass(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "qa") && strings.Contains(label, "pass")
}

func (g *GithubService) isDone(issue github.Issue) bool {
	label := g.getLabelString(issue.Labels)
	return strings.Contains(label, "done")
}

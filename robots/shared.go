package robots

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
)

// Robots is a slice variable
var Robots = make(map[string]Robot)

// Config is a variable with Configuration allocated
var Config = new(Configuration)

// ConfigDirectory variable uses the flag package to define an arg named c and type string
var ConfigDirectory = flag.String("c", ".", "Configuration directory (default .)")

// CaseInsensitiveSorter sorts String.
type CaseInsensitiveSorter []string

func (a CaseInsensitiveSorter) Len() int      { return len(a) }
func (a CaseInsensitiveSorter) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a CaseInsensitiveSorter) Less(i, j int) bool {
	return strings.ToLower(a[i]) < strings.ToLower(a[j])
}

func init() {

	// Try to load the configuration from the environment and fall back to files in the filesystem
	var c ConfigSpecification
	err := envconfig.Process("marvin", &c)

	if err != nil {
		log.Println(err.Error())
		loadConfigFromFile()
	} else {
		err = json.Unmarshal([]byte(c.Config), Config)
		if err != nil {
			log.Println("error parsing config from env: ", err)
			loadConfigFromFile()
		}
	}

	// This overrides the port in the configuration based on the environment variable PORT
	// for better Heroku happiness.
	if os.Getenv("PORT") != "" {
		Config.Port, _ = strconv.Atoi(os.Getenv("PORT"))
	}
}

func loadConfigFromFile() {
	// Fall back to reading from files if there is an error

	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "config.json")
	config, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error opening config: ", err)
	}

	err = json.Unmarshal(config, Config)
	if err != nil {
		log.Fatal("Error parsing loaded config: ", err)
	}

}

// RegisterRobot is responsible for storing the commands of the robots
func RegisterRobot(command string, r Robot) {
	if _, ok := Robots[command]; ok {
		log.Printf("There are two robots mapped to %s!", command)
	} else {
		log.Printf("Registered: %s", command)
		Robots[command] = r
	}
}

// Send is the method that does the request to the Slack hooks
func (i *IncomingWebhook) Send() error {
	webhook := url.URL{
		Scheme: "https",
		Host:   "hooks.slack.com",
		Path:   "/services/" + Config.WebHookPath,
	}

	jsonPayload, err := json.Marshal(i)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", webhook.String(), bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		message := fmt.Sprintf("ERROR: Non-200 Response from Slack Incoming Webhook API: %s", resp.Status)
		log.Println(message)
	}

	defer resp.Body.Close()

	return err
}

// BuildAttachments is a method that encapsulates the BuildAttachmentsShowRepo method
func BuildAttachments(issues []github.Issue, err error) []Attachment {
	return BuildAttachmentsShowRepo(issues, false, true, err)
}

// BuildAttachmentsShowRepo is the method responsible for return the attachments after build the array of attachments
func BuildAttachmentsShowRepo(issues []github.Issue, showrepo bool, showAssigned bool, err error) []Attachment {

	var attachments []Attachment

	if err == nil {
		if len(issues) > 0 {

			for _, issue := range issues {

				var assigned string

				if issue.Assignee != nil {
					assigned = "Assigned to *" + *issue.Assignee.Login + "*"
				} else {
					assigned = "*Unassigned*"
				}

				var color string = "#A0A0A0"

				var joinedLabels string = ""

				if issue.Labels != nil && len(issue.Labels) > 0 {
					color = *issue.Labels[0].Color

					first := true
					for _, label := range issue.Labels {
						if first == true {
							first = false
						} else {
							joinedLabels += ", "
						}

						joinedLabels += *label.Name
					}
				}

				var title string = "Issue #" + strconv.Itoa(*issue.Number) + ", " + *issue.Title

				var text string = ""
				if showAssigned {
					text += assigned
				}
				if issue.Milestone != nil {
					if showAssigned {
						text += " for "
					}
					text += *issue.Milestone.Title
				}

				if len(joinedLabels) > 0 {
					text += " - [" + joinedLabels + "]"
				}

				markdownFields := []MarkdownField{MarkdownFieldTitle, MarkdownFieldText}

				attachment := &Attachment{
					Title:      title,
					TitleLink:  *issue.HTMLURL,
					Text:       text,
					Color:      color,
					MarkdownIn: markdownFields,
				}

				attachments = append(attachments, *attachment)
			}
		} else {

			attachment := &Attachment{
				Text:  "No Issues Found.",
				Color: "#ff1010",
			}

			attachments = append(attachments, *attachment)
		}

	} else {
		attachment := &Attachment{
			Text:  "Error: " + err.Error(),
			Color: "#ff0000",
		}

		attachments = append(attachments, *attachment)
	}

	return attachments
}

// BuildAttachmentsShowPullRequests is responsible for returning the array of attachments of pull request data
func BuildAttachmentsShowPullRequests(openPRs []github.PullRequest, err error) []Attachment {
	var attachments []Attachment

	if err == nil {
		if len(openPRs) > 0 {
			for _, pullRequest := range openPRs {
				var numberOfDays = time.Since(*pullRequest.CreatedAt).Hours() / 24

				// Setup a color gradient to quickly show PR age for someone scanning the list
				var colour = "#ffe7e7"
				if numberOfDays > 30 {
					colour = "#ff1010"
				} else if numberOfDays > 7 {
					colour = "#ff5757"
				} else if numberOfDays > 3 {
					colour = "#ff9f9f"
				}

				var assigned string
				if pullRequest.User != nil {
					assigned = "_" + *pullRequest.User.Login + "_"
				} else {
					assigned = "_Unassigned_"
				}

				var title string = "PR #" + strconv.Itoa(*pullRequest.Number) + " - " + *pullRequest.Title
				var description string = "*" + strconv.FormatFloat(numberOfDays, 'f', 0, 64) + " days in " + *pullRequest.Head.Repo.Name + "* created by " + assigned
				markdownFields := []MarkdownField{MarkdownFieldTitle, MarkdownFieldText}
				attachment := &Attachment{
					Title:      title,
					TitleLink:  *pullRequest.HTMLURL,
					Text:       description,
					Color:      colour,
					MarkdownIn: markdownFields,
				}
				attachments = append(attachments, *attachment)
			}
		} else {
			attachment := &Attachment{
				Text:  "No active repositories have any open pull requests.",
				Color: "#A0A0A0",
			}
			attachments = append(attachments, *attachment)
		}
	} else {
		attachment := &Attachment{
			Text:  "Error: " + err.Error(),
			Color: "#ff0000",
		}
		attachments = append(attachments, *attachment)
	}

	return attachments
}

// BuildAttachmentsShowCommits is responsible for returning the array of attachments of show commits
func BuildAttachmentsShowCommits(repos map[string][]github.RepositoryCommit, err error) []Attachment {
	var attachments []Attachment
	sortedRepoNames := make([]string, len(repos))

	for repoName := range repos {
		sortedRepoNames = append(sortedRepoNames, repoName)
	}
	sort.Strings(sortedRepoNames)

	if err == nil {
		for _, repoName := range sortedRepoNames {
			repoCommits := repos[repoName]

			if len(repoCommits) > 0 {
				for _, commit := range repoCommits {
					var author string
					if commit.Author != nil {
						author = "_" + *commit.Author.Login + "_"
					} else {
						author = "_Unknown_"
					}

					var title string = repoName + "/" + (*commit.SHA)[0:7] + " - " + *commit.Commit.Message
					var description string = commit.Commit.Author.Date.Format("January _2 3:04PM") + " by " + author
					markdownFields := []MarkdownField{MarkdownFieldTitle, MarkdownFieldText}
					attachment := &Attachment{
						Title:      title,
						TitleLink:  *commit.HTMLURL,
						Text:       description,
						Color:      "#ff1010",
						MarkdownIn: markdownFields,
					}
					attachments = append(attachments, *attachment)
				}
			}
		}
	} else {
		attachment := &Attachment{
			Text:  "Error: " + err.Error(),
			Color: "#ff0000",
		}
		attachments = append(attachments, *attachment)
	}

	return attachments
}

// BuildAttachmentCommitSummaryByRepo is responsible for returning the array of attachments of commit summaries
func BuildAttachmentCommitSummaryByRepo(reposToCommits map[string][]github.RepositoryCommit, owner string, days int) []Attachment {
	var attachments []Attachment

	//Sort list by RepoName
	repos := make([]string, 0, len(reposToCommits))
	for key := range reposToCommits {
		repos = append(repos, key)
	}
	sort.Sort(CaseInsensitiveSorter(repos))

	for _, repoName := range repos {
		var commitList []string
		commitsToMaster := reposToCommits[repoName]

		for _, commit := range commitsToMaster {
			commitList = append(commitList, (*commit.SHA)[0:7])
		}
		commitListString := strings.Join(commitList, ", ")
		commitWording := "commits"
		if len(commitList) == 1 {
			commitWording = "commit"
		}
		attachment := &Attachment{
			Title:      repoName,
			TitleLink:  "https://www.github.com/" + owner + "/" + repoName + "/commits/master",
			Text:       strconv.Itoa(len(commitList)) + " " + commitWording + ": " + commitListString,
			Color:      colorForMasterCommitCount(len(commitList)),
			MarkdownIn: []MarkdownField{MarkdownFieldText},
		}
		attachments = append(attachments, *attachment)
	}

	return attachments
}

func colorForMasterCommitCount(commitCount int) string {
	if commitCount > 10 {
		return "#FF1010"
	} else if commitCount > 5 {
		return "#FF7222"
	} else if commitCount > 0 {
		return "#FFD334"
	} else {
		return "#A0A0A0"
	}
}

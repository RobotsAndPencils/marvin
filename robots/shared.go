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
	"strconv"
	"time"

	"github.com/RobotsAndPencils/marvin/githubservice"
	"github.com/google/go-github/github"
	"github.com/kelseyhightower/envconfig"
)

var Robots = make(map[string]Robot)
var Config = new(Configuration)
var ConfigDirectory = flag.String("c", ".", "Configuration directory (default .)")

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

func RegisterRobot(command string, r Robot) {
	if _, ok := Robots[command]; ok {
		log.Printf("There are two robots mapped to %s!", command)
	} else {
		log.Printf("Registered: %s", command)
		Robots[command] = r
	}
}

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
	return err
}

func BuildAttachments(issues []github.Issue, err error) []Attachment {
	return BuildAttachmentsShowRepo(issues, false, true, err)
}

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

func BuildAttachmentsShowPullRequests(reposWithPRs []githubservice.RepositoryPullRequest, err error) []Attachment {
	var attachments []Attachment

	if err == nil {
		if len(reposWithPRs) > 0 {
			for _, item := range reposWithPRs {
				for _, pullRequest := range item.PullRequests {
					var numberOfDays = time.Since(*pullRequest.CreatedAt).Hours() / 24
					if numberOfDays < 1 {
						break // These PRs are too new for us to care about
					}

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
					var description string = strconv.FormatFloat(numberOfDays, 'f', 0, 64) + " days open in *" + *item.Repository.Name + "* by " + assigned
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

func BuildAttachmentsShowCommits(repoCommits []github.RepositoryCommit, err error) []Attachment {
	var attachments []Attachment

	if err == nil {
		if len(repoCommits) > 0 {
			for _, commit := range repoCommits {
				var author string
				if commit.Author != nil {
					author = "_" + *commit.Author.Login + "_"
				} else {
					author = "_Unknown_"
				}

				var title string = "Commit " + (*commit.SHA)[0:7]
				var description string = "by " + author + " - " + *commit.Commit.Message
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
		} else {
			attachment := &Attachment{
				Text:  "No commits to master were found.",
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

func BuildAttachmentCommitSummary(repo string, masterCommitCount int, totalCommits int) Attachment {

	return Attachment{
		Title: "Commits to master on " + repo,
		Text:  strconv.Itoa(masterCommitCount) + " of " + strconv.Itoa(totalCommits) + " commits were direct to master",
		Color: "#A0A0A0",
	}
}

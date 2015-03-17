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
	return BuildAttachmentsShowRepo(issues, false, err)
}

func BuildAttachmentsShowRepo(issues []github.Issue, showrepo bool, err error) []Attachment {

	var attachments []Attachment

	if err == nil {
		if len(issues) > 0 {

			for _, issue := range issues {

				var name string

				if issue.Assignee != nil {
					name = *issue.Assignee.Login
				} else {
					name = "Unassigned"
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

				var title string = "Issue #" + strconv.Itoa(*issue.Number) + ", " + name + " [" + joinedLabels + "]"

				if issue.Milestone != nil {
					title += " - " + *issue.Milestone.Title
				}

				attachment := &Attachment{
					Title:     title,
					TitleLink: *issue.HTMLURL,
					Text:      *issue.Title,
					Color:     color,
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

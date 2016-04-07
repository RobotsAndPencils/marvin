package robots

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/RobotsAndPencils/marvin/githubservice"
	"github.com/kelseyhightower/envconfig"
)

type OpenPullRequestsBot struct {
}

var OpenPullRequestsConfig = new(GithubConfiguration)

// Loads the config file and registers the bot with the server for command /OpenPullRequests.
func init() {
	// Try to load the configuration from the environment and fall back to files in the filesystem
	var c ConfigSpecification
	err := envconfig.Process("github", &c)

	if err != nil {
		log.Println(err.Error())

		// Fall back to reading from files if there is an error
		loadOpenPullRequestsConfigFromFile()
	} else {
		err = json.Unmarshal([]byte(c.Config), OpenPullRequestsConfig)
		if err != nil {
			log.Println("error parsing config: ", err)
			loadOpenPullRequestsConfigFromFile()
		}
	}
	OpenPullRequests := &OpenPullRequestsBot{}
	RegisterRobot("openpullrequests", OpenPullRequests)
}

func loadOpenPullRequestsConfigFromFile() {
	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "github.json")
	if _, err := os.Stat(configFile); err == nil {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("ERROR: Error opening github config: %s", err)
			return
		}
		err = json.Unmarshal(config, OpenPullRequestsConfig)
		if err != nil {
			log.Printf("ERROR: Error parsing github config: %s", err)
			return
		}
	} else {
		log.Printf("WARNING: Could not find configuration file github.json in %s", *ConfigDirectory)
	}
}

func (r OpenPullRequestsBot) parsePayload(p *Payload) (daysPROpen int, projectLastActiveDays int) {
	output := strings.Split(strings.TrimSpace(p.Text), " ")

	var err error
	var validDaysPROpen int = 1
	if len(output) >= 1 {
		validDaysPROpen, err = strconv.Atoi(output[0])
		if err != nil {
			validDaysPROpen = 1 // Defaults to being open longer than 1 day
		}
	}

	var validDaysSinceLastProjectActivity int = 30
	if len(output) >= 2 {
		validDaysSinceLastProjectActivity, err = strconv.Atoi(output[1])
		if err != nil {
			validDaysSinceLastProjectActivity = 30 // Defaults to projects with activity in the last 30 days
		}
	}

	return validDaysPROpen, validDaysSinceLastProjectActivity
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r OpenPullRequestsBot) Run(p *Payload) string {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)

	daysPROpen, daysSinceLastProjectActivity := r.parsePayload(p)

	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.
	return "Finding pull requests that have been open longer than " + strconv.Itoa(daysPROpen) + " days in projects with activity in the last " + strconv.Itoa(daysSinceLastProjectActivity) + " days..."
}

func (r OpenPullRequestsBot) DeferredAction(p *Payload) {

	daysPROpen, daysSinceLastProjectActivity := r.parsePayload(p)

	service := githubservice.New(OpenPullRequestsConfig.PersonalAccessToken)
	pullRequests, err := service.OpenPullRequests(OpenPullRequestsConfig.Owner, daysPROpen, daysSinceLastProjectActivity)

	attachments := BuildAttachmentsShowPullRequests(pullRequests, err)

	var text string = "Pull requests *open for more than " + strconv.Itoa(daysPROpen) + " day(s)*"

	// Let's use the IncomingWebhook struct defined in definitions.go to form and send an
	// IncomingWebhook message to slack that can be seen by everyone in the room. You can
	// read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
	// You can also see what data is available from the command structure in definitions.go
	response := &IncomingWebhook{
		Channel:     p.ChannelID,
		Username:    "Marvin",
		Text:        text,
		IconEmoji:   ":robot:",
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
		Markdown:    true,
		Attachments: attachments,
	}

	response.Send()
}

func (r OpenPullRequestsBot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "This is a description for openPullRequestsBot which will be displayed on /c"
}

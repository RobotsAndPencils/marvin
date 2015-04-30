package robots

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/RobotsAndPencils/marvin/githubservice"
	"github.com/kelseyhightower/envconfig"
)

type ReadyForReviewBot struct {
}

var ReadyForReviewConfig = new(GithubConfiguration)

// Loads the config file and registers the bot with the server for command /ReadyForReview.
func init() {
	// Try to load the configuration from the environment and fall back to files in the filesystem
	var c ConfigSpecification
	err := envconfig.Process("github", &c)

	if err != nil {
		log.Println(err.Error())

		// Fall back to reading from files if there is an error
		loadReadyForReviewConfigFromFile()
	} else {
		err = json.Unmarshal([]byte(c.Config), ReadyForReviewConfig)
		if err != nil {
			log.Println("error parsing config: ", err)
			loadReadyForReviewConfigFromFile()
		}
	}
	ReadyForReview := &ReadyForReviewBot{}
	RegisterRobot("readyforreview", ReadyForReview)
}

func loadReadyForReviewConfigFromFile() {
	flag.Parse()
	configFile := filepath.Join(*ConfigDirectory, "github.json")
	if _, err := os.Stat(configFile); err == nil {
		config, err := ioutil.ReadFile(configFile)
		if err != nil {
			log.Printf("ERROR: Error opening github config: %s", err)
			return
		}
		err = json.Unmarshal(config, ReadyForReviewConfig)
		if err != nil {
			log.Printf("ERROR: Error parsing github config: %s", err)
			return
		}
	} else {
		log.Printf("WARNING: Could not find configuration file github.json in %s", *ConfigDirectory)
	}
}

// All Robots must implement a Run command to be executed when the registered command is received.
func (r ReadyForReviewBot) Run(p *Payload) string {
	// If you (optionally) want to do some asynchronous work (like sending API calls to slack)
	// you can put it in a go routine like this
	go r.DeferredAction(p)
	// The string returned here will be shown only to the user who executed the command
	// and will show up as a message from slackbot.

	return "Calculating ready for review for " + p.Text + "..."
}

func (r ReadyForReviewBot) DeferredAction(p *Payload) {

	service := githubservice.New(ReadyForReviewConfig.PersonalAccessToken)
	issues, err := service.ReadyForReview(ReadyForReviewConfig.Owner, strings.TrimSpace(p.Text))

	attachments := BuildAttachments(issues, err)

	// Let's use the IncomingWebhook struct defined in definitions.go to form and send an
	// IncomingWebhook message to slack that can be seen by everyone in the room. You can
	// read the Slack API Docs (https://api.slack.com/) to know which fields are required, etc.
	// You can also see what data is available from the command structure in definitions.go
	response := &IncomingWebhook{
		Channel:     p.ChannelID,
		Username:    "Marvin",
		Text:        "Ready for Review for repo *" + p.Text + "*",
		IconEmoji:   ":robot:",
		UnfurlLinks: true,
		Parse:       ParseStyleFull,
		Markdown:    true,
		Attachments: attachments,
	}

	response.Send()
}

func (r ReadyForReviewBot) Description() (description string) {
	// In addition to a Run method, each Robot must implement a Description method which
	// is just a simple string describing what the Robot does. This is used in the included
	// /c command which gives users a list of commands and descriptions
	return "This is a description for readyforreview Bot which will be displayed on /c"
}

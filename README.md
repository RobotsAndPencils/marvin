# Say Hi to Marvin

Marvin is Slack bot written in Go that allows Slack users to use Slash Commands to fetch a list of Github Issues that match a specific criteria.  We use a service called Waffle (http://waffle.io) that allows us to view our Issues in a kanban style board.  Each lane on Waffle has a corresponding label within Github issues.  We use these labels to keep track of which stage each issue is in across our entire company.

Example:

type `/backlog RepoName` into Slack and it will display the current backlog items for that project.

or `/readyforqa RepoName` and you’ll see all the items currently in progress.

![](https://dl.dropboxusercontent.com/s/l984qm2t9j2yfao/D3D7D390-F586-4C72-BA54-45251A252C1D-5045-00009A0112DCB887.gif?dl=0)

`/assigned * GithubName` will fetch all tasks assigned to that user under the specified organization.
or
```/assigned RepoName GitHub```  Will fetch all tasks assigned to that user under that repo.

# Running Locally

Setup your go environment and make sure you have the `$GOPATH` variable defined as well as include `$GOPATH/bin` in your `$PATH`.

Pull down the latest source code from Github for Marvin and after you should be able to see the files under `$GOPATH/src` directory

```
go get github.com/RobotsAndPencils/marvin
```

Create the `config.json` and `github.json` files and put your specific slack and github information into it. Save these files in the `$GOPATH/src/github.com/RobotsAndPencils/marvin` directory.

## config.json

```
{ 
        "domain": "YOUR SLACK DOMAIN HERE", 
        "port": 4444, 
        "webhookpath": "RNP SLACK WEBHOOKPATH HERE"
}
```

**domain** is the subdomain of the slack url used by your company. For example it would be the `mycompany` part of `mycompany.slack.com`

**webhookpath** This will need to be set up in Slack's incoming webhooks integration. If the integration has already been set up you can find the value in Slack settings: `Integrations > Configured Integrations > Incoming WebHooks > #channel > Webhook URL`. Where `#channel` is the slack channel that the webhook is set up to post to.

## github.json

```
{ 
        "owner": "YOUR GITHUB ID HERE",
        "personalAccessToken": "YOUR GITHUB ACCESS TOKEN HERE" 
}
```

**owner** is the userid or organization id that you want to use that has access to all of the repos you want the bot to be able to query.

**personalAccessToken** can be found by going to `settings/profile` on github.com and selecting `Generate New Token` from the `Personal Access Token` tab.

Then, you can run/test the programs locally after initializing 

```
make init
make run
make test
```

To test, use a web posting tool like `Postman` to sent a POST to `http://localhost:4444/slack` with `x-www-form-urlencoded` with at least a `command` and `text` parameters. For example:

```
command /assigned
text myrepository somelogin
```

This should, by default, output information in the default channel you set up when you created the Incoming Webhook.

# Deploying to Heroku

The app can be easily deployed to Heroku, but you need to set up a few config variables to make it work. First, you need to set `BUILDPACK_URL` to be 

```
https://github.com/kr/heroku-buildpack-go.git
```

Then, set `GITHUB_CONFIG` variable to the following JSON with your variables replaced:

```
{ 
	"owner": "YOUR ORG NAME", 
  	"personalAccessToken" : "PERSONAL ACCESS TOKEN WITH ACCESS TO ALL REPOS OF ORG" 
}
```
Your Org Name should be substituted by your organization name, that you want to be able to run the queries for.

The personal access token should be from a user with read access to all of the repos in your org. Don't give it too much power. Create one of these in your `Personal Settings > Applications > Personal Access Tokens` section of your GitHub settings.

Then, you need to set `MARVIN_CONFIG` to json like the following:

```
{
	"domain": "SLACK DOMAIN - EVERYTHING BEFORE .slack.com", 
	"webhookpath" : "INCOMING WEBHOOK PATH AFTER /services/"
}
```

The webhook path is everything after the /services/ in an incoming webhook that you've created. It'll look like three randomized strings of characters with slashes between.

Make sure you also update your GoDeps. `godep save`, then commit the changes to your git repo.

# Setting up your Slack

You need to create a bunch of Slack Slash commands, each pointing at the same URL with the following command names and parameters:

```
/backlog [repo]
/sprint [repo]
/inprogress [repo]
/readyforqa [repo]
/qapass [repo]
/assigned [repo|*] [login]
```

The URL you need to configure will be `https://herokudomain.herokuapp.com/slack`.

Also, you need to create an Incoming Webhook integration and use the end part of the webhook path for parts of the configuration above.

# Acknowledgements

This code is heavily influenced by: `https://github.com/trinchan/slackbot` with only slight divergences. The incoming webhook was changed to a post, and the configuration can come from either the environment or files on disk.

# License

The MIT License (MIT)

Copyright (c) 2015, Robots and Pencils

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.


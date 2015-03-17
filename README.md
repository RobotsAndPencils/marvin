# marvin
Slack bot written in Go

When we move to waffle, it will open the doors for PMs to build tools that will streamline our development process.

One of the benefits is the standardization of the boards across the entire company.  If every project is set up the same, it means that we can built tools that can leverage this formatting and allow us to improve our daily standup and review processes.

One of the tools I would like to build is a bot that fetches names and assignees of issues in Github by a simple command in slack.

Ex:

type `/backlog blueotter` into slack and it will display the current backlog items for that project.


or type `/inprogress blueotter` and you’ll see all the items currently in progress.

This would help the company as a whole keep on top of who is working on what since we can expand this further in future versions that would allow us to type

```/assigned * login```
or
```/assigned repo login```

and it will show all the tasks he’s currently assigned to across the entire company.


This would save Project Managers a lot of time.

Consider this:

Even if this saved each PM 10 mins a day, that would amount to 44 hours a year x 5 project managers = 220 hours per year.

This would not only be used by project managers, but robots and pencils as well to quickly see what they are assigned to at any given point.

If you consider it saving each developer and pencil 10 mins a day, that would save the entire team roughly 3000 hours a year.
Thanks!

# Running Locally

Edit the `config.json` and `github.json` files and put your specific slack and github information into it. Then, you can run

```
make init
make run
```

To test, use a web posting tool like `Postman` to sent a POST to `http://localhost:4444/slack` with `x-www-form-urlencoded` with at least a `command` and `text` parameters. For example:

```
command /assigned
text myrepository somelogin
```

This should, by default, output information in the defauly channel you set up when you created the Incoming Webhook.

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


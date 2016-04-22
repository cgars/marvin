package main

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/gicmo/webhooks"
	"github.com/gicmo/webhooks/github"
)

func (b *Bot) HandlePullRequest(payload interface{}) {

	if !b.conn.Connected() {
		return
	}

	pl := payload.(github.PullRequestPayload)

	name := pl.PullRequest.Base.Repo.FullName
	action := pl.Action
	number := pl.Number
	title := pl.PullRequest.Title
	from := pl.PullRequest.Head.Label
	to := pl.PullRequest.Base.Label
	sender := pl.Sender.Login

	b.conn.Privmsgf("#gnode", "[%s#%d] '%s' [%s → %s] %s (%s)\n", name, number,
		title, from, to, action, sender)
}

func (b *Bot) HandleStatus(payload interface{}) {
	pl := payload.(github.StatusPayload)

	state := pl.State

	if state == "pending" {
		return
	}

	name := pl.Name
	cks := pl.SHA

	comps := strings.Split(pl.Context, "/")
	service := comps[0]
	if len(comps) > 1 {
		service = comps[1]
	}

	out := bytes.NewBufferString("")
	out.WriteString(fmt.Sprintf("[%s] %.7s %s %s", name, cks, service, state))

	if state == "failure" {
		out.WriteString(fmt.Sprintf(" [%s]", pl.TragetURL))
	}

	if service == "coveralls" {
		out.WriteString(" (")
		out.WriteString(pl.Desctiption)
		out.WriteString(")")
	}

	out.WriteString("\n")

	if b.conn.Connected() {
		b.conn.Privmsg("#gnode", out.String())
	}
}

func (b *Bot) SetupGithub() {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	if secret == "" {
		fmt.Fprintf(os.Stderr, "No webook secret. No travis integration. :()")
		return
	}

	hook := github.New(&github.Config{Secret: secret})
	hook.RegisterEvents(b.HandlePullRequest, github.PullRequestEvent)
	hook.RegisterEvents(b.HandleStatus, github.StatusEvent)

	go func() {
		port := os.Getenv("GITHUB_WEBHOOK_PORT")
		if port == "" {
			port = "2323"
		}

		fmt.Printf("Listening on :%s/webhooks\n", port)
		err := webhooks.Run(hook, ":"+port, "/webhooks")

		if err != nil {
			fmt.Printf("Error starting webhook listener: %v\n", err)
		}
	}()
}

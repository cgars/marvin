package main

import (
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

	txt := Text("[").Fg(LightBlue, "%s", name).Text("#").Fg(Blue, "%d", number).Text("]")
	txt.Text(" '%s' ", title)
	txt.Fg(LightGrey, "[%s → %s]", from, to).Fg(Green, " %s ", action).Fg(Orange, "(%s)", sender)

	txt.Send(b.conn, "#gnode")
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

	out := Text("[%s] %.7s %s ", name, cks, service)

	if state == "failure" {
		out.Fg(LightRed, "FAIL").Text(" [%s]", pl.TragetURL)
	} else {
		out.Fg(LightGreen, "%s", state)
	}

	if service == "coveralls" {
		out.Text(" (%s)", pl.Desctiption)
	}

	if b.conn.Connected() {
		out.Send(b.conn, "#gnode")
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

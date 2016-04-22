package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/G-Node/marvin/mensa"
	"github.com/G-Node/marvin/quotes"
	irc "github.com/fluffle/goirc/client"

	"github.com/gicmo/webhooks"
	"github.com/gicmo/webhooks/github"
)

type Bot struct {
	conn *irc.Conn

	quit chan bool
}

func (b *Bot) onConncted(conn *irc.Conn, line *irc.Line) {
	conn.Join("#gnode")
}

func (b *Bot) onPrivMessage(conn *irc.Conn, line *irc.Line) {
	text := line.Text()
	target := line.Target()

	fmt.Printf("[D] {pm}: [%s] %s\n", target, text)

	if strings.HasPrefix(text, "mensa") {
		mc := &mensa.Client{Address: "http://openmensa.org/api/v2"}

		var meals []mensa.Meal
		if strings.Contains(text, "tomorrow") {
			meals, _ = mc.MealsForTomorrow("134")
			// ignored error for now
		} else {
			meals, _ = mc.MealsForToday("134")
		}
		if len(meals) == 0 {
			conn.Privmsgf(target, "No milk today, my love has gone away...")
			return
		}
		if strings.Contains(text, "beilagen") {
			b.postMeals(conn, target, meals, []string{"Beilagen"})
			return
		}
		b.postMeals(conn, target, meals, []string{"Tagesgericht", "Aktionsessen", "Biogericht"})
	}

	if strings.Contains(text, "nix") {
		conn.Privmsg(target, "https://youtu.be/Go4SI5ie7qE")
	}

	if strings.Contains(text, "gnode learn") {
		idx := strings.LastIndex(text, "gnode learn")
		err := quotes.LearnQuote(text[idx+11:])
		if err != nil {
			conn.Privmsg(target, "I think you ought to know I'm feeling very depressed.")
			return
		}
		conn.Privmsgf(target, "I just learned smth. new!")
	}
}

func (b *Bot) postMeals(conn *irc.Conn, target string, meals []mensa.Meal, catToUse []string) {

	for _, meal := range meals {
		category := meal.Category
		if !stringInSlice(category, catToUse) {
			continue
		}
		var prices []string
		for key, value := range meal.Prices {
			if value != 0. {
				prices = append(prices, fmt.Sprintf("%s:%.2f€",
					key, value))
			}
		}
		notes := mensa.Emojify(strings.Join(meal.Notes, ", "))
		conn.Privmsgf(target, "%s [%s] [%s]", meal.Name, notes,
			mensa.Emojify(strings.Join(prices, ", ")))
	}
}

func (b *Bot) onJoin(conn *irc.Conn, line *irc.Line) {
	target := line.Target()
	nick := line.Nick
	fmt.Printf("[D] {pm}: [%s] %s\n", line.Nick, line.Text())
	quote, err := quotes.GetRandomQuote()
	if err != nil {
		fmt.Println(err)
	}
	conn.Privmsgf(target, "Welcome %s! Did you know that %s said \"%s\"", nick, quote.Author, quote.Txt)
}

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func NewBot() *Bot {
	cfg := irc.NewConfig("gnode", "marvin", "Metal Man")
	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{ServerName: "irc.freenode.net"}
	cfg.Server = "irc.freenode.net:7000"
	cfg.NewNick = func(n string) string { return n + "^" }
	cfg.Version = "0.1"
	cfg.QuitMessage = "Oh god, I am so depressed!"

	c := irc.Client(cfg)
	b := &Bot{conn: c}
	b.quit = make(chan bool)

	c.HandleFunc(irc.CONNECTED, b.onConncted)
	c.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) { b.quit <- true })

	c.HandleFunc(irc.PRIVMSG, b.onPrivMessage)
	c.HandleFunc(irc.JOIN, b.onJoin)

	return b
}

func main() {

	b := NewBot()

	if err := b.conn.Connect(); err != nil {
		fmt.Printf("Connection error: %s\n", err.Error())
	}
	rand.Seed(time.Now().UTC().UnixNano())
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	if secret != "" {
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

	<-b.quit
}

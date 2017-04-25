package main

import (
	"crypto/tls"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/G-Node/marvin/mensa"
	"github.com/G-Node/marvin/quotes"
	irc "github.com/fluffle/goirc/client"
	"os"
)

type Bot struct {
	conn *irc.Conn

	quit chan bool
}

func (b *Bot) onConncted(conn *irc.Conn, line *irc.Line) {os.Getenv("GITHUB_WEBHOOK_SECRET")
	channel := os.Getenv("MARVIN_CHANNEL")
	if channel ==""{
		channel = "#gnode"
	}
	conn.Join(channel)
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
		b.postMeals(conn, target, meals, []string{"Tagesgericht", "Aktionsessen", "Biogericht", "Aktion"})
		return
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
	messageFilter := strings.NewReplacer("[]", "")
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
		message := fmt.Sprintf("%s [%s] [%s]", meal.Name, notes,
			mensa.Emojify(strings.Join(prices, ", ")))
		conn.Privmsg(target, messageFilter.Replace(message))
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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if strings.Contains(a, b) {
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

	b.SetupGithub()

	<-b.quit
}

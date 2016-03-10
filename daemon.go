package main

import (
    "crypto/tls"
    "fmt"
    "strings"

    irc "github.com/fluffle/goirc/client"
    "github.com/g-node/marvin/mensa"    
)

type Bot struct {    
    conn *irc.Conn
    
    quit chan bool
}


func (b *Bot) onConncted(conn *irc.Conn, line *irc.Line) {
    conn.Join("#gnode");   
}

func (b *Bot) onPrivMessage(conn *irc.Conn, line *irc.Line) {
    text := line.Text()
    target := line.Target()
    
    fmt.Printf("[D] {pm}: [%s] %s\n", target, text)
    
    if strings.HasPrefix(text, "mensa") {
        mc := &mensa.Client{Address: "http://openmensa.org/api/v2"}
        
        var meals []mensa.Meal
        if (strings.Contains(text, "tomorrow")) {
            meals, _ = mc.MealsForTomorrow("134")
            // ignored error for now
        } else {
            meals, _ = mc.MealsForToday("134")
        }
        if len(meals) == 0 {
            conn.Privmsgf(target, "No milk today, my love has gone away...")
            return
        }
        
        for _, meal := range meals {
            category := meal.Category
            if !(strings.HasPrefix(category, "Tagesgericht") ||
                 strings.HasPrefix(category, "Aktionsessen")) {
                continue
            }
            var prices []string 
            for key, value := range meal.Prices {
                if value != 0.{
                    prices = append(prices,fmt.Sprintf("%s:%.2f€",
                                    key,value))
                }                
            }
            notes := mensa.Emojify(strings.Join(meal.Notes, ", "))
            conn.Privmsgf(target, "%s [%s] [%s]", meal.Name, notes, 
                          mensa.Emojify(strings.Join(prices, ", ")))
        }
    }
    if (strings.Contains(text, "nix")) {
            conn.Privmsg(target, "https://youtu.be/Go4SI5ie7qE")
        }
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
        
    return b
}


func main() {
 
    b := NewBot()
  
    if err := b.conn.Connect(); err != nil {
        fmt.Printf("Connection error: %s\n", err.Error())
    }
    
    <-b.quit
}

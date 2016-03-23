package main

import (
    "os"
    "fmt"
    "bytes"
    "strings"
    
    "github.com/gicmo/webhooks"
    "github.com/gicmo/webhooks/github"
)


func main() {
    hook := github.New(&github.Config{Secret: os.Getenv("GITHUB_WEBHOOK_SECRET")})
    hook.RegisterEvents(HandlePullRequest, github.PullRequestEvent)
    hook.RegisterEvents(HandleStatus, github.StatusEvent)

    port := os.Getenv("GITHUB_WEBHOOK_PORT")
    if port == "" {
        port = "2323"
    }
    err := webhooks.Run(hook, ":" + port, "/webhooks")

    if err != nil {
        fmt.Println(err)
    }
}

func HandlePullRequest(payload interface{}) {

    pl := payload.(github.PullRequestPayload)

    name := pl.PullRequest.Base.Repo.FullName 
    action := pl.Action
    number := pl.Number
    title := pl.PullRequest.Title
    from := pl.PullRequest.Head.Label
    to := pl.PullRequest.Base.Label
    sender := pl.Sender.Login
    cks := pl.PullRequest.Head.SHA
    
    fmt.Printf("[%s#%d] %.7s '%s' [%s → %s] %s (%s)\n", name, number, cks,
    title, from, to, action, sender)
      
}

func HandleStatus(payload interface{}) {
    pl := payload.(github.StatusPayload)
    
    state := pl.State
    
    if (state == "pending") {
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
    print(out.String())
}
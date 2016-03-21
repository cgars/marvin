package main

import (
    "os"
    "fmt"
    
    "github.com/gicmo/webhooks"
    "github.com/gicmo/webhooks/github"
)


func main() {
    hook := github.New(&github.Config{Secret: os.Getenv("GITHUB_WEBHOOK_SECRET")})
    hook.RegisterEvents(HandlePullRequest, github.PullRequestEvent)

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
    
    fmt.Printf("[%s#%d] '%s' [%s → %s] %s (%s)\n", name, number, 
    title, from, to, action, sender)
      
}
package main

import (
	"encoding/json"
	"fmt"
    "time"
	"io/ioutil"
	"net/http"
    irc "github.com/fluffle/goirc/client"

)
type Client struct {
    Address string
}






func main() {
    client := &Client{Address: "http://openmensa.org/api/v2"}
    res, err := client.Meals("134", time.Now())
    if err != nil {
        println(err)
    }    
    fmt.Printf("%v\n", res)
}
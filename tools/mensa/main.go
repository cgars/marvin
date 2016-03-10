package main

import (
	"fmt"
    "time"
	"github.com/g-node/marvin/mensa"

)

func main() {
    client := &mensa.Client{Address: "http://openmensa.org/api/v2"}
    res, err := client.Meals("134", time.Now())
    if err != nil {
        println(err)
    }    
    fmt.Printf("%v\n", res)
}
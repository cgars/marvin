package main

import (
	"fmt"
    "time"
	"github.com/g-node/marvin/mensa"
    "strings"
)

func main() {
    client := &mensa.Client{Address: "http://openmensa.org/api/v2"}
    meals, err := client.Meals("134", time.Now())
    if err != nil {
        println(err)
        return
    } 
    
    for _, meal := range meals {
        notes := mensa.Emojify(strings.Join(meal.Notes, ", "))
        fmt.Printf("%s [%s] (%s)\n", meal.Name, meal.Category, notes)
    }
    
}
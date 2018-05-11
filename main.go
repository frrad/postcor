package main

import (
	"log"

	"github.com/frrad/postcor/client"
	"github.com/frrad/settings"
)

func main() {
	set := client.Settings{
		Username: "your@email.com",
		Password: "secret12312",
	}
	x, err := settings.NewSettings(&set, []string{"~/.postcor"})
	if err != nil {
		log.Fatal(err)
	}

	c, err := client.NewClient(set, x.Save)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(c.GetPage("https://na.preva.com/exerciser-api//exerciser-account"))
}

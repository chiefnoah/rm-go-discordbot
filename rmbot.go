package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiefnoah/rm-go-discordbot/config"
	"log"
	"time"
)

func main() {

	cfg := config.LoadConfig()

	discord, err := discordgo.New(cfg.AppConfig.AuthToken)
	if err != nil {
		log.Fatal("Unable to authenticate with discord")
		return
	}
	err = discord.Open()
	if err != nil {
		log.Fatal("Unable to open connection: ", err)
		return
	}
	defer discord.Close()

	for {
		discord.UpdateStatus(0, "moderating")
		time.Sleep(10 * time.Second)
	}

}

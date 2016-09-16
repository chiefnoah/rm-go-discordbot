package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiefnoah/rm-go-discordbot/config"
	"log"
	"time"
	"github.com/chiefnoah/rm-go-discordbot/commands"
)

func main() {
	//Load config from file
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

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		log.Print("Author: ", m.Author.Username, "\n Content: ", m.Content)
		if len(m.Content) > 0 && string(m.Content[0]) == cfg.AppConfig.CommandPrefix {
			commands.ParseCommand(s, m.Message)
		}
	})

	for {
		discord.UpdateStatus(0, "moderating")
		time.Sleep(10 * time.Second)
	}

}

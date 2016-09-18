package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/chiefnoah/rm-go-discordbot/config"
	"log"
	"time"
	"github.com/chiefnoah/rm-go-discordbot/commands"
	"github.com/chiefnoah/rm-go-discordbot/database"
)

func main() {

	cfg := config.LoadConfig()
	database.Init()
	defer database.End()

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

	discord.AddHandler(func(s *discordgo.Session, m *discordgo.GuildMemberAdd) {
		newMessage := "Clan Overall Rules\nDon’t be a dick\nRespect and listen to the admins\nHave all the funs.\n" + cfg.AppConfig.CommandPrefix+"help for a list of commands for the management bot.\nRead the rest of the rules here: https://docs.google.com/document/d/1xi1nT6JGxpSwU-6YxcY8eX37MF4l3gzQbVnAqswrIKk/edit"
		message, err := s.ChannelMessageSend("225086487752867840", newMessage) //TODO: change this from a constant to something loaded from the DB
		if err != nil || message == nil {
			log.Print("Unable to send message to discord: ", err)
		}
	})

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

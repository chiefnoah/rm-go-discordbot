package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"log"
	"math/rand"
	"strconv"
	"github.com/chiefnoah/rm-go-discordbot/database"
	"github.com/chiefnoah/rm-go-discordbot/config"
	"encoding/base64"
	"io/ioutil"
)

/****************************************************************************************************************
*														*
*						Define commands here!						*
*														*
****************************************************************************************************************/
var helpCommand = CommandProcess{
	Triggers: map[string]interface{}{"h": nil, "help": nil, "wat": nil},
	Run: help,
	AdditionalParams: []string{},
	DeleteCommand: false,
	Description: "Prints this dialog",
}

var tempChannelCommand = CommandProcess{
	Triggers: map[string]interface{}{"gChannel": nil},
	Run: tempChannel,
	AdditionalParams: []string{},
	Description: "Generates a temporarty voice channel of given name",
	DeleteCommand: false,
}

var optRolesCommand = CommandProcess{
	Triggers: map[string]interface{}{"jR": nil, "joinRole": nil},
	Run: optIn,
	AdditionalParams: []string{},
	Description: "Allows Users to opt-in to hidden roles",
	DeleteCommand: false,
}

var d20Command = CommandProcess{
	Triggers: map[string]interface{}{"d20": nil, "roll": nil},
	Run: rollD20,
	AdditionalParams: []string{},
	Description: "Rolls a D20",
	DeleteCommand: false,
}

var setProfile = CommandProcess{
	Triggers: map[string]interface{}{"cpu": nil},
	Run: setGoddess,
	AdditionalParams: []string{},
	Description: "Sets the profile name and image of the bot",
	DeleteCommand: false,
}

var getRolesCommand = CommandProcess{
	Triggers: map[string]interface{}{"r": nil, "roles": nil},
	Run: roleInfo,
	AdditionalParams: []string{},
	Description: "Fetches roles and descriptors of roles",
	DeleteCommand: false,
}

//Commands MUST be specified here to be checked.
var enabledCommands []CommandProcess = []CommandProcess{tempChannelCommand, getRolesCommand, d20Command, optRolesCommand, setProfile}


//Wraps command triggers, additional parameters, and explicitly defines the function to be called when a command is typed
//Triggers: a map that's only used for the keys. Use contains() to check if the map contains
//Run: the function that is run whenever a command is typed
//AdditionalParams default parameters and other default parameters passed with the command. Allows for multiple commands with specific parameters to use the same function
//DeleteCommand: If true, after the bot processes a command it deletes the message command
type CommandProcess struct {
	Triggers         map[string]interface{}                                       //Maps for fast lookup, we don't actually care about what they hold
	Run              func(*discordgo.Session, *discordgo.Message, []string, bool) //I explicity define a function that implements Command so we can just loop through all the CommandProcesses and call Run generically
	AdditionalParams []string
	Description      string
	DeleteCommand    bool
}

func ParseCommand(s *discordgo.Session, c *discordgo.Message) {
	command := strings.Fields(c.Content)[0][1:]//takes the first word (has to be the command), and drops the prefix
	if contains(helpCommand.Triggers, command) {
		helpCommand.Run(s, c, []string{}, true)
		return
	}
	for _, v := range enabledCommands {
		if contains(v.Triggers, command) {
			v.Run(s, c, v.AdditionalParams, v.DeleteCommand)
		}
	}
}

func tempChannel(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	if len(m.Content) < 2 {
		return
	}
	channelName := strings.Split(m.Content, " ")[1:]
	//Add error checks

	newChannel := discordgo.Channel{}

	newChannel.Name = channelName[0]
	newChannel.Type = "voice"

	s.State.ChannelAdd(&newChannel)

	curChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Print("Unable to fetch channel")
		return
	}

	s.GuildMemberMove(curChannel.GuildID, m.Author.ID, newChannel.ID)

	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func optIn(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	if len(m.Content) < 2 {
		return
	}
	role := strings.Split(m.Content, " ")[1:][0]

	//IMPORTANT
	//TODO: ADD A CHECK TO SEE IF THEY CAN JOIN THE ROLE
	//IMPORTANT

	curChannel, err := s.Channel(m.ChannelID)
	if err != nil {
		log.Print("Unable to fetch channel")
		return
	}
	mem, err := s.GuildMember(curChannel.GuildID, m.Author.ID)
	if err != nil {
		log.Print("Unable to fetch guild member")
		return
	}
	s.GuildMemberEdit(curChannel.GuildID, m.Author.ID, append(mem.Roles, role))
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func roleInfo(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	args := strings.Split(m.Content, " ")[1:]
	if len(args) == 0 {
		getManagedRoles(s, m, extraArgs, deleteCommand)
		return
	}
	messageContent := "Fill in commands"

	message, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil || message == nil {
		log.Print("Unable to send message to discord: ", err)
	}
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func rollD20(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	messageContent := strconv.Itoa(rand.Intn(21))
	log.Print("Rolled: ", messageContent)
	message, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil || message == nil {
		log.Print("Unable to send message to discord: ", err)
	}
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func getManagedRoles(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	messageContent := "Roles available: "
	managedRoles, err := database.GetManagedRoles(m.Author.ID)
	if err != nil {
		messageContent = "No roles available to join :("
	}
	for _, mr := range managedRoles {
		messageContent += "\n*  " + mr.Name + " - " + mr.Description
	}

	message, err := s.ChannelMessageSend(m.ChannelID, messageContent)
	if err != nil || message == nil {
		log.Print("Unable to send message to discord: ", err)
	}
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func help(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	helpMessage := ""

	for _, v := range enabledCommands {
		for k, _ := range v.Triggers {
			helpMessage += k + ", "
		}
		helpMessage = helpMessage[:len(helpMessage) - 2] //Trim off the last ", "
		helpMessage += v.Description + "\n"
	}

	message, err := s.ChannelMessageSend(m.ChannelID, helpMessage)
	if err != nil || message == nil {
		log.Print("Unable to send message to discord: ", err)
	}
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func setGoddess(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	args := strings.Split(m.Content, " ")[1:]
	if len(args) < 1 {
		//we could return an error message, or just delete the command and be done with it idk

		return
	}
	cfg := config.LoadConfig()
	goddess := strings.Split(m.Content, " ")[1:][0]
	log.Print("Attempting to set profile to goddess: ", goddess)
	botUser, err := s.User("@me")
	if err != nil {
		log.Print("Unable to get bot user!")
		return
	}
	//log.Printf("Guild: %+s\nBot: %+s", curGuild, botUser)
	name := ""
	avatarImgDate := "data:image/jpeg;base64,"
	for p, v := range cfg.CPUConfig.GoddessNames {
		if (v == goddess) {
			rnd := rand.Intn(len(cfg.CPUConfig.GoddessImages[p]))
			data, err := ioutil.ReadFile("./images/" + cfg.CPUConfig.GoddessImages[p][rnd])
			if err != nil {
				log.Print("Unable to load profile image for goddess: ", err)
				return
			}
			avatarImgDate += base64.StdEncoding.EncodeToString(data)
			name = cfg.CPUConfig.GoddessNames[p]
		}
	}
	//This means we got were able to load a goddess name and image
	if len(name) > 0 {
		guilds, err := s.UserGuilds()
		if err != nil {
			log.Print("Unable to retrieve the guilds the user is in: ",err)
		}
		for _, g := range guilds {

			newUser, err := s.UserUpdate("", "", botUser.Username, avatarImgDate, "")
			if err != nil {
				log.Print("Unable to update user: ", err)
				return
			}
			err = s.GuildMemberNickname(g.ID, newUser.ID, name + "-Bot")
			if err != nil {
				log.Print("Unable to set nickname: ", err)
			}
			log.Print("Updated user: ", newUser.Username)
		}
	} else {
		log.Print("Unable to find set goddess: ", args)
	}
}

func contains(set map[string]interface{}, s string) bool {
	_, ok := set[s]
	return ok
}
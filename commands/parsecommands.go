package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"log"
	"math/rand"
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
}

var tempChannelCommand = CommandProcess{
	Triggers: map[string]interface{}{"gChannel": nil},
	Run: tempChannel,
	AdditionalParams: []string{},
	Description: "Generates a temporarty voice channel of given name",
	DeleteCommand: false,
}

var optRoles = CommandProcess{
	Triggers: map[string]interface{}{"jR": nil, "joinRole": nil},
	Run: optIn,
	AdditionalParams: []string{},
	Description: "Allows Users to opt-in to hidden roles",
	DeleteCommand: false,
}

var d20 = CommandProcess{
	Triggers: map[string]interface{}{"d20": nil, "roll": nil},
	Run: rollD20,
	AdditionalParams: []string{},
	Description: "Rolls a D20",
	DeleteCommand: false,
}

var getRoles = CommandProcess{
	Triggers: map[string]interface{}{"r": nil, "roles": nil},
	Run: roleInfo,
	AdditionalParams: []string{},
	Description: "Fetches roles and descriptors of roles",
	DeleteCommand: false,
}

//Commands MUST be specified here to be checked.
var enabledCommands []CommandProcess = []CommandProcess{helpCommand, tempChannelCommand, getRoles, d20, optRoles}


//Wraps command triggers, additional parameters, and explicitly defines the function to be called when a command is typed
//Triggers: a map that's only used for the keys. Use contains() to check if the map contains
//Run: the function that is run whenever a command is typed
//AdditionalParams default parameters and other default parameters passed with the command. Allows for multiple commands with specific parameters to use the same function
//DeleteCommand: If true, after the bot processes a command it deletes the message command
type CommandProcess struct {
	Triggers         map[string]interface{}                                       //Maps for fast lookup, we don't actually care about what they hold
	Run              func(*discordgo.Session, *discordgo.Message, []string, bool) //I explicity define a function that implements Command so we can just loop through all the CommandProcesses and call Run generically
	AdditionalParams []string
	Description	 string
	DeleteCommand    bool
}

func ParseCommand(s *discordgo.Session, c *discordgo.Message) {
	command := strings.Fields(c.Content)[0][1:]//takes the first word (has to be the command), and drops the prefix
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

	newChannel := discordgo.ChannelCreate{}

	newChannel.Name = channelName
	newChannel.Type = "voice"

	s.State.ChannelAdd(newChannel)

	s.GuildMemberMove(nil,m.Author,newChannel.ID)

	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func optIn(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	if len(m.Content) < 2 {
		return
	}
	role := strings.Split(m.Content, " ")[1:]

	//IMPORTANT
	//TODO: ADD A CHECK TO SEE IF THEY CAN JOIN THE ROLE
	//IMPORTANT

	curChannel, err1 := s.Channel(m.ChannelID)
	if err1 != nil {
		log.Fatal("Unable to fetch channel")
		return
	}
	mem, err2 := s.GuildMember(curChannel.GuildID, m.Author.ID)
	if err2 != nil {
		log.Fatal("Unable to fetch guild member")
		return
	}
	s.GuildMemberEdit(curChannel.GuildID,m.Author.ID,append(mem.Roles,role))
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func roleInfo(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
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
	messageContent := "" + rand.Intn(21)

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
		for _, b := range enabledCommands[v].Triggers {
			append(helpMessage, enabledCommands[v].Triggers[b] + ", " )
		}
		append(helpMessage, enabledCommands[v].Description + "\n")
	}

	message, err := s.ChannelMessageSend(m.ChannelID, helpMessage)
	if err != nil || message == nil {
		log.Print("Unable to send message to discord: ", err)
	}
	if deleteCommand {
		s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func contains(set map[string]interface{}, s string) bool {
	_, ok := set[s]
	return ok
}
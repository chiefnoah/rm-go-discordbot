package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
	"log"
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
	Triggers: map[string]interface{}{"temp": nil, "channel": nil},
	Run: tempChannel,
	AdditionalParams: []string{},
	DeleteCommand: false,
}

var getRoles = CommandProcess{
	Triggers: map[string]interface{}{"r": nil, "roles": nil},
	Run: roleInfo,
	AdditionalParams: []string{},
	DeleteCommand: false,
}

//Commands MUST be specified here to be checked.
var enabledCommands []CommandProcess = []CommandProcess{helpCommand, tempChannelCommand, getRoles}

//Array of tempChannels Created
var genChannels []discordgo.Channel
var numGenChannels int = 0

//Wraps command triggers, additional parameters, and explicitly defines the function to be called when a command is typed
//Triggers: a map that's only used for the keys. Use contains() to check if the map contains
//Run: the function that is run whenever a command is typed
//AdditionalParams default parameters and other default parameters passed with the command. Allows for multiple commands with specific parameters to use the same function
//DeleteCommand: If true, after the bot processes a command it deletes the message command
type CommandProcess struct {
	Triggers         map[string]interface{}                                       //Maps for fast lookup, we don't actually care about what they hold
	Run              func(*discordgo.Session, *discordgo.Message, []string, bool) //I explicity define a function that implements Command so we can just loop through all the CommandProcesses and call Run generically
	AdditionalParams []string
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

	if(numGenChannels + 1 >= len(genChannels)){
		// genChannels is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newGenChannels := make([]int, len(genChannels), 2*len(genChannels)+1)
		copy(newGenChannels, genChannels)
		genChannels = newGenChannels
	}
	genChannels[numGenChannels] = discordgo.ChannelCreate{}

	genChannels[numGenChannels].Name = channelName
	genChannels[numGenChannels].Type = "Voice"

	s.State.ChannelAdd(genChannels[numGenChannels])

	s.GuildMemberMove(nil,m.Author,genChannels[numGenChannels].ID)
	numGenChannels++
}

func optIn(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {

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

func help(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {
	messageContent := "Fill in commands"

	message, err := s.ChannelMessageSend(m.ChannelID, messageContent)
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
package commands

import (
	"github.com/bwmarrin/discordgo"
	"strings"
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



//Commands MUST be specified here to be checked.
var enabledCommands []CommandProcess = []CommandProcess{helpCommand}

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

func help(s *discordgo.Session, m *discordgo.Message, extraArgs []string, deleteCommand bool) {

}

func contains(set map[string]interface{}, s string) bool {
	_, ok := set[s]
	return ok
}
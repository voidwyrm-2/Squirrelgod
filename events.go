package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	msgLow := strings.ToLower(m.Content)
	msgClean := strings.TrimSpace(m.Content)

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	} else if ContainsAnyOf(msgLow, "squirrel") &&
		!strings.Contains(msgLow, "squirrelgod") {
		sendReply(s, m.Reference(), "https://tenor.com/view/squirrel-rotating-ring-gif-7442174")
		return
	} else if strings.Contains(msgLow, "brazil") {
		sendReply(s, m.Reference(), "BRAZIL MENTION!!!!")
		return
	} else if ContainsAnyOf(msgLow, "neuro-sama", "vedal") {
		sendReply(s, m.Reference(), "Someone tell Vedal there is a problem with my AI.")
		return
	}

	if !strings.HasPrefix(msgClean, "$") {
		return
	}

	commandHandler(s, m.Reference(), msgClean)
}

func commandHandler(s *discordgo.Session, ref *discordgo.MessageReference, msgClean string) error {
	args := strings.Split(msgClean[1:], " ")
	comm := strings.ToLower(strings.TrimSpace(args[0]))
	args = args[1:]

	if comm == "help" {
		return sqHelp(s, encapSendReply(s, ref), args)
	}

	if command, ok := commands[comm]; ok {
		return command.run(s, encapSendReply(s, ref), args)
	}

	return nil
}

//func perceptionCheck(s *discordgo.Session, m *discordgo.Message) {
//}

func messageReact(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	fmt.Println(m.Emoji.ID, m.UserID)
}

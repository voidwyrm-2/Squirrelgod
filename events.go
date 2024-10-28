package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	} else if ContainsAnyOf(strings.ToLower(m.Content), "squirrel") &&
		!strings.Contains(strings.ToLower(m.Content), "squirrelgod") {
		err := SendReply(s, m, "https://tenor.com/view/squirrel-rotating-ring-gif-7442174")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if strings.Contains(strings.ToLower(m.Content), "brazil") {
		err := SendReply(s, m, "BRAZIL MENTION!!!!")
		if err != nil {
			fmt.Println(err.Error())
		}
	} else if ContainsAnyOf(strings.ToLower(m.Content), "neuro-sama", "vedal") {
		err := SendReply(s, m, "Someone tell Vedal there is a problem with my AI.")
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func messageReact(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	fmt.Println(m.Emoji.ID, m.UserID)
}

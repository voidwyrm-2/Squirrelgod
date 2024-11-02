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

	end, err := perceptionCheck(s, strings.TrimSpace(msgLow), m.Author.ID, m.Reference())
	if err != nil {
		fmt.Println(err.Error())
		return
	} else if end {
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

var (
	pCheck = false
	p1     = ""
)

func perceptionCheck(s *discordgo.Session, content, authID string, ref *discordgo.MessageReference) (bool, error) {
	if content == "$pcheck" {
		pCheck = true
		return true, SendReplyErr(s, ref, `I feel it in my fingers, I feel it in my toes
These people mean to harm us and they got to go
So c'mon get 'em now
You picked the wrong day to mess around with my tight crew, ooh, ooh
There's no escaping it, I can perceive you
Here's what we're gonna do
Me and my boys are gonna mess you up`)
	} else if content == "i rolled a one" {
		if p1 != "" {
			if p1 == authID {
				return true, nil
			}
			p1 = ""
			return true, SendReplyErr(s, ref, "CRAP\nMy boys are otherwise engaged\nSo I'm gonna bring it all myself- hey")
		} else {
			p1 = authID
		}
	}

	return false, nil
}

func messageReactAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	if r.UserID == s.State.User.ID {
		return
	}
	// fmt.Printf("given: '%s'(%v)\n", r.Emoji.Name, isValidOffering(r.Emoji.Name))

	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		fmt.Println("error from messageReactAdd:", err.Error())
		return
	}

	if isValidOffering(r.Emoji.Name) && msg.Author.ID == s.State.User.ID {
		// fmt.Println("offering given")
		offeringCount++
	}
}

func messageReactRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	if r.UserID == s.State.User.ID {
		return
	}
	// fmt.Printf("removed: '%s'(%v)"\n, r.Emoji.Name, isValidOffering(r.Emoji.Name))

	msg, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		fmt.Println("error from messageReactRemove:", err.Error())
		return
	}

	if isValidOffering(r.Emoji.Name) && msg.Author.ID == s.State.User.ID {
		// fmt.Println("offering removed")
		if offeringCount > 0 {
			offeringCount--
		}
	}
}

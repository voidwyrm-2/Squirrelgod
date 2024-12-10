package main

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type sqCommand struct {
	run         func(*discordgo.Session, *discordgo.Channel, func(msg string), []string) error
	description string
}

func newSQCommand(description string, run func(s *discordgo.Session, ch *discordgo.Channel, reply func(msg string), _ []string) error) sqCommand {
	return sqCommand{description: description, run: run}
}

var sqHelp = func(_ *discordgo.Session, reply func(msg string), _ []string) error {
	commandList := []string{"`help`: Lists all commands"}
	for name, c := range commands {
		commandList = append(commandList, fmt.Sprintf("`%s`: %s", name, c.description))
	}

	reply("**Squirrelgod commands:**\n" + strings.Join(commandList, "\n"))

	return nil
}

var commands = map[string]sqCommand{
	"source": newSQCommand("Replies with the link to the source repo", func(_ *discordgo.Session, ch *discordgo.Channel, reply func(msg string), _ []string) error {
		if sourceLink != "" {
			reply("<" + sourceLink + ">")
		} else {
			reply("no source repo link was given")
		}
		return nil
	}),
	"echo": newSQCommand("A simple response test", func(s *discordgo.Session, ch *discordgo.Channel, reply func(msg string), args []string) error {
		reply(strings.Join(args, " "))
		return nil
	}),
	"offerings": newSQCommand("Replies with the amount of offerings given", func(s *discordgo.Session, ch *discordgo.Channel, reply func(msg string), _ []string) error {
		if offeringCount == 69 || offeringCount == 42 {
			reply(fmt.Sprintf("Offerings given so far: %v\n(nice)", offeringCount))
		} else {
			reply(fmt.Sprintf("Offerings given so far: %v", offeringCount))
		}
		return nil
	}),
	"onmsgs": newSQCommand("Replies with the size of the pool of messages that are randomly selected to be said when Squirrelgod comes online", func(s *discordgo.Session, ch *discordgo.Channel, reply func(msg string), _ []string) error {
		reply(fmt.Sprint(len(onlineAnnounceMessages)))
		return nil
	}),
	"origin": newSQCommand("Shows the origin of the messages that Squirrelgod said when he came online; only works in specific circumstances", func(s *discordgo.Session, ch *discordgo.Channel, reply func(msg string), _ []string) error {
		if slices.Contains(channelsThatCanShowOrigins, ch.ID) {
			return nil
		}

		if origin, ok := onlineAnnounceMessageOrigins[usedAnnounceMessage]; ok {
			reply(origin)
		} else {
			reply("no origin found")
		}
		return nil
	}),
}

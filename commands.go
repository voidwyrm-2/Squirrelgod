package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type sqCommand struct {
	run         func(*discordgo.Session, func(msg string), []string) error
	description string
}

func newSQCommand(description string, run func(s *discordgo.Session, reply func(msg string), _ []string) error) sqCommand {
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
	"source": newSQCommand("Replies with the link to the source repo", func(s *discordgo.Session, reply func(msg string), _ []string) error {
		if sourceLink != "" {
			reply("<" + sourceLink + ">")
		} else {
			reply("no source repo link was given")
		}
		return nil
	}),
}

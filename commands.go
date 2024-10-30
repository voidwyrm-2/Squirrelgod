package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type sqCommand struct {
	run         func(s *discordgo.Session, reply func(msg string), _ []string) error
	description string
}

func newSQCommand(description string, run func(s *discordgo.Session, reply func(msg string), _ []string) error) sqCommand {
	return sqCommand{description: description, run: run}
}

var sqHelp = func(s *discordgo.Session, reply func(msg string), _ []string) error {
	commandList := []string{"`help`: Lists all commands"}
	for name, c := range commands {
		commandList = append(commandList, fmt.Sprintf("`%s`: %s", name, c.description))
	}

	reply("**Squirrelgod commands:**\n" + strings.Join(commandList, "\n"))

	return nil
}

var commands = map[string]sqCommand{}

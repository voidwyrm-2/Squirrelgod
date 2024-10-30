package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func sendReply(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	err := SendReplyErr(s, m, msg)
	fmt.Println(err.Error())
}

func SendReplyErr(s *discordgo.Session, m *discordgo.MessageCreate, msg string) error {
	_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m.Reference())
	return err
}

func ContainsAnyOf(str string, substrs ...string) bool {
	if strings.Contains(str, substrs[0]) {
		return true
	}

	if len(substrs) < 2 {
		return false
	}
	return ContainsAnyOf(str, (substrs[1:])...)
}

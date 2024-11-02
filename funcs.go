package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func readFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	content := ""
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return content, nil
}

func writeFile(filename string, data string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func encapSendReply(s *discordgo.Session, m *discordgo.MessageReference) func(msg string) {
	return func(msg string) {
		sendReply(s, m, msg)
	}
}

func sendReply(s *discordgo.Session, m *discordgo.MessageReference, msg string) {
	err := SendReplyErr(s, m, msg)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func SendReplyErr(s *discordgo.Session, m *discordgo.MessageReference, msg string) error {
	_, err := s.ChannelMessageSendReply(m.ChannelID, msg, m)
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

func isValidOffering(offering string) bool {
	return offering == "🥜" || offering == "🌰"
}

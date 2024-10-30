package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/voidwyrm-2/goconf"
)

var (
	botToken                      = ""
	installLink                   = ""
	onlineAnnounceChannel         = ""
	onlineAnnounceMessages        = []string{}
	offeringCount          uint64 = 0
)

func sgInit() error {
	conf, err := goconf.Load("config.txt")
	if err != nil {
		return err
	}

	if tok, ok := conf["token"]; !ok {
		return errors.New("missing key 'token' in config")
	} else {
		botToken = strings.TrimSpace(tok.(string))
	}

	if insLink, ok := conf["install_link"]; ok {
		installLink = strings.TrimSpace(insLink.(string))
	}

	if anncChan, ok := conf["announce_channel"]; ok {
		onlineAnnounceChannel = strings.TrimSpace(anncChan.(string))
	}

	if anncMsgs, ok := conf["online_messages"]; ok {
		onlineAnnounceMessages = strings.Split(strings.TrimSpace(anncMsgs.(string)), "\n")
	}

	if ofco, err := readFile("offeringCount.txt"); err != nil {
		return err
	} else {
		n, err := strconv.ParseUint(strings.TrimSpace(ofco), 10, 0)
		if err != nil {
			return err
		}
		offeringCount = n
		fmt.Printf("loaded offeringCount as '%v'\n", offeringCount)
	}

	return nil
}

func sgExit() {
	if err := writeFile("offeringCount.txt", fmt.Sprintf("%v", offeringCount)); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	fmt.Println("===INIT===")
	err := sgInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer sgExit()

	// Create a new Discord session using the provided bot token.
	fmt.Println("creating Discord session...")
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// register event handlers
	fmt.Println("registering event handlers...")
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReact)

	fmt.Println("adding intents...")
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	fmt.Println("opening websocket...")
	atStart := time.Now()
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	atEnd := time.Now()
	fmt.Printf("opened web socket in %s\n", atEnd.Sub(atStart).String())

	fmt.Println("sending start-up message...")
	if onlineAnnounceChannel != "" {
		msg := onlineAnnounceMessages[rand.Intn(len(onlineAnnounceMessages))]
		fmt.Println("message rolled, got `" + msg + "`")
		_, err := dg.ChannelMessageSend(onlineAnnounceChannel, msg)
		if err != nil {
			fmt.Println("error while sending start-up message: " + err.Error())
			return
		}
	} else {
		fmt.Println("could not send start-up message, channel not given")
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("init finished!")
	fmt.Println("===END INIT===")

	fmt.Println("signal loop:")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt) //, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	fmt.Println()
	fmt.Println("ending session...")
	dg.Close()
	fmt.Println("session ended")
}

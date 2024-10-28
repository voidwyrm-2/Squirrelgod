package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/voidwyrm-2/goconf"
)

var botToken = ""

func sgInit() error {
	conf, err := goconf.Load("config.txt")
	if err != nil {
		return err
	}

	if tok, ok := conf["token"]; !ok {
		return errors.New("missing key 'token' in config")
	} else {
		botToken = tok.(string)
	}

	return nil
}

func main() {
	fmt.Println("===INIT===")
	err := sgInit()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

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

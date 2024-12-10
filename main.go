package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/bwmarrin/discordgo"
	"github.com/voidwyrm-2/goconf"
)

var (
	botToken                            = ""
	installLink                         = ""
	sourceLink                          = ""
	onlineAnnounceChannel               = ""
	onlineAnnounceMessages              = []string{}
	usedAnnounceMessage                 = ""
	onlineAnnounceMessageOrigins        = map[string]string{}
	channelsThatCanShowOrigins          = []string{}
	offeringCount                uint64 = 0
)

func sgInit() error {
	fmt.Println("reading config and stat data...")
	defer fmt.Println("finished reading config and stat data")

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

	if sLink, ok := conf["source_link"]; ok {
		sourceLink = strings.TrimSpace(sLink.(string))
	}

	if anncChan, ok := conf["announce_channel"]; ok {
		onlineAnnounceChannel = strings.TrimSpace(anncChan.(string))
	}

	if showOChan, ok := conf["channels_show_origins"]; ok {
		channelsThatCanShowOrigins = strings.Split(strings.TrimSpace(showOChan.(string)), "\n")
	}

	/*
		if anncMsgs, ok := conf["online_messages"]; ok {
			getmo := func(msg string) (string, string) {
				split := strings.Split(msg, "<origin>")
				if len(split) == 1 {
					split = append(split, "")
				}
				return strings.TrimSpace(split[0]), strings.TrimSpace(split[1])
			}

			for _, msg := range strings.Split(strings.TrimSpace(anncMsgs.(string)), "\n") {
				m, o := getmo(msg)
				onlineAnnounceMessages = append(onlineAnnounceMessages, o)
				if o != "" {
					onlineAnnounceMessageOrigins[m] = o
				}
			}
		}
	*/

	var announceMsgInfo struct {
		Messages [][]string
	}

	if content, err := readFile("announce_msgs.toml"); err != nil {
		return err
	} else if _, err = toml.Decode(content, &announceMsgInfo); err != nil {
		return err
	}

	for _, m := range announceMsgInfo.Messages {
		onlineAnnounceMessages = append(onlineAnnounceMessages, m[0])
		if len(m) > 1 {
			if m[1] != "" {
				onlineAnnounceMessageOrigins[m[0]] = m[1]
			}
		}
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

	if len(onlineAnnounceMessages) == 0 {
		onlineAnnounceMessages = append(onlineAnnounceMessages, "[PLACEHOLDER]")
	}

	return nil
}

func sgExit() {
	fmt.Println("saving data...")
	defer fmt.Println("data saved")
	if err := writeFile("offeringCount.txt", fmt.Sprintf("%v", offeringCount)); err != nil {
		fmt.Println(err.Error())
	}
}

func main() {
	fmt.Println("===INIT===")
	if err := sgInit(); err != nil {
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

	fmt.Println("adding intents...")
	dg.Identify.Intents = (discordgo.IntentsGuildMessages | discordgo.IntentGuildMessageReactions)

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

	// register event handlers
	fmt.Println("registering event handlers...")
	dg.AddHandler(messageCreate)
	dg.AddHandler(messageReactAdd)
	dg.AddHandler(messageReactRemove)
	fmt.Println("event handlers registered")

	if slices.Contains(os.Args, "nm") {
		fmt.Println("start-up message ignored because of a passed flag")
	} else {
		fmt.Println("sending start-up message...")
		if onlineAnnounceChannel != "" {
			msg := onlineAnnounceMessages[rand.Intn(len(onlineAnnounceMessages))]
			usedAnnounceMessage = msg
			fmt.Println("message rolled, got `" + usedAnnounceMessage + "`")
			_, err := dg.ChannelMessageSend(onlineAnnounceChannel, msg)
			if err != nil {
				fmt.Println("error while sending start-up message: " + err.Error())
				return
			}
		} else {
			fmt.Println("could not send start-up message, channel not given")
		}
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("init finished")
	fmt.Println("===END INIT===")

	fmt.Println("==RUNTIME LOGS==")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt) //, os.Kill)
	<-sc

	fmt.Println("==END RUNTIME LOGS==")

	// Cleanly close down the Discord session.
	fmt.Println("ending session...")
	dg.Close()
	fmt.Println("session ended")
}

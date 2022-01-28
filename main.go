package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const Version = "v0.0.0-alpha"

//TODO: move bot instantiation into this as a wrapper that can be called in init. issue rn is I need dg to be referenceable outside of init, otherwise I would spawn it there. current issue is idk the type of the dg bot. realllly should, but its a pointer to a Session object defined somewhere deep in the lib.
func botGen() {

}

func init() {
	// Discord Authentication Token
	// Print out a fancy logo!
	fmt.Printf(`Science Defender %-16s\/`+"\n\n", Version)

	//Load dotenv file from .
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//Bot session instance
	var dg, err = discordgo.New()

	//Load Token from env (simulated with godotenv)
	dg.Token = os.Getenv("BOT_TOKEN")
	if dg.Token == "" {
		log.Fatal("Error loading token from env file")
		os.Exit(1)
	}

	//Add Event Handler Functions
	dg.AddHandler(messageCreateHandler) //use for message create events

	//Register Bot Intents with Discord
	//worth noting MakeIntent is a no-op, but I want it there for doing something with pointers later
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Clean up
	dg.Close()

	// Exit Normally.
	//exit

}

//<----------------------> HANDLERS <---------------------->

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	// if m.Author.ID == s.State.User.ID {
	// 	return
	// }
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}

}

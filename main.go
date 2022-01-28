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

func main() {
	//Declarations
	var err error

	// Print out a fancy logo!
	fmt.Printf(`Science Defender %-16s\/`+"\n\n", Version)

	//Load dotenv file from .
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//create session
	var Session, _ = discordgo.New()
	Session.Token = ""

	//Load Token from env (simulated with godotenv)
	Session.Token = os.Getenv("BOT_TOKEN")
	if Session.Token == "" {
		log.Fatal("Error loading token from env file")
		return
	}

	// Open a websocket connection to Discord
	err = Session.Open()
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
	Session.Close()

	// Exit Normally.
	//exit

}

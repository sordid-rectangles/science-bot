package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/Necroforger/dgrouter"

	"github.com/Necroforger/dgrouter/exrouter"
)

const Version = "v0.0.0-alpha"

//var dg *discordgo.Session
var BOTID string
var PREFIX string
var TOKEN string
var FITE string
var ADMIN string

func init() {
	// Print out a fancy logo!
	fmt.Printf(`Science Defender! %-16s\/`+"\n\n", Version)

	//Load dotenv file from .
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Load Token from env (simulated with godotenv)
	TOKEN = os.Getenv("BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Error loading token from env file")
		os.Exit(1)
	}

	//Set some constants
	//TODO: move this into the .env and other more cleaned up config areas
	PREFIX = "?"
	FITE = ""
	ADMIN = "444543604547911680"
}

func main() {

	//Instantiation of core tools

	//Configure discordgo session bot
	var dg, err = discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("Error creating discordgo session!")
		os.Exit(1)
	}

	//configure dgrouter instance
	var router = exrouter.New()

	//Add regex routes to router
	// Add some commands
	router.On("ping", func(ctx *exrouter.Context) {
		ctx.Reply("pong")
	}).Desc("responds with pong")

	router.On("avatar", func(ctx *exrouter.Context) {
		ctx.Reply(ctx.Msg.Author.AvatarURL("2048"))
	}).Desc("returns the user's avatar")

	// Match the regular expression user(name)?
	router.OnMatch("username", dgrouter.NewRegexMatcher("user(name)?"), func(ctx *exrouter.Context) {
		ctx.Reply("Your username is " + ctx.Msg.Author.Username)
	})

	// Match the regular expression user(name)?
	router.OnMatch("fite", dgrouter.NewRegexMatcher(`(fite)+:[\w]+#\d\d\d\d?`), func(ctx *exrouter.Context) {
		ctx.Reply("Hit") //print debug
		ctx.Reply(ctx.Msg.Author.ID)
		if ctx.Msg.Author.ID == string(ADMIN) {
			FITE = strings.Split(ctx.Msg.Content, ":")[1]
			ctx.Reply("Now configured to argue with: " + FITE)
		}
	}).Desc("Configures who this bot will fight with. Must be a greenlisted user to use")

	router.Default = router.On("help", func(ctx *exrouter.Context) {
		var text = ""
		for _, v := range router.Routes {
			text += v.Name + " : \t" + v.Description + "\n"
		}
		ctx.Reply("```" + text + "```")
	}).Desc("prints this help menu")

	//Add Event Handler Functions
	dg.AddHandler(messageCreateHandler) //use for message create events
	dg.AddHandler(func(_ *discordgo.Session, m *discordgo.MessageCreate) {
		router.FindAndExecute(dg, PREFIX, dg.State.User.ID, m.Message)
	})

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
	// if m.Content == "ping" {
	// 	s.ChannelMessageSend(m.ChannelID, "Pong!")
	// }

	// // If the message is "pong" reply with "Ping!"
	// if m.Content == "pong" {
	// 	s.ChannelMessageSend(m.ChannelID, "Ping!")
	// }

}

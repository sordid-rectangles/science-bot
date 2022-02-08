package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"

	"github.com/sordid-rectangles/science-bot/auth"
	"github.com/sordid-rectangles/science-bot/responses"
)

type FiteUser struct {
	_name string
	_mode int
}

const Version = "v0.0.0-alpha"

var dg *discordgo.Session
var BOTID string
var PREFIX string
var TOKEN string
var FITE map[string]FiteUser = make(map[string]FiteUser)
var ADMINS map[string]string = make(map[string]string)
var OWNER map[string]string = make(map[string]string)

var (
	GuildID        string
	BotToken       string
	RemoveCommands bool
)

func init() {
	// Print out a fancy logo!
	fmt.Printf(`Science Defender! %-16s\/`+"\n\n", Version)

	//init maps
	auth.SetMaps(&ADMINS, &OWNER)

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

	OWNER_ID := os.Getenv("OWNER_ID")
	if OWNER_ID == "" {
		log.Fatal("Error loading admin id from env file")
		os.Exit(1)
	}
	OWNER_NAME := os.Getenv("OWNER_NAME")
	if OWNER_NAME == "" {
		log.Fatal("Error loading admin id from env file")
		os.Exit(1)
	}

	PREFIX = os.Getenv("CMD_PREFIX")
	if PREFIX == "" {
		log.Fatal("Error loading admin id from env file")
		os.Exit(1)
	}

	BotToken = TOKEN
	RemoveCommands = true

	//Set some constants
	//TODO: move this into the .env and other more cleaned up config areas
	for k := range auth.Users.OWNER {
		delete(auth.Users.OWNER, k)
	}

	auth.SetOwner(OWNER_ID, OWNER_NAME)
	// //Configure discordgo session bot
	// var dg, err = discordgo.New("Bot " + TOKEN)
	// if err != nil {
	// 	log.Fatal("Error creating discordgo session!")
	// 	os.Exit(1)
	// }
	fmt.Printf("Owners: %v \n", auth.Users.OWNER)
}

var (
	//<------------------->Command Object Creation<------------------->
	commands = []*discordgo.ApplicationCommand{
		{
			Name: "whoami",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Returns the discord id of the user who called the command",
		},
		{
			Name: "admins",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Returns the list of the users registered as admins",
		},
		{
			Name:        "register-admin",
			Description: "Register a new bot admin",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-select",
					Description: "select a user to register",
					Required:    true,
				},
			},
		},
		{
			Name:        "remove-admin",
			Description: "Remove an existing bot admin",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-select",
					Description: "select a user to remove",
					Required:    true,
				},
			},
		},
		{
			Name:        "register-fite-user",
			Description: "Register a new user for the bot to fite with",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-select",
					Description: "select a user to register",
					Required:    true,
				},
				{
					Name:        "fite-type",
					Description: "Select the fite mode",
					Type:        discordgo.ApplicationCommandOptionInteger,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "simple",
							Value: 0,
						},
						{
							Name:  "censor",
							Value: 1,
						},
						{
							Name:  "nlp",
							Value: 2,
						},
						{
							Name:  "simple-censor",
							Value: 3,
						},
						{
							Name:  "nlp-censor",
							Value: 4,
						},
					},
					Required: true,
				},
			},
		},
		{
			Name:        "remove-fite-user",
			Description: "Remove an existing registered fite user",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-select",
					Description: "select a user to remove",
					Required:    true,
				},
			},
		},
	}

	//<------------------->Command Handler Function Creation<------------------->
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"whoami": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			id := i.Member.User.ID
			content := "Your discord user id is: " + string(id)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"admins": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string
			if len(auth.Users.ADMINS) == 0 {
				content = "No Admins currently registered ;-("
			} else {
				content = "Current admin users are: \n"
				for id, name := range auth.Users.ADMINS {
					content += fmt.Sprintf("name: %s, ID: %s \n", name, id)
				}
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: content,
				},
			})
		},
		"register-admin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(s)

			content := `%s registered as an admin`

			auth.RegisterAdmin(user.ID, user.Username)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, user),
				},
			})
		},
		"remove-admin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(s)

			content := `%s removed as an admin`

			auth.RemoveAdmin(user.ID)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, user),
				},
			})
		},
		"register-fite-user": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			user := i.ApplicationCommandData().Options[0].UserValue(s)
			mode := i.ApplicationCommandData().Options[1].IntValue()

			_ = registerFite(user.ID, user.Username, int(mode))

			content := `%s registered as a target in mode %d `

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf(content, user, mode),
				},
			})
		},
	}
)

func init() {
	var err error
	dg, err = discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("Error creating discordgo session!")
		os.Exit(1)
	}
	// dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
	// 		h(s, i)
	// 	}
	// })
}

func main() {
	var err error
	//Configure discordgo session bot
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) { log.Println("Bot is up!") })
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	//Instantiation of core tools

	//Add Event Handler Functions
	dg.AddHandler(messageCreateHandler) //use for message create events

	//Register Bot Intents with Discord
	//worth noting MakeIntent is a no-op, but I want it there for doing something with pointers later
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	// Open a websocket connection to Discord
	err = dg.Open()
	if err != nil {
		log.Printf("error opening connection to Discord, %s\n", err)
		os.Exit(1)
	}

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)

	GuildID = ""
	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	// sc := make(chan os.Signal, 1)
	// signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	// <-sc

	// // Clean up
	// dg.Close()
	defer dg.Close()

	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Println("Gracefully shutdowning")

	// Exit Normally.
	//exit

}

//<----------------------> HANDLERS <---------------------->

func messageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	//println("hit")
	fiter, ok := FITE[string(m.Author.ID)]
	//print(ok)
	if ok {
		//println("hit ok")
		// s.ChannelMessageSendReply(m.ChannelID, "Ummmmmm actually that is incorrect, also don't care, also ratio", m.MessageReference)
		var r string
		var err error
		switch fiter._mode {

		//cases defined in the register fite user slash command
		case 0: //simple response mode
			r, err = responses.RandSimple()
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}
			s.ChannelMessageSendReply(m.ChannelID, r, m.Reference())
		case 1: //censor
			s.ChannelMessageDelete(m.ChannelID, m.Reference().MessageID)
		case 2: //nlp
			r, err = responses.GenProse(m.Content)
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}
			s.ChannelMessageSendReply(m.ChannelID, r, m.Reference())
		case 3: //simple-censor
			r, err = responses.RandSimple()
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}
			s.ChannelMessageSendReply(m.ChannelID, r, m.Reference())
			s.ChannelMessageDelete(m.ChannelID, m.Reference().MessageID)
		case 4: //nlp-censor
			r, err = responses.GenProse(m.Content)
			if err != nil {
				fmt.Println(err)
				log.Println(err)
			}
			s.ChannelMessageSendReply(m.ChannelID, r, m.Reference())
			s.ChannelMessageDelete(m.ChannelID, m.Reference().MessageID)
		}

		if err != nil {
			fmt.Println(err)
			log.Println(err)
		} else {
			//println("hit res")
			println(r)

			//println(mess)
			//println(err)
		}
	}
}

func registerFite(id string, name string, mode int) error {
	_, ok := FITE[id]
	if ok {
		FITE[id] = FiteUser{_name: name, _mode: mode}
		return nil
	} else {
		FITE[id] = FiteUser{_name: name, _mode: mode}
		//todo: also update the db
		return nil
	}
}

//Some misc example code I'm going to keep here for reference for now
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

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

const Version = "v0.0.1-alpha"

var dg *discordgo.Session
var BOTID string
var PREFIX string
var TOKEN string
var FITE map[string]FiteUser = make(map[string]FiteUser)
var ADMINS map[string]string = make(map[string]string)
var OWNER map[string]string = make(map[string]string)
var GUILDID string = "" //really only useful to guild-scope application commands, which is useful in dev to get them to update instantly rather than 1hr

func init() {
	// Print out a fancy logo!
	fmt.Printf(`Science Defender! %-16s\/`+"\n\n", Version)

	//init maps
	auth.SetMaps(&ADMINS, &OWNER)

	//Load dotenv file from .
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, trying env variables")
	}
	//Load Token from env (simulated with godotenv)
	TOKEN = os.Getenv("BOT_TOKEN")
	if TOKEN == "" {
		log.Fatal("Error loading token from env")
		os.Exit(1)
	}

	OWNER_ID := os.Getenv("OWNER_ID")
	if OWNER_ID == "" {
		log.Fatal("Error loading admin id from env")
		os.Exit(1)
	}
	OWNER_NAME := os.Getenv("OWNER_NAME")
	if OWNER_NAME == "" {
		log.Fatal("Error loading admin id from env")
		os.Exit(1)
	}

	// PREFIX = os.Getenv("CMD_PREFIX")
	// if PREFIX == "" {
	// 	log.Fatal("Error loading admin id from env file")
	// 	os.Exit(1)
	// }

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
		{
			Name:        "register-response",
			Description: "Register a new bot response",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "response-string",
					Description: "response string to be added to the simpleRes list",
					Required:    true,
				},
			},
		},
		{
			Name: "bot-responses",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Returns the list of strings registered as responses",
		},
		{
			Name: "fite-targets",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Returns the list of users registered as targets",
		},
	}

	//<------------------->Command Handler Function Creation<------------------->
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"whoami": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content := "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})

			} else {
				id := i.Member.User.ID
				content := "Your discord user id is: " + string(id)
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"admins": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					if len(auth.Users.ADMINS) == 0 {
						content = "No Admins currently registered ;-("
					} else {
						content = "Current admin users are: \n"
						for id, name := range auth.Users.ADMINS {
							content += fmt.Sprintf("name: %s, ID: %s \n", name, id)
						}
					}
				} else {
					content = `Peasant, you are not authorized.`
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"register-admin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content := "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {

				user := i.ApplicationCommandData().Options[0].UserValue(s)

				var content string
				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					content = user.Username + ` registered as an admin`
					auth.RegisterAdmin(user.ID, user.Username)
				} else {
					content = `Peasant, you are not authorized.`
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"remove-admin": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content := "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				user := i.ApplicationCommandData().Options[0].UserValue(s)
				var content string

				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					content = user.Username + `removed as an admin`
					auth.RemoveAdmin(user.ID)
				} else {
					content = `Peasant, you are not authorized.`
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"register-fite-user": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content := "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				user := i.ApplicationCommandData().Options[0].UserValue(s)
				mode := i.ApplicationCommandData().Options[1].IntValue()
				var content string

				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					_ = registerFite(user.ID, user.Username, int(mode))
					content = fmt.Sprintf(`%s registered as a target in mode %d`, user.Username, mode)
				} else {
					content = `Peasant, you are not authorized.`
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"remove-fite-user": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content := "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				user := i.ApplicationCommandData().Options[0].UserValue(s)
				var content string

				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					_ = removeFite(user.ID)
					content = `%s removed as a target`
				} else {
					content = `Peasant, you are not authorized.`
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content, user),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"register-response": func(s *discordgo.Session, i *discordgo.InteractionCreate) {

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content := "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				res := i.ApplicationCommandData().Options[0].StringValue()

				var content string

				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					responses.AddResponse(res)
					content = res + ` registered as a response`
				} else {
					content = `Peasant, you are not authorized.`
				}

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"bot-responses": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					if len(responses.SimpleRes) == 0 {
						content = "No responses currently registered ;-("
					} else {
						content = "Current responses are: \n"
						for i := range responses.SimpleRes {
							res := responses.SimpleRes[i]
							content += fmt.Sprintf("Response %d: %s \n", i, res)
						}
					}
				} else {
					content = `Peasant, you are not authorized.`
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

		},
		"fite-targets": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			var content string

			check, err := comesFromDM(s, i)
			if check {
				log.Println("Message in dm")

				content = "I can only be used in servers ;-("

				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: fmt.Sprintf(content),
					},
				})
			} else {
				a, _ := auth.IsAuthed(i.Member.User.ID)
				if a {
					if len(responses.SimpleRes) == 0 {
						content = "No responses currently registered ;-("
					} else {
						content = "Current fite users are: \n"
						for _, user := range FITE {
							content += fmt.Sprintf("User: %s Mode: %d \n", user._name, user._mode)
						}
					}
				} else {
					content = `Peasant, you are not authorized.`
				}
				s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: content,
					},
				})
			}
			if err != nil {
				log.Printf("Error checking if interaction is DM: %s \n", err)
			}

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

	for _, v := range commands {
		_, err := dg.ApplicationCommandCreate(dg.State.User.ID, GUILDID, v)

		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}

	//Experimenting with overwrite to ensure fresh commands.
	// _, err = dg.ApplicationCommandBulkOverwrite(dg.State.User.ID, GUILDID, commands)
	// if err != nil {
	// 	log.Panicf("Cannot bulk overwrite commands!!")
	// }

	// Wait for a CTRL-C
	log.Printf(`Now running. Press CTRL-C to exit.`)

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
	if m.Author.ID == s.State.User.ID {
		return
	}

	//If the sender is a FITE target, grab their mode and respond accordingly
	fiter, ok := FITE[string(m.Author.ID)]
	if ok {
		var r string
		var err error
		switch fiter._mode {

		//cases defined in the register fite user slash command. TODO: make this defined by an interface or something
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
			//Message handler ran without error. maybe do something idk
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

func removeFite(id string) error {
	_, ok := FITE[id]
	if ok {
		delete(FITE, id)
	}
	return nil
}

func comesFromDM(s *discordgo.Session, i *discordgo.InteractionCreate) (bool, error) {
	channel, err := s.State.Channel(i.ChannelID)
	if err != nil {
		if channel, err = s.Channel(i.ChannelID); err != nil {
			return false, err
		}
	}

	return channel.Type == discordgo.ChannelTypeDM, nil
}

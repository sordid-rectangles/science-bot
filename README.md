# science-bot
A silly discord bot to do silly things

## Reference

Godotenv:
https://github.com/joho/godotenv

DiscordGo
https://github.com/bwmarrin/discordgo/tree/29269347e820c4011fd277948eb8b13308b61bb9

## Configuration
This bot uses godotenv to import environment settings, so to specify the bot's access token create a .env file with the following:

BOT_TOKEN: "<token here>"

Currently the global ADMIN variable in main sets who can talk to the bot. This will change to be configurable
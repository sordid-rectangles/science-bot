# science-bot
A silly discord bot to do silly things

## Reference

Godotenv:
https://github.com/joho/godotenv

dgrouter:
https://github.com/Necroforger/dgrouter/tree/e66453b957c1bcce881b9dabe3d1fe2627aec394

DiscordGo:
https://github.com/bwmarrin/discordgo/tree/29269347e820c4011fd277948eb8b13308b61bb9

Note on DiscordGo:
I am using a nonstandard version of discordgo, as I need slash command support, and currently discord is making lots of breaking changes the discordgo team is having a tough time keeping up with.
ref: https://github.com/bwmarrin/discordgo/wiki/FAQ#application-commands-release-notice

## Configuration
This bot uses godotenv to import environment settings, so to specify the bot's access token, owner, and command prefix create a .env file with the following:

BOT_TOKEN: yourtokenhere
OWNER_NAME: yourusernamehere
OWNER_ID: yourdiscorduseridhere
CMD_PREFIX: yourdesiredprefixhere

## Registering Admins


## Commands

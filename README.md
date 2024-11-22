![Build](https://github.com/wokkipannu/Wokkibot-Go/actions/workflows/build.yml/badge.svg)

# Wokkibot-Go
Wokkibot rewritten yet again. This time in GO using DisGo.

The main purpose of Wokkibot was originally to play music. This has shifted a lot and these days while music playing is still possible, it is not the main focus.
The bot currently still requires lavalink to run, but this will ideally be changed soon so that it can function without it. This would obviously disable the music related commands.

# Commands
##### Music related commands
* `/play` - Play a video or a song from URL or search by given text
* `/skip` - Skip currently palying song
* `/seek` - Skip to a timestamp in current song
* `/volume` - Set the volume of the current song
* `/queue` - List all songs in queue
* `/disconnect` - Disconnect from voice channel and clear queue
##### Other commands
* `/download` - Download a video from URL using yt-dlp and convert using ffmpeg if needed
* `/friday` - Post a random friday celebration clip from SQLite database
* `/roll` - Roll a dice
* `/pizza` - Get random pizza toppings (Currently not working due to API being down)
* `/user` - Get information about a user
* `/flip` - Flip a coin
* `/trivia` - Start a trivia game
* `/joke` - Get a random joke
* `/inspect` - Inspect an image using AI
* `/settings`
  * `/settings commands`
    * `/settings commands add` - Add a custom command
    * `/settings commands remove` - Remove a custom command
    * `/settings commands list` - List all custom commands
  * `/settings llm` - LLM/Ollama specific configuration
    * `/settings llm system-message` - Change the system message for LLM
    * `/settings llm model` - Change the model for LLM
    * `/settings llm history-count` - Change the amount of messages to remember and send for LLM. This includes messages from user and bot, so one response is 2 messages in history.
    * `/settings llm api-url` - Change the API URL for LLM
    * `/settings llm enabled` - Enable or disable responding to @Wokkibot
  * `/settings friday`
    * `/settings friday add` - Add a friday celebration clip
    * `/settings friday remove` - Remove a friday celebration clip
    * `/settings friday list` - List all friday celebration clips
##### Context menu commands
* Quote - Post a message quote as an embed
* Eval - Evaluate code
* Pin - Pin a message to pins channel. This will eventually be configurable using the settings, currently set in config.json

# Setup
* Get [Lavalink](https://github.com/freyacodes/Lavalink)
* Setup config.json
* Run `go run main.go` or build

# Configuration
### config.json example
```
{
 "token": "", // Discord bot token
 "guildid": "", // Discord guild id if you want to restrict commands to a specific guild
 "nodes": [
  {
   "name": "",
   "address": "localhost:2333",
   "password": "youshallnotpass",
   "secure": false,
   "session_id": ""
  }
 ],
 "trivia_token": "",
 "ai_settings": {
  "model": "llama3.2",
  "system": "You are a discord bot",
  "enabled": true,
  "api_url": "http://127.0.0.1:11434",
  "history_count": 10
 },
 "admins": [
  "SOME_DISCORD_USER_ID"
 ],
 "pin_channel": "SOME_DISCORD_CHANNEL_ID"
}
```

### custom_commands.json
>[!NOTE]
> Soon to be removed and moved to SQLite

Create a custom_commands.json file in the same directory as the bot and give it an empty array to start with.
```
[]
```

# TODO
- [x] Add SQLite database for storing data
- [ ] Make custom commands guild based (currently global)
- [ ] Store Trivia token in database for each guild
- [ ] Place music related commands under a music/player command as subcommands
- [x] Store /friday command clips in database
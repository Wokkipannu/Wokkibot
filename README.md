![Build](https://github.com/wokkipannu/Wokkibot-Go/actions/workflows/build.yml/badge.svg)

# Wokkibot-Go
Wokkibot rewritten yet again. This time in GO using DisGo.

The main purpose of Wokkibot was originally to play music. This has shifted a lot and these days while music playing is still possible, it is not the main focus.

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
* `/settings`
  * `/settings commands`
    * `/settings commands add` - Add a custom command
    * `/settings commands remove` - Remove a custom command
    * `/settings commands list` - List all custom commands
  * `/settings friday`
    * `/settings friday add` - Add a friday celebration clip
    * `/settings friday remove` - Remove a friday celebration clip
    * `/settings friday list` - List all friday celebration clips
  * `/settings guild`
    * `/settings guild pinchannel` - Set the pin channel
* `/minesweeper` - Start a minesweeper game
##### Context menu commands
* Quote - Post a message quote as an embed
* Eval - Evaluate code
* Pin - Pin a message to pins channel

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
 "admins": [
  "SOME_DISCORD_USER_ID"
 ],
 "lavalink_enabled": true
}
```

### custom_commands.json
>[!NOTE]
> Soon to be removed and moved to SQLite

Create a custom_commands.json file in the same directory as the bot and give it an empty array to start with.
```
[]
```
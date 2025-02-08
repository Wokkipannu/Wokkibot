![Build](https://github.com/wokkipannu/Wokkibot-Go/actions/workflows/build.yml/badge.svg)

# Wokkibot
Wokkibot is a multi purpose Discord bot built with Go on the [DisGo](https://github.com/disgoorg/disgo) library.

- [C# version of Wokkibot](https://github.com/Wokkipannu/Wokkibot-CSharp) (Not maintained)
- [Original JavaScript version of Wokkibot](https://github.com/Wokkipannu/WokkibotJS) (Not maintained)

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
  * `/settings lavalink`
    * `/settings lavalink toggle` - Toggle lavalink on/off
* `/minesweeper` - Start a minesweeper game
##### Context menu commands
* Quote - Post a message quote as an embed
* Eval - Evaluate code
* Pin - Pin a message to pins channel

# Setup
* Get [Lavalink](https://github.com/freyacodes/Lavalink) and required plugins ([https://github.com/lavalink-devs/youtube-source#plugin](YouTube-Source), [https://github.com/topi314/LavaSrc](LavaSrc) and [https://github.com/topi314/LavaSearch](LavaSearch))
* Setup config.json
* Run `go run main.go` or build

# Configuration
### config.json example
```
{
 "token": "BOT_TOKEN", // Discord bot token
 "guildid": "GUILD_ID", // Discord guild id if you want to restrict commands to a specific guild
 "admins": [
  "ADMIN_USER_ID"
 ],
 "lavalink": {
  "enabled": true,
  "nodes": [
   {
    "name": "Lavalink",
    "address": "localhost:2333",
    "password": "youshallnotpass",
    "secure": false,
    "session_id": ""
   }
  ]
 }
}
```

# Custom commands
Custom commands can include attributes to do specific things. Currently the supported attributes are:
* `{{time|timezone}}` - Returns the current time in the specified timezone
* `{{random|choices;separated;by;semicolons}}` - Returns a random choice from the provided choices
* `{{user|attribute}}` - Returns a user attribute

Example:
```
{{time|Europe/Helsinki}}
{{random|Hello;World;This;Is;A;Random;Choice}}
{{user|name}}
{{user|id}}
{{user|avatar}}
{{user|mention}}
{{user|created}}
```

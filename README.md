![Build](https://github.com/wokkipannu/Wokkibot-Go/actions/workflows/build.yml/badge.svg)

# Wokkibot-Go
Wokkibot rewritten yet again. This time in GO using DisGo.

# Commands
##### Music related commands
* /play
* /skip
* /seek
* /volume
* /queue
* /disconnect
##### Other commands
* /friday
* /roll
* /pizza
* /user
* /flip
* /trivia
* /joke
* /settings
  * /settings commands
    * /settings commands add
    * /settings commands remove
    * /settings commands list
##### Context menu commands
* Quote
* Eval

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
 ]
}
```

### custom_commands.json
Create a custom_commands.json file in the same directory as the bot and give it an empty array to start with.
```
[]
```

# TODO
- [ ] Add SQLite database for storing data
- [ ] Make custom commands guild based (currently global)
- [ ] Store Trivia token in database for each guild
- [ ] Place music related commands under a music/player command as subcommands
- [ ] Store /friday command clips in database
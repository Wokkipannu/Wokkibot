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
* /trivia
##### Other commands
* /friday
* /roll
* /pizza
* /user
* /flip
##### Context menu commands
* Quote
* Eval

# Setup
* Get [Lavalink](https://github.com/freyacodes/Lavalink)
* Setup config.json
* Run `go run main.go` or build

# config.json example
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
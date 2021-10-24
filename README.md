# Wokkibot
No longer updating this. Continuing development on CSharp version [Wokkibot CSharp](https://github.com/Wokkipannu/Wokkibot-CSharp).

A Discord bot written using DiscordJS and Commando.

# Requirements
- Node.js 10.0.0 or newer is required due to DiscordJS master branch being used

# Installation
- `git clone https://github.com/Wokkipannu/Wokkibot.git`
- `cd Wokkibot`
- `npm install`
- Create config.json in root with the following
```
{
  "TOKEN": "YOUR-DISCORD-BOT-TOKEN",
  "OWNER": "YOUR-ACCOUNT-ID",
  "PREFIX": "?",
  "GOOGLE_API_KEY": "YOUR-GOOGLE-API-KEY"
}
```
- `npm start`

# Commands
- (prefix)help - Get command list
- (prefix)play <url/keyword> - Play a song
- (prefix)skip - Skip currently playing song
- (prefix)queue - Display song queue
- (prefix)remove <id> - Remove specific song from queue
- (prefix)volume <1-100> - Change dispatcher volume
- (prefix)weather <location> - Get weather for given location
- (prefix)accountage <Optional: @User> - Get date when account was created
- (prefix)purge <limit> - Delete x messages from channel
- (prefix)dice <Optional: Times> <Optional: Max> - Roll the dice
- (prefix)8ball <Question?> - Ask magic 8ball a question

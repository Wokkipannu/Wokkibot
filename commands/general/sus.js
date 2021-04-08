const { Command } = require('discord.js-commando');

module.exports = class SusCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'sus',
      group: 'general',
      memberName: 'sus',
      description: 'See how sus you are',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'user',
          prompt: 'Maximum roll?',
          type: 'integer',
          default: 100
        }
      ]
    });
  }

  run(msg, { user }) {
    user = msg.mentions.users ? msg.mentions.users : msg.author

    let susLevel = Math.floor(Math.random() * 100)

    return msg.channel.send(`${user} is ${susLevel}% sus!`)
  }
}
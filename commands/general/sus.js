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
          prompt: 'Who to check for sus?',
          type: 'string',
          default: ''
        }
      ]
    });
  }

  run(msg, { user }) {
    user = user ? msg.mentions.users.first() : msg.author

    if (!user || user === undefined) {
      return msg.reply('Invalid target! Use `!sus` or `!sus @someone` to check specific user')
    }

    let susLevel = Math.floor(Math.random() * 101)

    return msg.channel.send(`${user.username} is ${susLevel}% sus!`)
  }
}
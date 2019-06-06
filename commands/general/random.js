const { Command } = require('discord.js-commando');

module.exports = class RandomCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'random',
      group: 'general',
      memberName: 'random',
      description: 'Roll random from given options',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'list',
          prompt: 'Enter list of items to select from divided with commas',
          type: 'string'
        }
      ]
    });
  }

  run(msg, { list }) {
    list = list.split(",");
    if (list.length > 0) {
      const random = Math.floor(Math.random() * list.length);
      msg.reply(list[random]);
    }
    else {
      msg.reply('Your list sucks');
    }

    return undefined;
  }
}
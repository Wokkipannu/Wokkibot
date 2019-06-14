const { Command } = require('discord.js-commando');

module.exports = class SkipCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'skip',
      group: 'music',
      memberName: 'skip',
      description: 'Skip currently playing song',
      guildOnly: true,
      clientPermissions: ['CONNECT', 'SPEAK', 'SEND_MESSAGES']
    });
  }

  async run(msg) {
    const queue = await this.queue.get(msg.guild.id);
    if (!queue) return msg.reply('You can not skip nonexistent song');

    const unskippable = ["108299947257925632", "108617380552273920", "117985849257230345"];

    if (unskippable.includes(queue.songs[0].requester.id) && !unskippable.includes(msg.author.id)) {
      return msg.reply(`You can not skip song requested by ${queue.songs[0].requester}. Ask them to skip!`);
    }

    const title = queue.songs[0].title;
    queue.connection.dispatcher.end('skipped');
    return msg.channel.send(`${title} was skipped`);
  }

  get queue() {
    return this.client.registry.resolveCommand('music:play').queue;
  }
}
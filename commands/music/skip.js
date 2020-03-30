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
    if (!queue) return msg.reply('No songs in queue');

    queue.connection.dispatcher.end();
    return msg.channel.send('Song skipped');
  }

  get queue() {
    return this.client.registry.resolveCommand('music:play').queue;
  }
}
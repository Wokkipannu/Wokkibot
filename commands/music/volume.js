const { Command } = require('discord.js-commando');

module.exports = class VolumeCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'volume',
      group: 'music',
      memberName: 'volume',
      description: 'Change current song volume',
      guildOnly: true,
      clientPermissions: ['CONNECT', 'SPEAK', 'SEND_MESSAGES'],
      args: [
        {
          key: 'volume',
          prompt: 'Enter volume',
          type: 'integer'
        }
      ]
    });
  }

  async run(msg, { volume }) {
    const queue = await this.queue.get(msg.guild.id);
    if (!queue) return msg.reply('No songs in queue');

    try {
      queue.connection.dispatcher.setVolume(volume / 100);
      return msg.channel.send('Volume changed');
    }
    catch(error) {
      console.error(error);
      return msg.reply('Volume change failed');
    }
  }

  get queue() {
    return this.client.registry.resolveCommand('music:play').queue;
  }
}
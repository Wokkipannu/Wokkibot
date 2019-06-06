const { Command } = require('discord.js-commando');
const { MessageEmbed } = require('discord.js');

module.exports = class QueueCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'queue',
      group: 'music',
      memberName: 'queue',
      description: 'Display song queue',
      guildOnly: true,
      clientPermissions: ['SEND_MESSAGES']
    });
  }

  async run(msg) {
    const queue = await this.queue.get(msg.guild.id);
    if (!queue) return msg.reply('There are no songs in queue');

    const queueEmbed = new MessageEmbed()
      .setColor('#1a2b3c')
      .setTitle('Song queue');
    
    queue.songs.forEach((song, index) => {
      queueEmbed.addField(`Song ${index + 1} ${index === 0 ? '(Now playing)' : ''}`, `[${song.title}](https://www.youtube.com/watch?v=${song.id})\nDuration: ${this.timeString(queue.songs[index].duration)}\nRequested by ${song.requester}`, false);
    });

    return msg.channel.send(queueEmbed);
  }

  timeString(seconds, forceHours = false) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor(seconds % 3600 / 60);

    return `${forceHours || hours >= 1 ? `${hours}:` : ''}${hours >= 1 ? `0${minutes}`.slice(-2) : minutes}:${`0${Math.floor(seconds % 60)}`.slice(-2)}`;
  }

  get queue() {
    return this.client.registry.resolveCommand('music:play').queue;
  }
}
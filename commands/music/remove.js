const { Command } = require('discord.js-commando');

module.exports = class RemoveCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'remove',
      group: 'music',
      memberName: 'remove',
      description: 'Remove a song from the queue',
      guildOnly: true,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'index',
          prompt: 'Song number in queue',
          type: 'integer'
        }
      ]
    });
  }

  async run(msg, { index }) {
    const queue = await this.queue.get(msg.guild.id);
    if (!queue) return msg.reply('There are no songs in queue');

    if (index <= 1) return msg.reply('You can not remove first song. Use skip instead.');
    if (index > queue.songs.length) return msg.reply('No song with given index');

    const unremovable = ["108299947257925632", "108617380552273920", "117985849257230345"];

    if (unremovable.includes(queue.songs[index].requester.id) && !unremovable.includes(msg.author.id)) {
      return msg.reply(`You can not remove song requested by ${queue.songs[index].requester}. Ask them to remove!`);
    }

    queue.songs.splice(index - 1, 1);

    return msg.channel.send(`Removed song ${index} from queue`);
  }

  get queue() {
    return this.client.registry.resolveCommand('music:play').queue;
  }
}
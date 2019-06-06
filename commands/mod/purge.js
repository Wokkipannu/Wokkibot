const { Command } = require('discord.js-commando');
const { MessageEmbed } = require('discord.js');

module.exports = class PurgeCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'purge',
      aliases: ['clear', 'clean'],
      group: 'mod',
      memberName: 'purge',
      description: 'Purge messages from chat',
      clientPermissions: ['MANAGE_MESSAGES', 'SEND_MESSAGES', 'EMBED_LINKS'],
      userPermissions: ['MANAGE_MESSAGES'],
      args: [
        {
          key: 'limit',
          prompt: 'How many messages should be deleted?',
          type: 'integer'
        }
      ]
    });
  }

  async run(msg, { limit }) {
    await msg.channel.messages.fetch({ limit: limit + 1 }).then(messages => {
      msg.channel.bulkDelete(messages.array().reverse()).then(msgs => {
        const embed = new MessageEmbed()
          .setColor('#00ff1d')
          .setTitle('Purge')
          .setDescription(`Deleted ${msgs.size - 1} messages`)
          .setFooter('This message will be deleted automatically in 5 seconds');
        
        msg.channel.send(embed).then(msg => msg.delete({ timeout: 5000 }));
      }).catch(err => {
        return [this.client.logger.error('Purge command error', err),msg.channel.send('Could not bulk delete messages')];
      });
    }).catch(err => {
      return [this.client.logger.error('Purge command error', err),msg.channel.send('Could not fetch messages')];
    });

    return undefined;
  }
}
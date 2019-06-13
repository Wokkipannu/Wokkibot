const { Command } = require('discord.js-commando');

module.exports = class UnknownCommandCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'unknown-command',
      group: 'util',
      memberName: 'unknown-command',
      description: '',
      unknown: true,
      hidden: true
    });
  }

  run(msg) {
    if (msg.content.includes('vittu')) {
      return msg.reply('Vittu säkin olet lapsellinen? Kysyn kiltisti jotain niin tommonen vitu chättibodyguard ei voi ottaa chat banneja pois mitkä oon saanu aivan vitun turhaan kun heitin läppää?');
    }
    else return;
  }
}
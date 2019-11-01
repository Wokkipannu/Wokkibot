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
    return;
  }
}
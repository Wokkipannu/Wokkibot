const { Command } = require('discord.js-commando');

module.exports = class ClassCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'class',
      group: 'general',
      memberName: 'class',
      description: 'Roll random class',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES']
    });
  }

  run(msg) {
    const classes = ["Druid", "Hunter", "Mage", "Priest", "Rogue", "Shaman", "Warlock", "Warrior"];
    const random = classes[Math.floor(Math.random() * classes.length)];

    return msg.reply(`you play ${random}`);
  }
}
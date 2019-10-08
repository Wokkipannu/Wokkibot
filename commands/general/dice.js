const { Command } = require('discord.js-commando');

module.exports = class DiceCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'dice',
      aliases: ['die', 'roll'],
      group: 'general',
      memberName: 'dice',
      description: 'Roll the dice',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'max',
          prompt: 'Maximum roll?',
          type: 'integer',
          default: 100
        },
        {
          key: 'times',
          prompt: 'How many times?',
          type: 'integer',
          default: 1
        }
      ]
    });
  }

  run(msg, { max, times }) {
    let rolls = [];

    for (let i = 0; i < times; i++) {
      rolls.push(Math.floor(Math.random() * max));
    }

    if (rolls.length > 1) {
      const sum = rolls.reduce((a, b) => a + b);
      msg.reply(`You rolled ${rolls.join(", ")} (1-${max}) (Total: ${sum})`);
    } 
    else {
      msg.reply(`You rolled ${rolls.join(", ")} (1-${max})`);
    }
    
    return undefined;
  }
}
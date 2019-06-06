const { Command } = require('discord.js-commando');

module.exports = class DiceCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'dice',
      aliases: ['die'],
      group: 'general',
      memberName: 'dice',
      description: 'Roll the dice',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'times',
          prompt: 'How many times?',
          type: 'integer',
          default: 1
        },
        {
          key: 'max',
          prompt: 'Maximum roll?',
          type: 'integer',
          default: 100
        }
      ]
    });
  }

  run(msg, { times, max }) {
    let rolls = [];

    for (let i = 0; i < times; i++) {
      rolls.push(Math.floor(Math.random() * max));
    }

    if (rolls.length > 1) {
      const sum = rolls.reduce((a, b) => a + b);
      msg.reply(`You rolled ${rolls.join(", ")} (Total: ${sum})`);
    } 
    else {
      msg.reply(`You rolled ${rolls.join(", ")}`);
    }
    
    return undefined;
  }
}
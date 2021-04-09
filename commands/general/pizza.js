const { Command } = require('discord.js-commando')
const toppings = require('../../toppings.json')

module.exports = class PizzaCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'pizza',
      group: 'general',
      memberName: 'pizza',
      description: 'Get random pizza toppings',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'amount',
          prompt: 'How many toppings?',
          type: 'integer',
          default: 4
        }
      ]
    })
  }

  run (msg, { amount }) {
    if (amount > 10) amount = 10
    if (amount < 1) return msg.reply('Syö pizzasi ilman täytteitä tai anna suurempi luku')

    let selectedToppings = []

    for (let i = 0; i < amount; i++) {
      selectedToppings.push(toppings[i])
    }

    let lastTopping = selectedToppings.splice(selectedToppings.length - 1, 1)[0]

    if (amount === 1) {
      return msg.reply(`Pizzaasi tuli täyte: ${selectedToppings[0]}`)
    } else {
      return msg.reply(`Pizzaasi tuli täytteet: ${selectedToppings.join(', ')} ja ${lastTopping}`)
    }
  }
}
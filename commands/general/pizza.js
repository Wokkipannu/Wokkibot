const { Command } = require('discord.js-commando')
const toppings = require('../../toppings.json')

module.exports = class PizzaCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'pizza',
      group: 'general',
      memberName: 'pizza',
      aliases: ['täytteet', 'toppings', 'taytteet', 'pitsa'],
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

    let myToppings = toppings

    let selectedToppings = []

    for (let i = 0; i < amount; i++) {
      let randomTopping = myToppings[Math.floor(Math.random() * myToppings.length)]
      selectedToppings.push(randomTopping)
      myToppings = myToppings.filter(t => t !== randomTopping)
    }

    if (amount === 1) {
      return msg.reply(`Pizzaasi tuli täyte: ${selectedToppings[0]}`)
    } else {
      let lastTopping = selectedToppings.splice(selectedToppings.length - 1, 1)[0]
      return msg.reply(`Pizzaasi tuli täytteet: ${selectedToppings.join(', ')} ja ${lastTopping}`)
    }
  }
}
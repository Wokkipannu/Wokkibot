const { Command } = require('discord.js-commando');

module.exports = class UnknownCommandCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'quote',
      group: 'general',
      memberName: 'quote',
      description: 'Get a random quote'
    });
  }

  run(msg) {
    const possibilities = [
      "mutta nyt :zipper_mouth:",
      "Vittu säkin olet lapsellinen? Kysyn kiltisti jotain niin tommonen vitu chättibodyguard ei voi ottaa chat banneja pois mitkä oon saanu aivan vitun turhaan kun heitin läppää?",
      "http://www.tilaapullo.com/",
      "Hei, olen Tohtori Gerhard",
      "Jopa 8cm pidempi penis kahdessa viikossa",
      "mee lahjottaan tonni hyväntekeväisyyteen ja tuu sitten neuvoon elämässä",
      "new to this channel <----"
    ];

    const random = Math.floor(Math.random() * possibilities.length);
    const quote = possibilities[random];

    return msg.reply(quote);
  }
}
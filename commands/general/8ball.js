const { Command } = require('discord.js-commando');

module.exports = class EightballCommand extends Command {
  constructor(client) {
    super(client, {
      name: '8ball',
      group: 'general',
      memberName: '8ball',
      description: 'Get an answer to your question from the magic 8ball!',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'question',
          prompt: 'What would you like to ask the magic 8ball?',
          type: 'string'
        }
      ]
    });
  }

  run(msg, { question }) {
    if (!question.endsWith("?")) return msg.reply('Is that a question?');

    const answers = [
      'It is certain.', 'It is decidedly so.', 'Without a doubt.',
      'Yes - definitely.', 'You may rely on it.', 'As I see it, yes.',
      'Most likely.', 'Outlook good.', 'Yes', 'Signs point to yes.',
      'Reply hazy, try again.', 'Ask again later.', 'Better not tell you now.',
      'Cannot predict now.', 'Concentrate and ask again.', 'Don\'t count on it.',
      'My reply is no.', 'My sources say no.', 'Outlook not so good.', 'Very doubtful.'
    ];

    const random = Math.floor(Math.random() * answers.length);
    const answer = answers[random];

    return msg.reply(answer);
  }
}
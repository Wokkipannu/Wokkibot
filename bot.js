const Commando = require('discord.js-commando');
const winston = require('winston');
const path = require('path');

const client = new Commando.Client({
  owner: process.env.OWNER,
  commandPrefix: process.env.PREFIX
});

client.logger = winston.createLogger({
  format: winston.format.combine(
    winston.format.timestamp({
      format: 'DD.MM.YYYY HH:mm:ss'
    }),
    winston.format.printf(info => `[${info.timestamp}] ${info.level}: ${info.message}`)
  ),
  transports: [
    new winston.transports.Console(),
    new winston.transports.File({ filename: 'wokkibot.log' })
  ]
});

client
  .on('ready', () => {
    client.logger.info(`Logged in as ${client.user.tag}`);
    client.user.setActivity('you', { type: 'WATCHING'} );
  })
  .on('warn', client.logger.error)
  .on('error', client.logger.warn)
  .on('commandRun', (cmd, promise, msg, args) => client.logger.info(`${msg.author.tag} (${msg.author.id}) ran command ${cmd.groupID}:${cmd.memberName}`))
  .on('commandError', (cmd, err) => client.logger.error(`Error occurred when running command ${cmd.groupID}:${cmd.memberName}`, err));

client.registry
  .registerGroups([
    ['music', 'Music commands'],
    ['general', 'General commands'],
    ['mod', 'Moderation commands']
  ])
  .registerDefaults()
  .registerCommandsIn(path.join(__dirname, 'commands'));

client.login(process.env.TOKEN);
/**
 * Wokkibot
 * 
 * A multipurpose discord bot with a focus in playng music in voice channels
 */

require('dotenv').config()

const { CommandoClient } = require('discord.js-commando');
const path = require('path');

const client = new CommandoClient({
  commandPrefix: process.env.PREFIX,
  owner: process.env.OWNER
});

client.registry
  .registerDefaultTypes()
  .registerGroups([
    ['general', 'General commands'],
    ['music', 'Music playing commands']
  ])
  .registerDefaultGroups()
  .registerDefaultCommands()
  .registerCommandsIn(path.join(__dirname, 'commands'));

client.once('ready', () => {
  console.log(`Logged in as ${client.user.tag}`);
  client.user.setActivity('you', { type: 'WATCHING' });
});

client.on('error', console.error);

client.login(process.env.TOKEN);
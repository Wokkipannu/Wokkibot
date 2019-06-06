const { Command } = require('discord.js-commando');

module.exports = class WowCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'wow',
      group: 'general',
      memberName: 'wow',
      description: 'Get time until WoW Classic release',
      guildOnly: false,
      clientPermissions: ['SEND_MESSAGES']
    });
  }

  run(msg) {
    let target = new Date('8/27/2019');
    target.setHours(target.getHours() - 2);

    const difference = target.getTime() - new Date().getTime();

    const days = this.getDays(difference / 1000);
    const hours = this.getHours(difference / 1000);
    const minutes = this.getMinutes(difference / 1000);
    const seconds = this.getSeconds(difference / 1000);

    let reply = "";
    reply += days !== 1 ? `${days} days, ` : `${days} day, `;
    reply += hours !== 1 ? `${hours} hours, ` : `${hours} hour, `;
    reply += minutes !== 1 ? `${minutes} minutes, ` : `${minutes} minute, `;
    reply += seconds !== 1 ? `${seconds} seconds` : `${seconds} second`;

    return msg.reply(`**${reply}** until WoW Classic`);
  }

  getDays(seconds) {
    const days = Math.floor(seconds / 86400);
    return days;
  }
  
  getHours(seconds) {
    const hours = Math.floor(seconds % 86400 / 3600);
    return hours;
  }
  
  getMinutes(seconds) {
    const minutes = Math.floor(seconds % 3600 / 60);
    return minutes;
  }
  
  getSeconds(seconds) {
    return Math.floor(seconds % 60);
  }
}
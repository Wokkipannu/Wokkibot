const { Command } = require('discord.js-commando');

module.exports = class IQCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'iq',
      group: 'fun',
      memberName: 'iq',
      description: 'How high is your IQ',
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'user',
          prompt: 'User to check IQ from',
          type: 'string',
          default: ''
        }
      ]
    });
  }

  run(msg, { user }) {
    // Get user data
    user = user ? msg.mentions.users.first(1)[0] : msg.author;

    const PepoClap = this.client.emojis.find(emoji => emoji.name === "feelspepoclap");
    const PepoDance = this.client.emojis.find(emoji => emoji.name === "PepoDance");
    const PepoSpin = this.client.emojis.find(emoji => emoji.name === "feelspepospin");

    let randomIQ = Math.floor(this.normalRandom() * 200 + 100) + 1;
    while (Number.isNaN(randomIQ)) {
      randomIQ = Math.floor(this.normalRandom() * 200 + 100) + 1;
    }

    this.client.db.getUser(user.id)
      .then(data => {
        if (data && data.iq) {
          if (data.iq < 90) return msg.channel.send(`${user}'s IQ is **${data.iq}** ${PepoSpin}`);
          else if (data.iq > 90 && data.iq < 110) return msg.channel.send(`${user}'s IQ is **${data.iq}** ${PepoClap}`);
          else if (data.iq > 110) return msg.channel.send(`${user}'s IQ is **${data.iq}** ${PepoDance}`);
          else return msg.channel.send(`${user}'s IQ is **${data.iq}**`);
        }
        else if(user.id == 128685552450011137){
           this.client.db.updateUser(user.id, "iq", 200);
           if (data.iq < 90) return msg.channel.send(`${user}'s IQ is **${data.iq}** ${PepoSpin}`);
           else if (data.iq > 90 && data.iq < 110) return msg.channel.send(`${user}'s IQ is **${data.iq}** ${PepoClap}`);
           else if (data.iq > 110) return msg.channel.send(`${user}'s IQ is **${data.iq}** ${PepoDance}`);
           else return msg.channel.send(`${user}'s IQ is **${data.iq}**`);
        }
        else {
          this.client.db.updateUser(user.id, "iq", randomIQ);
          if (randomIQ < 90) return msg.channel.send(`${user}'s IQ is **${randomIQ}** ${PepoSpin}`);
          else if (randomIQ > 90 && randomIQ < 110) return msg.channel.send(`${user}'s IQ is **${randomIQ}** ${PepoClap}`);
          else if (randomIQ > 110) return msg.channel.send(`${user}'s IQ is **${randomIQ}** ${PepoDance}`);
          else return msg.channel.send(`${user}'s IQ is **${randomIQ}**`);
        }
      })
      .catch(e => {
        return [this.client.logger.error(e),msg.channel.send('An error occurred. More information logged to console.')];
      });
  }

  normalRandom() {
    var val, u, v, s, mul;

    if(this.spareRandom !== null) {
      val = this.spareRandom;
      this.spareRandom = null;
    }
    else {
      do {
          u = Math.random()*2-1;
          v = Math.random()*2-1;

          s = u*u+v*v;
      } while(s === 0 || s >= 1);

      mul = Math.sqrt(-2 * Math.log(s) / s);

      val = u * mul;
      this.spareRandom = v * mul;
    }
    
    return val / 14;
  }
}

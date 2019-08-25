const { Command } = require('discord.js-commando');
const SuperAgent = require('superagent');

module.exports = class RemoveCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'set',
      group: 'wow',
      memberName: 'set',
      description: 'Set a specified WoW setting',
      guildOnly: true,
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'key',
          prompt: 'Enter name, race, class, spec, profession1 or profession2',
          type: 'string'
        },
        {
          key: 'value',
          prompt: 'What should this be set to?',
          type: 'string'
        }
      ]
    });

    this.token = ''
  }

  async run(msg, { key, value}) {
    //if (msg.channel.id !== "610585419326685242") return;

    if (!this.token) {
      await SuperAgent
        .post(`${process.env.API}/api/v1/users/login`)
        .send({ username: process.env.API_USERNAME, password: process.env.API_PASSWORD })
        .end(async (err, res) => {
          if (err) return msg.reply(`Unable to login to API. Try again later. ${err}`);
          this.token = res.body.data;
        })
    }

    let players = [];
    let player = {};
    let races = ['N/A', 'Orc', 'Tauren', 'Troll', 'Undead']
    let classes = ['N/A', 'Druid', 'Hunter', 'Mage', 'Priest', 'Rogue', 'Shaman', 'Warlock', 'Warrior']
    let professions = ['N/A', 'Blacksmithing', 'Engineering', 'Herbalism', 'Mining', 'Leatherworking', 'Tailoring', 'Enchanting', 'Alchemy', 'Skinning']

    // Get all players from the API
    await SuperAgent
      .get(`${process.env.API}/api/v1/players`)
      .end(async (err, res) => {
        if (err) return msg.reply(`Failed to get players from API`);
        players = await res.body.data;

        // Check if the player already exists
        player = await players.find(p => p.discord === msg.author.id);
        if (!player) {
          // Add player
          player = await this.createPlayer(msg.author.username, msg.author.id);
        }

        /**
         * Change name
         */
        if (key === 'name') {
          player.name = value;
        }
        /**
         * Change race
         */
        else if (key === 'race') {
          if (races.includes(value)) {
            player.race = races.find(race => race.toLowerCase() === value.toLowerCase());
          }
          else {
            return msg.reply(`Valid races are: ${races.join(', ')}`);
          }
        }
        /**
         * Change class
         */
        else if (key === 'class') {
          if (classes.includes(value)) {
            player.class = classes.find(c => c.toLowerCase() === value.toLowerCase());
          }
          else {
            return msg.reply(`Valid classes are: ${classes.join(', ')}`);
          }
        }
        /**
         * Change spec
         */
        else if (key === 'spec') {
          player.spec = value;
        }
        /**
         * Change profession 1
         */
        else if (key === 'profession1') {
          if (professions.includes(value)) {
            player.prof1 = professions.find(prof => prof.toLowerCase() === value.toLowerCase());
          }
          else {
            return msg.reply(`Valid professions are: ${professions.join(', ')}`);
          }
        }
        /**
         * Change profession 2
         */
        else if (key === 'profession2') {
          if (professions.includes(value)) {
            player.prof2 = professions.find(prof => prof.toLowerCase() === value.toLowerCase());
          }
          else {
            return msg.reply(`Valid professions are: ${professions.join(', ')}`);
          }
        }
        /**
         * Invalid key response
         */
        else {
          return msg.reply(`Correct usage is **!set <race/class/spec/profession1/profession2> <value>**`);
        }

        await SuperAgent
          .put(`${process.env.API}/api/v1/players/${player._id}?token=${this.token}`)
          .send(player)
          .end((err, res) => {
            if (err) return msg.reply(`Update player failed: ${err}`);
            return msg.reply(`User updated`);
          });
      });
  }

  async createPlayer(name, id) {
    return new Promise((resolve, reject) => {
      let player = {
        name: name,
        race: 'N/A',
        class: 'N/A',
        spec: 'N/A',
        prof1: 'N/A',
        prof2: 'N/A',
        dkp: 10,
        discord: id,
        token: this.token
      };
  
      SuperAgent
        .post(`${process.env.API}/api/v1/players`)
        .send(player)
        .end((err, res) => {
          if (err) return reject(err);
          return resolve(res.body.data);
        });
    });
  }
}
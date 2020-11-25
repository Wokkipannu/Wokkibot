const { Command } = require('discord.js-commando');
const { Collection } = require('discord.js');
const ytdl = require("discord-ytdl-core");

module.exports = class PlayCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'play',
      group: 'music',
      memberName: 'play',
      description: 'Play audio in a voice channel',
      guildOnly: true,
      clientPermissions: ['CONNECT', 'SPEAK', 'SEND_MESSAGES'],
      args: [
        {
          key: 'url',
          prompt: 'Enter URL or keyword',
          type: 'string'
        },
        {
          key: 'bass',
          prompt: '',
          type: 'string',
          default: '10'
        }
      ]
    });

    this.queue = new Collection();
  }

  async run(msg, { url, bass }) {
    if (!msg.member.voice.channel) return msg.channel.send("You must be connected to voice channel to play");

    let queue = await this.queue.get(msg.guild.id);

    if (url.match(/^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$/g)) {
      const song = {
        url,
        bass
      };

      if (queue) {
        await queue.songs.push(song);
        msg.reply("Added to queue");
      }
      else {
        await this.queue.set(msg.guild.id, { id: msg.guild.id, msg, songs: [song], connection: undefined });
        queue = await this.queue.get(msg.guild.id);
        msg.reply("Added to queue");
        return this.play(queue);
      }
    }
    else {
      return msg.channel.send("Play command must be supplied with a youtube link");
    }
  }

  async play(queue) {
    const { msg } = queue;

    if (queue.songs.length === 0) {
      msg.channel.send("No more songs in queue");
      queue.connection.channel.leave();
      return this.queue.delete(queue.id);
    }

    const { url, bass } = queue.songs[0];

    let stream = await ytdl(url, {
      filter: "audioonly",
      opusEncoded: false,
      fmt: "mp3",
      encoderArgs: ['-af', `bass=g=${bass},dynaudnorm=f=200`]
    });

    if (!queue.connection || queue.connection === undefined) {
      await msg.member.voice.channel.join().then(connection => queue.connection = connection);
    }

    let dispatcher = queue.connection.play(stream, {
      type: "unknown",
      volume: 0.1
    })
    .on("finish", () => {
      queue.songs.shift();
      return this.play(queue);
    });

  }
}
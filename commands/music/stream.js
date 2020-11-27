const { Command } = require('discord.js-commando');
const ytdl = require("discord-ytdl-core");

module.exports = class StreamCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'stream',
      group: 'music',
      memberName: 'stream',
      description: 'Play a stream',
      guildOnly: true,
      clientPermissions: ['CONNECT', 'SPEAK', 'SEND_MESSAGES'],
      args: [
        {
          key: 'url',
          prompt: 'URL to stream',
          type: 'string'
        }
      ]
    })
  }

  async run(msg, { url }) {
    if (!msg.member.voice.channel) return msg.channel.send("You must be connected to voice channel to play");

    let stream = ytdl.arbitraryStream(url, {
      opusEncoded: false,
      fmt: "mp3",
      encoderArgs: ['-af', 'bass=g=10,dynaudnorm=f=200']
    });

    msg.member.voice.channel.join()
      .then(connection => {
        let dispatcher = connection.play(stream, {
          type: 'unknown'
        })
        .on('finish', () => {
          msg.guild.me.voice.channel.leave();
        });
      });
  }
}
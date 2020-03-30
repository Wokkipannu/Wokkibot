const { Command } = require('discord.js-commando');
const { Collection } = require('discord.js');
const ytdl = require('ytdl-core-discord');

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
        }
      ]
    });

    this.queue = new Collection();
  }

  async run(msg, { url }) {
    const channel = msg.member.voice.channel;
    if (!channel) return msg.reply('You can only use play command when in a voice channel');

    msg.reply('Handling song request').then(async message => {
      // Get queue
      const queue = await this.queue.get(msg.guild.id);

      // If URL is youtube URL
      if (url.match(/^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$/g)) {
        await ytdl.getInfo(url)
          .then(async info => {
            const song = {
              type: 'yt',
              id: info.video_id,
              url: info.video_url,
              title: info.title,
              length: info.length_seconds,
              requestedBy: msg.author
            };

            if (queue) {
              queue.songs.push(song);
            }
            else {
              await channel.join()
                .then(async connection => {
                  const newQueue = {
                    connection: connection,
                    songs: [song]
                  };

                  await this.queue.set(msg.guild.id, newQueue);
                  this.play(connection);
                })
                .catch(error => {
                  console.error(error);
                  message.edit('Failed to join your voice channel');
                });
            }
            message.edit(`${info.title} added to queue`);
          })
          .catch(error => {
            console.error(error);
            return message.edit('Error occurred while trying to get song info');
          });
      }
      else {
        if (url.endsWith('.mp3')) {
          const song = {
            type: 'mp3',
            id: 'N/A',
            url: url,
            title: 'Unknown',
            length: 0,
            requestedBy: msg.author
          };

          if (queue) {
            queue.songs.push(song);
          }
          else {
            await channel.join()
              .then(async connection => {
                const newQueue = {
                  connection: connection,
                  songs: [song]
                };

                await this.queue.set(msg.guild.id, newQueue);
                this.play(connection);
              })
              .catch(error => {
                console.error(error);
                message.edit('Failed to join your voice channel');
              });
          }

          message.edit(`mp3 file added to queue`);
        }
        else {
          return message.edit('URL has to be YouTube video or mp3 file');
        }
      }
    });
  }

  async play(connection) {
    const queue = await this.queue.get(connection.channel.guild.id);
    if (queue.songs.length === 0) {
      connection.channel.leave();
      this.queue.delete(connection.channel.guild.id);
    }
    else {
      console.log(`Starting to play song ${queue.songs[0].title}`);
      if (queue.songs[0].type === 'yt') {
        connection.play(await ytdl(queue.songs[0].url), { type: 'opus', volume: 0.1 })
          .on('finish', () => {
            console.log('Song ended or it was skipped');
            queue.songs.shift();
            this.play(connection);
          })
          .on('error', (error) => {
            console.error(error);
            queue.songs.shift();
            this.play(connection);
          });
      }
      else {
        connection.play(queue.songs[0].url)
          .on('finish', () => {
            console.log('Song ended or it was skipped');
            queue.songs.shift();
            this.play(connection);
          })
          .on('error', (error) => {
            console.error(error);
            queue.songs.shift();
            this.play(connection);
          });
      }
    }
  }
}
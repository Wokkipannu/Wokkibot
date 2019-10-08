const { Command } = require('discord.js-commando');
const { Util, MessageEmbed } = require('discord.js');
const YouTube = require('simple-youtube-api');
const ytdl = require('ytdl-core-discord');

const GOOGLE_API_KEY = process.env.G_API_KEY

module.exports = class PlayCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'play',
      group: 'music',
      memberName: 'play',
      description: 'Play audio from YouTube URL or keyword',
      clientPermissions: ['CONNECT', 'SPEAK', 'SEND_MESSAGES', 'EMBED_LINKS'],
      guildOnly: true,
      args: [
        {
          key: 'url',
          prompt: 'Enter URL or keyword',
          type: 'string'
        }
      ]
    });

    this.queue = new Map();
    this.youtube = new YouTube(GOOGLE_API_KEY);
  }

  async run(msg, { url }) {
    const voiceChannel = msg.member.voice.channel;
    if (!voiceChannel) return msg.reply('You must connect to a voice channel first');

    if (url.match(/^(https?\:\/\/)?(www\.youtube\.com|youtu\.?be)\/.+$/g)) {
      let video = await this.youtube.getVideo(url);
      this.handleVideo(msg, video);
    }
    else {
      let videos = await this.youtube.searchVideos(url, 1)
        .catch(error => {
          if(error) this.client.logger.error('Play command error', error);
        });
      let video = await this.youtube.getVideoByID(videos[0].id);
      this.handleVideo(msg, video);
    }
  }

  async handleVideo(msg, video) {
    const voiceChannel = await msg.member.voice.channel;
    const queue = await this.queue.get(msg.guild.id);

    const song = {
      id: video.id,
      title: Util.escapeMarkdown(video.title),
      url: `https://www.youtube.com/watch?v=${video.id}`,
      thumbnail: `https://img.youtube.com/vi/${video.id}/mqdefault.jpg`,
      duration: video.durationSeconds ? video.durationSeconds : video.duration / 1000,
      requester: msg.author
    }

    if (!queue) {
      const queueConstruct = {
        voiceChannel: voiceChannel,
        connection: null,
        songs: []
      };
      queueConstruct.songs.push(song);
      this.queue.set(msg.guild.id, queueConstruct);

      let connection = await voiceChannel.join();
      queueConstruct.connection = connection;
      this.play(msg);
    }
    else {
      let queueEmbed = new MessageEmbed()
        .setColor('#1a2b3c')
        .setTitle('Song added to queue')
        .setDescription(`[${song.title}](https://www.youtube.com/watch?v=${song.id})\n**Duration:** ${this.timeString(song.duration)}\n**Requested by:** ${msg.author}`)
        .setImage(song.thumbnail);
      
      msg.channel.send(queueEmbed);
      queue.songs.push(song);
    }
  }

  async play(msg) {
    // Get guild's queue from our map
    const queue = await this.queue.get(msg.guild.id);
    // If the queue has no songs left, leave voicechannel and delete guilds queue
    if (queue.songs.length === 0) {
      queue.voiceChannel.leave();
      return this.queue.delete(msg.guild.id);
    }
    // Create embed that we'll post to the channel
    const playEmbed = new MessageEmbed()
      .setColor('#1a2b3c')
      .setTitle('Now playing')
      .setDescription(`[${queue.songs[0].title}](https://www.youtube.com/watch?v=${queue.songs[0].id})\n**Duration:** ${this.timeString(queue.songs[0].duration)}\n**Requested by:** ${queue.songs[0].requester}`)
      .setImage(queue.songs[0].thumbnail);
    // Try playing the song
    try {
      const dispatcher = queue.connection.play(await ytdl(queue.songs[0].url), { type: 'opus', volume: 0.2 })
        .on('end', reason => {
          if (reason === 'skipped') msg.channel.send('Song skipped');
          queue.songs.shift();
          this.play(msg);
        })
        .on('error', error => {
          queue.songs.shift();
          this.client.logger.error('Dispatcher error', error);
          this.play(msg);
        });
    }
    catch(error) {
      msg.reply('Could not play song');
      this.client.logger.error(error);
      queue.songs.shift();
      this.play(msg);
    }
    // Send the embed to the channel
    msg.channel.send(playEmbed);
  }

  timeString(seconds, forceHours = false) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor(seconds % 3600 / 60);

    return `${forceHours || hours >= 1 ? `${hours}:` : ''}${hours >= 1 ? `0${minutes}`.slice(-2) : minutes}:${`0${Math.floor(seconds % 60)}`.slice(-2)}`;
  }
}
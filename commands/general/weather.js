const { Command } = require('discord.js-commando');
const { MessageEmbed } = require('discord.js');
const weather = require('weather-js');

module.exports = class WeatherCommand extends Command {
  constructor(client) {
    super(client, {
      name: 'weather',
      group: 'general',
      memberName: 'weather',
      description: 'Get weather for a location',
      clientPermissions: ['SEND_MESSAGES'],
      args: [
        {
          key: 'location',
          prompt: 'Enter location',
          type: 'string'
        }
      ]
    });
  }

  run(msg, { location }) {
    msg.channel.send('Searching weather...').then(message => {
      weather.find({
        search: location,
        degreeType: 'C'
      }, (error, result) => {
        if (error) return [this.client.logger.error('Weather API error', error),message.edit('Weather API error')];
        if (result.length < 1) return message.edit('No weather data for location');

        const embed = new MessageEmbed()
          .setColor('#00ff1d')
          .setTitle(`Weather in ${result[0].location.name}`)
          .addField('Temperature', `${result[0].current.temperature}°C (Feels like ${result[0].current.feelslike}°C)`, true)
          .addField('Sky Text', `${result[0].current.skytext}`, true)
          .addField('Wind Speed', `${result[0].current.windspeed}`, true)
          .setThumbnail(result[0].current.imageUrl);

        message.edit('', embed);
      })
    });

    return undefined;
  }
}
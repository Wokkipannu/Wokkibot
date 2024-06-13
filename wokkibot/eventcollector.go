package wokkibot

import (
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/snowflake/v2"
)

type EventCollector struct {
	client    bot.Client
	channelID snowflake.ID
	eventCh   chan *events.MessageCreate
	stopCh    chan struct{}
}

func NewEventCollector(client bot.Client, channelID snowflake.ID) *EventCollector {
	ec := &EventCollector{
		client:    client,
		channelID: channelID,
		eventCh:   make(chan *events.MessageCreate),
		stopCh:    make(chan struct{}),
	}
	client.EventManager().AddEventListeners(ec)
	return ec
}

func (ec *EventCollector) HandleEvent(event *events.MessageCreate) {
	if event.ChannelID == ec.channelID && event.Message.Content != "" {
		select {
		case ec.eventCh <- event:
		case <-ec.stopCh:
		}
	}
}

func (ec *EventCollector) OnEvent(event bot.Event) {
	switch e := event.(type) {
	case *events.MessageCreate:
		ec.HandleEvent(e)
	}
}

func (ec *EventCollector) Stop() {
	close(ec.stopCh)
	close(ec.eventCh)
	ec.client.EventManager().RemoveEventListeners(ec)
}

func (ec *EventCollector) Events() <-chan *events.MessageCreate {
	return ec.eventCh
}

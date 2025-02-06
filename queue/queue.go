package queue

import (
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Queue struct {
	Tracks []lavalink.Track
}

type QueueManager struct {
	queues map[snowflake.ID]*Queue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{
		queues: make(map[snowflake.ID]*Queue),
	}
}

func (q *Queue) Add(track ...lavalink.Track) {
	q.Tracks = append(q.Tracks, track...)
}

func (q *Queue) Next() (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	track := q.Tracks[0]
	q.Tracks = q.Tracks[1:]
	return track, true
}

func (q *Queue) Skip() (lavalink.Track, bool) {
	if len(q.Tracks) == 0 {
		return lavalink.Track{}, false
	}
	q.Tracks = q.Tracks[1:]
	return q.Tracks[0], true
}

func (q *Queue) Clear() {
	q.Tracks = make([]lavalink.Track, 0)
}

func (q *QueueManager) Get(guildID snowflake.ID) *Queue {
	queue, ok := q.queues[guildID]
	if !ok {
		queue = &Queue{
			Tracks: make([]lavalink.Track, 0),
		}
		q.queues[guildID] = queue
	}
	return queue
}

func (q *QueueManager) Delete(guildID snowflake.ID) {
	delete(q.queues, guildID)
}

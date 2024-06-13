package utils

import (
	"fmt"

	"github.com/disgoorg/disgolink/v3/lavalink"
)

// Format a duration in minutes and seconds
func FormatDuration(duration lavalink.Duration) string {
	if duration == 0 {
		return "0 minutes 0 seconds"
	}

	minutes := duration.Minutes()
	seconds := duration.SecondsPart()

	minutesText := "minute"
	secondsText := "second"
	if minutes > 1 {
		minutesText += "s"
	}
	if seconds == 1 {
		secondsText += "s"
	}

	return fmt.Sprintf("%d %s %d %s", minutes, minutesText, seconds, secondsText)
}

// Format a position in minutes and seconds
func FormatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}

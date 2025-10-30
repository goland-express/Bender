package player

import (
	"bender/internal/youtube"
	"errors"
	"io"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	errNotStreamableTrack = "track has no stream yet"
)

type track struct {
	mu             sync.RWMutex
	s              io.ReadCloser
	metadata       youtube.Metadata
	guildID        string
	voiceChannelID string
	channelID      string
	authorID       string
	msgRef         *discordgo.MessageReference
}

func (t *track) Read(buf []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.s == nil {
		return 0, errors.New(errNotStreamableTrack)
	}

	n, err := t.s.Read(buf)

	return n, err
}

func (t *track) Close() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.s == nil {
		return errors.New(errNotStreamableTrack)
	}

	err := t.s.Close()

	return err
}

func (t *track) stream() io.ReadCloser {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.s
}

func (t *track) setStream(stream io.ReadCloser) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.s = stream
}

package player

import (
	"errors"
	"sync"

	"github.com/bwmarrin/discordgo"
)

const (
	errPlayerNotInitialized = "player is not initialized for this guild"
)

type playerStatus int
type playerEvent int

const (
	statusIddle playerStatus = iota
	statusPlaying
	statusPaused
)

const (
	eventPlay playerEvent = iota
	eventNext
	eventStop
)

var players = playerList{items: make(map[string]*player)}

type player struct {
	mu sync.RWMutex

	ct *track
	q  *queue
	vc *discordgo.VoiceConnection
	s  playerStatus

	eventChannel chan playerEvent
}

func Init(s *discordgo.Session, guildID string) {
	players.setPlayer(guildID, newPlayer())

	go handler(s, guildID)
}

func Deinit(s *discordgo.Session, guildID string) {
}

func newPlayer() *player {
	return &player{
		q:            newQueue(),
		eventChannel: make(chan playerEvent),
	}
}

type playerList struct {
	mu    sync.RWMutex
	items map[string]*player
}

func (pl *playerList) player(guildID string) (*player, error) {
	pl.mu.RLock()
	player, ok := pl.items[guildID]
	pl.mu.RUnlock()

	if !ok {
		return nil, errors.New(errPlayerNotInitialized)
	}

	return player, nil
}

func (pl *playerList) setPlayer(guildID string, player *player) {
	pl.mu.Lock()
	defer pl.mu.Unlock()

	pl.items[guildID] = player
}

func (p *player) voiceConnection() *discordgo.VoiceConnection {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.vc
}

func (p *player) setVoiceConnection(vc *discordgo.VoiceConnection) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.vc = vc
}

func (p *player) currentTrack() *track {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.ct
}

func (p *player) next() (*track, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	track, ok := p.q.dequeue()

	if !ok {
		p.ct = nil

		return nil, false
	}

	p.ct = track

	return track, true
}

func (p *player) stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.q.clear()
	p.ct = nil
}

func (p *player) play(t *track) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.q.isEmpty() {
		p.ct = t
		return
	}

	p.q.enqueue(t)
}

func (p *player) status() playerStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.s
}

func (p *player) setStatus(status playerStatus) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.s = status
}

func (p *player) hasNext() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.q.isEmpty()
}

func (p *player) hasTrack() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.ct != nil
}

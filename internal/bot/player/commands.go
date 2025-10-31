package player

import (
	"bender/internal/bot/commander"
	"bender/internal/youtube"
	"io"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type PlayQuery struct {
	GuildID   string
	ChannelID string
	UserID    string
	Query     string
	MsgRef    *discordgo.MessageReference
}

func CommandPlay(ctx *commander.Context) {
	s := ctx.Session
	msgRef := ctx.Messenger.RootMessage().Reference()

	query := strings.Join(ctx.Args, " ")
	player, err := players.player(ctx.GuildID)

	if err != nil {
		log.Println("Error getting player:", err)
		return
	}

	vc := player.voiceConnection()
	userVoiceState, err := s.State.VoiceState(ctx.GuildID, ctx.AuthorID)

	if err != nil {
		ctx.Messenger.Reply("You need to be in a voice channel.")
		return
	}

	if vc != nil && userVoiceState == nil {
		ctx.Messenger.Reply("You need to connect to a voice channel first.")
		return
	}

	if vc != nil && userVoiceState.ChannelID != vc.ChannelID {
		ctx.Messenger.Reply("You need to connect to the same channel as me.")
		return
	}

	if vc == nil {
		vc, err = s.ChannelVoiceJoin(ctx.GuildID, userVoiceState.ChannelID, false, true)

		if err != nil {
			log.Println("Error joining voice channel:", err)
			return
		}

		player.setVoiceConnection(vc)
	}

	var metadata *youtube.Metadata
	var stream io.ReadCloser

	if !player.hasTrack() {
		stream, metadata, err = youtube.FetchStreamWithMetadata(query)
	} else {
		metadata, err = youtube.FetchMetadata(query)
	}

	if err != nil {
		ctx.Messenger.Reply("It was not possible to fetch this song.")
		log.Println("Error fetching from youtube:", err)
		return
	}

	if player.hasTrack() {
		ctx.Messenger.Reply("**%s** was added to playlist.", metadata.Title)
	}

	track := &track{
		s:              stream,
		guildID:        ctx.GuildID,
		voiceChannelID: userVoiceState.ChannelID,
		authorID:       ctx.AuthorID,
		channelID:      ctx.ChannelID,
		msgRef:         msgRef,
		metadata:       *metadata,
	}

	player.play(track)
}

func CommandStop(ctx *commander.Context) {
	p, err := players.player(ctx.GuildID)

	if err != nil {
		return
	}

	p.eventChannel <- eventStop
}

func CommandSkip(ctx *commander.Context) {
	p, err := players.player(ctx.GuildID)

	if err != nil {
		return
	}

	p.eventChannel <- eventNext
}

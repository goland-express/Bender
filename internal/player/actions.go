package player

import (
	"bender/internal/youtube"
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

type PlayQuery struct {
	GuildID   string
	ChannelID string
	UserID    string
	Query     string
	MsgRef    *discordgo.MessageReference
}

func Play(s *discordgo.Session, q PlayQuery) error {
	player, err := players.player(q.GuildID)

	if err != nil {
		return err
	}

	vc := player.voiceConnection()
	userVoiceState, err := s.State.VoiceState(q.GuildID, q.UserID)

	if err != nil {
		return err
	}

	if vc != nil && userVoiceState == nil {
		return errors.New("user is not connected to a voice channel")
	}

	if vc != nil && userVoiceState.ChannelID != vc.ChannelID {
		return errors.New("user is not connected to the same voice channel")
	}

	if vc == nil {
		vc, err = s.ChannelVoiceJoin(q.GuildID, userVoiceState.ChannelID, false, true)

		if err != nil {
			return err
		}

		player.setVoiceConnection(vc)
	}

	var metadata *youtube.Metadata
	var stream io.ReadCloser

	startTime := time.Now()

	if !player.hasTrack() {
		stream, metadata, err = youtube.FetchStreamWithMetadata(q.Query)
	} else {
		metadata, err = youtube.FetchMetadata(q.Query)
	}

	if err != nil {
		return err
	}

	elapsed := time.Since(startTime)

	if _, err := s.ChannelMessageSendReply(q.ChannelID, fmt.Sprintf("elapsed %f seconds", elapsed.Seconds()), q.MsgRef); err != nil {
		panic(err)
	}

	if player.hasTrack() {
		msgContent := fmt.Sprintf("**%s** was added to playlist.", metadata.Title)

		if _, err := s.ChannelMessageSendReply(q.ChannelID, msgContent, q.MsgRef); err != nil {
			log.Println("Error sending song added to playlist message: ", err)
		}
	}

	track := &track{
		s:              stream,
		guildID:        q.GuildID,
		voiceChannelID: userVoiceState.ChannelID,
		authorID:       q.UserID,
		channelID:      q.ChannelID,
		msgRef:         q.MsgRef,
		metadata:       *metadata,
	}

	player.play(track)

	return nil
}

func Stop(guildID string) error {
	p, err := players.player(guildID)

	if err != nil {
		return err
	}

	p.eventChannel <- eventStop

	return nil
}

func Skip(guildID string) error {
	p, err := players.player(guildID)

	if err != nil {
		return err
	}

	p.eventChannel <- eventNext

	return nil
}

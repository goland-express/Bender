package player

import (
	"bender/internal/youtube"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/pion/opus/pkg/oggreader"
)

func handler(s *discordgo.Session, guildID string) {
	var err error

	p, err := players.player(guildID)

	if err != nil {
		log.Printf("Error getting player for guild %s: %s\n", guildID, err)
	}

	for {
		currentTrack := p.currentTrack()

		if currentTrack == nil && p.hasNext() {
			currentTrack, _ = p.next()
		}

		if currentTrack == nil {
			time.Sleep(time.Millisecond * 500)
			continue
		}

		vc := p.voiceConnection()

		if vc == nil {
			vc, err = s.ChannelVoiceJoin(guildID, currentTrack.voiceChannelID, false, true)

			if err != nil {
				log.Println("Error entering voice channel on player handler: ", err)
				continue
			}

			p.setVoiceConnection(vc)
		}

		if currentTrack.stream() == nil {
			stream, err := youtube.FetchStream(currentTrack.metadata.Url)

			if err != nil {
				if _, err := s.ChannelMessageSendReply(currentTrack.channelID, "It was not possible to fetch the stream from current track.", currentTrack.msgRef); err != nil {
					log.Println("Error sending current playing message with reply: ", err)
				}

				p.next()
				log.Println("Error fetching stream from youtube: ", err)

				continue
			}

			currentTrack.setStream(stream)
		}

		if err != nil {
			log.Println("Error fetching a stream from youtube: ", err)
			continue
		}

		oggStream, _, err := oggreader.NewWith(currentTrack)

		if err != nil {
			log.Println("Error initializing oggreader: ", err)
			continue
		}

		msgContent := fmt.Sprintf("Playing **%s** (%.0f seconds)", currentTrack.metadata.Title, currentTrack.metadata.Duration)

		if _, err := s.ChannelMessageSendReply(currentTrack.channelID, msgContent, currentTrack.msgRef); err != nil {
			log.Println("Error sending current playing message with reply: ", err)
		}

	oggReader:
		for {
			segments, _, err := oggStream.ParseNextPage()

			if err == io.EOF {
				_, ok := p.next()

				if !ok {
					msgContent := "Playlist has no more songs."

					if _, err := s.ChannelMessageSend(currentTrack.channelID, msgContent); err != nil {
						log.Println("Error sending current playing message: ", err)
					}
				}

				break
			}

			if err != nil {
				log.Println("Error parsing ogg page: ", err)
				break
			}

			for _, segment := range segments {
				select {
				case event := <-p.eventChannel:
					{
						if event == eventNext {
							currentTrack.Close()
							p.next()

							break oggReader
						}

						if event == eventStop {
							currentTrack.Close()
							p.stop()

							break oggReader
						}
					}
				default:
				}

				vc.OpusSend <- segment
			}
		}
	}
}

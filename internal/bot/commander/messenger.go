package commander

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Messenger struct {
	session     *discordgo.Session
	rootMessage *discordgo.Message
}

func NewMessenger(session *discordgo.Session, rootMessage *discordgo.Message) *Messenger {
	return &Messenger{
		session:     session,
		rootMessage: rootMessage,
	}
}

func (m *Messenger) Send(format string, a ...any) (*discordgo.Message, error) {
	content := fmt.Sprintf(format, a...)

	msg, err := m.session.ChannelMessageSend(m.rootMessage.ChannelID, content)

	if err != nil {
		log.Printf("Error sending message via commander messenger: %s (%s)", err, content)
		return nil, err
	}

	return msg, nil
}

func (m *Messenger) Reply(format string, a ...any) (*discordgo.Message, error) {
	content := fmt.Sprintf(format, a...)

	msg, err := m.session.ChannelMessageSendReply(m.rootMessage.ChannelID, content, m.rootMessage.Reference())

	if err != nil {
		log.Printf("Error sending message reply via commander messenger: %s (%s)", err, content)
		return nil, err
	}

	return msg, nil
}

func (m *Messenger) RootMessage() *discordgo.Message {
	return m.rootMessage
}

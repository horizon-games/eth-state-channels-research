package matcher

import (
	"encoding/json"
	"log"
)

const (
	UNSUBSCRIBE Code = -100 // a internal PubSubManager message type
)

type PubSubManager struct {
	Channels map[string]chan *Message // subkey -> listening channel
	*SessionManager
}

func NewPubSubManager(smgr *SessionManager) *PubSubManager {
	return &PubSubManager{
		SessionManager: smgr,
		Channels:       make(map[string]chan *Message),
	}
}

func (mgr *PubSubManager) Publish(channel string, msg Message) error {
	count, err := mgr.RedisClient.Publish(channel, msg).Result()
	if err != nil {
		return err
	}
	log.Printf("%d messages broadcast to channel %s", count, channel)
	return nil
}

func (mgr *PubSubManager) Subscribe(key string, channel chan *Message) {
	if mgr.Channels[key] != nil {
		return
	}
	mgr.Channels[key] = channel
	go func() {
		ps := mgr.RedisClient.Subscribe(key)
		defer ps.Close()
		for {
			message, err := ps.ReceiveMessage()
			if err != nil {
				log.Printf("Error receiving redis message or timeout, continuing: %s", err.Error())
				continue
			}
			msg := &Message{}
			err = json.Unmarshal([]byte(message.Payload), msg)
			if err != nil {
				log.Printf("Error unmarshalling subkey channel message payload %s %s", message.Payload, err.Error())
				continue
			}
			if msg.Code == UNSUBSCRIBE || msg.Code == TERMINATE {
				ps.Unsubscribe(key)
				delete(mgr.Channels, key)
				if msg.Code == TERMINATE { // bubble up termination message so others can take action
					mgr.Channels[key] <- msg
				}
				break
			}
			mgr.Channels[key] <- msg
		}
	}()
}

func (mgr *PubSubManager) Unsubscribe(key string) {
	mgr.Publish(key, Message{Meta: &Meta{Code: UNSUBSCRIBE}})
}

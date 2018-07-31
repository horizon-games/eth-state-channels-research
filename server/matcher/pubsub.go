package matcher

import (
	"encoding/json"
	"log"

	"github.com/horizon-games/arcadeum/server/lib/util"
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
	log.Printf("Attempting to send message %v to channel %v\n", msg, channel)
	msgJson, err := util.Jsonify(msg)
	if err != nil {
		log.Printf("Error jsonifying message %s\n", err.Error())
		return err
	}
	count, err := mgr.RedisClient.Publish(channel, msgJson).Result()
	if err != nil {
		log.Printf("Error sending message to channel %v, %v\n", channel, err.Error())
		return err
	}
	log.Printf("%d messages broadcast to channel %v\n", count, channel)
	return nil
}

func (mgr *PubSubManager) Subscribe(key string, channel chan *Message) {
	if mgr.Channels[key] != nil {
		return
	}
	mgr.Channels[key] = channel
	go func() {
		log.Printf("Subscribing to channel %s", key)
		ps := mgr.RedisClient.Subscribe(key)
		defer ps.Close()
		for {
			log.Printf("Waiting for message on channel %s", key)
			message, err := ps.ReceiveMessage()
			log.Printf("Received message %s on channel %s", message, key)
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

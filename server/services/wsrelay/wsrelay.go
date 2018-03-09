package wsrelay

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/websocket"
	"github.com/horizon-games/arcadeum/server/services/arcadeum"
	"github.com/horizon-games/arcadeum/server/services/matcher"
	"github.com/horizon-games/arcadeum/server/services/wsrelay/config"
)

type Server struct {
	Matcher *matcher.Service
}

type MessageRequest struct {
	PlayerConn *websocket.Conn // sender
	*matcher.Message
}

var relay = make(chan *MessageRequest)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow cross origin connections
	},
}

func NewServer(cfg *config.Config) *Server {
	return &Server{
		Matcher: matcher.NewService(&cfg.ENV, &cfg.MatcherConfig, &cfg.ETHConfig, &cfg.ArcadeumConfig),
	}
}

func (s *Server) HandleMessages() {
	for {
		// Grab the next message from the relay channel
		msg := <-relay

		// Grab game and player session info
		session := s.Matcher.FindSession(msg.MatchID)
		player := session.FindPlayer(msg.PlayerConn)
		if session != nil || player != nil {
			msg.PlayerConn.WriteJSON(matcher.NewError(fmt.Sprintf("Unknown match ID %d", msg.MatchID)))
			continue
		}

		if msg.Code == matcher.SIGNED_TIMESTAMP { // Verified signed timestamp
			req := &arcadeum.VerifyTimestampRequest{}
			err := json.Unmarshal([]byte(msg.Payload), req)
			if err != nil {
				msg.PlayerConn.WriteJSON(matcher.NewError(err.Error()))
				continue
			}
			verified, err := s.Matcher.VerifyTimestamp(session.GameID, session.MatchID, req, player)
			if err != nil {
				msg.PlayerConn.WriteJSON(matcher.NewError(err.Error()))
				continue
			}
			if !verified {
				msg.PlayerConn.WriteJSON(matcher.NewError("Timestamp signature not verified."))
				continue
			}
			player.Verified = verified
			player.TimestampSig = req.Signature // set the verified signature
			err2 := s.Matcher.BeginVerifiedMatch(session)
			if err2 != nil {
				msg.PlayerConn.WriteJSON(matcher.NewError(fmt.Sprintf("Error sending begin match payload. %s", err2.Error())))
				continue
			}
		} else if !session.IsVerified() { // don't relay messages unless both players have proved their timestamp signature
			msg.PlayerConn.WriteJSON(matcher.NewError("Match session not verified."))
		} else { // verified
			opponent := session.GetOpponent(msg.PlayerConn)
			if opponent != nil {
				opponent.Conn.WriteJSON(msg.Message) // relay the message to the opponent
			}
		}

	}
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	log.Println("Opening WS connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	request := &matcher.MatchRequest{
		Conn:  ws,
		Token: matcher.Context(r),
	}
	go s.Matcher.FindMatch(request)

	for {
		var msg matcher.Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			ws.WriteJSON(matcher.NewError("Unrecognized message format."))
		} else {
			// relay the newly received message to the opponent
			relay <- &MessageRequest{PlayerConn: ws, Message: &msg}
		}
	}
}

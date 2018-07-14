package server

import (
	"log"
	"net/http"

	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/websocket"
	"github.com/horizon-games/arcadeum/server/config"
	"github.com/horizon-games/arcadeum/server/matcher"
	"github.com/pkg/errors"
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

func New(cfg *config.Config) (*Server, error) {
	if cfg.ENV.WorkingDir == "" {
		return nil, errors.New("Missing working directory in config")
	}
	return &Server{
		Matcher: matcher.NewService(
			&cfg.ENV,
			&cfg.MatcherConfig,
			&cfg.ETHConfig,
			&cfg.ArcadeumConfig,
			&cfg.RedisConfig),
	}, nil
}

func (s *Server) Start() error {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`¯\_(ツ)_/¯`))
	})

	// Configure websocket route
	r.With(matcher.AddTokenContext).HandleFunc("/ws", s.HandleConnections)

	go s.HandleMessages()
	go s.Matcher.HandleMatchResponses()

	// Start the server on localhost
	log.Printf("ARCADEUM Server started :%d; connect at /ws", s.Matcher.ENV.Port)

	if s.Matcher.ENV.TLSEnabled {
		return http.ListenAndServeTLS(fmt.Sprintf(":%d", s.Matcher.ENV.Port),
			s.Matcher.ENV.TLSCertFile, s.Matcher.ENV.TLSKeyFile, r)
	} else {
		return http.ListenAndServe(fmt.Sprintf(":%d", s.Matcher.ENV.Port), r)
	}

	return nil
}

func (s *Server) HandleMessages() {
	for {
		msg := <-relay
		err := s.Matcher.OnMessage(msg.Message)
		if err != nil {
			msg.PlayerConn.WriteJSON(matcher.NewError(err.Error()))
		}
	}
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	log.Println("Opening WS connection")
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer ws.Close()
		s.FindMatch(matcher.Context(r), ws)

		for {
			var msg matcher.Message
			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				if websocket.IsCloseError(err, 1001) {
					break
				}
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					break
				}
				log.Printf("client is gone")
				return
			} else {
				relay <- &MessageRequest{PlayerConn: ws, Message: &msg}
			}
		}
	}()
}

func (s *Server) FindMatch(token *matcher.Token, conn *websocket.Conn) {
	channel := make(chan *matcher.Message)
	s.Matcher.SubscribeToSubKey(token.SubKey, channel)
	go OnMessage(conn, channel)
	go s.Matcher.FindMatch(token)
}

func OnMessage(ws *websocket.Conn, messages chan *matcher.Message) {
	for {
		msg := <-messages
		log.Printf("GOT PUBLISHED MESSAGE, sending to client: %s", msg)
		err := ws.WriteJSON(msg)
		if err != nil {
			log.Printf("Error sending message to client over websocket %s", err.Error())
		}
		if msg.Code == matcher.TERMINATE {
			break
		}
	}
}

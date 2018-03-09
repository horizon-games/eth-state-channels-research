package matcher

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/horizon-games/arcadeum/server/services/crypto"
)

// Request token.
// This token is returned as a base64 string by a client requesting to play a game.
type Token struct {
	GameID          uint32           `json:"gameID"`        // globally unique game ID
	SubKey          common.Address   `json:"subkey,string"` // public address of account owner of subkey
	SubKeySignature crypto.Signature `json:"signature"`     // Signed signature of SubKey to prove SubKey ownership
	Seed            []byte           `json:"seed,string"`   // game "deck"; unmarshalled base64; never share with opponents!
}

// Unmarshal token from query parameter and put it in request context for later retrieval and authentication logic.
func AddTokenContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Adding token context")
		tokenb64 := r.FormValue("token")
		if len(tokenb64) == 0 {
			writeUnauthorized(w)
			return
		}
		decoded, err := base64.StdEncoding.DecodeString(tokenb64)
		if err != nil {
			log.Printf("decode error: %s", err.Error())
			writeUnauthorized(w)
			return
		}
		token := &Token{}
		er := json.Unmarshal(decoded, &token)
		if er != nil {
			log.Printf("parse error: %s", er.Error())
			writeUnauthorized(w)
			return
		}
		ctx := context.WithValue(
			r.Context(),
			"Token",
			token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func writeUnauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Unauthorized or invalid authorization token."))
}

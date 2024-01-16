package token

import "time"

//Maker represents the issuer of the token, which issues and verifies tokens.
type Maker interface {
	// Make produces a new token with the given duration for the given username
	Make(username string, duration time.Duration) (string, error)

	//Verify returns the payload of the token if its valid or an error otherwise
	Verify(token string) (*Payload, error)
}
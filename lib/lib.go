package lib

import "github.com/rs/xid"

// GenerateXID generate unique id
func GenerateXID() string {
	return xid.New().String()
}

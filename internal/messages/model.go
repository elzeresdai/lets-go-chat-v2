package messages

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Message struct {
	Message string    `json:"message"`
	UserId  uuid.UUID `json:"userId"`
}

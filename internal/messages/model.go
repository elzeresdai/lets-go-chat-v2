package messages

import (
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
)

type Message struct {
	ID      string `json:"id"`
	message string `json:"message"`
	userId  uuid.UUID
}

type MessageRequest struct {
	userId  uuid.UUID
	message string
}

package domain

import "time"

type Reply struct {
	ReplyType  string    `json:"replyType"`
	SagaID     string    `json:"sagaId"`
	Success    bool      `json:"success"`
	Reason     string    `json:"reason,omitempty"`
	OccurredAt time.Time `json:"occurredAt"`
}

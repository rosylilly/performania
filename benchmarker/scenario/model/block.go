package model

type Block struct {
	BlockerID string `json:"blocker_id"`
	BlockedID string `json:"blocked_id"`
}

package events

import (
	"net/mail"
	"time"
)

type Data struct {
	// Database ID
	ID string
	// Message ID
	MessageID string
	// Subject
	Subject string
	// From address
	From *mail.Address
	// To addresses
	To []*mail.Address
	// Cc addresses
	Cc []*mail.Address
	// Bcc addresses
	Bcc []*mail.Address
	// ReplyTo addresses
	ReplyTo []*mail.Address
	// // List-Unsubscribe header information
	// // swagger:ignore
	// ListUnsubscribe ListUnsubscribe
	// // Message date if set, else date received
	Date time.Time
	// Message tags
	Tags []string
	// Message body text
	Text string
	// Message body HTML
	HTML string
	// Message size in bytes
	Size int
	// // Inline message attachments
	// Inline []Attachment
	// // Message attachments
	// Attachments []Attachment
}

type MercureMessage struct {
	ID    int    `json:"mail_id"`
	To    string `json:"to"`
	From  string `json:"from"`
	State string `json:"state"`
	Raw   string `json:"raw"`
	Data  Data   `json:"Data"`
}

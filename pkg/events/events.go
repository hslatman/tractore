package events

// TODO: add some type identifier, so that the frontend can distinguish
// between the types published to Mercure? Or use Mercure specific events,
// which then get serialized incl. their type?

type IncomingEmail struct {
	ID   int    `json:"incoming_mail_id"`
	To   string `json:"to"`
	From string `json:"from"`
	Raw  string `json:"raw"`
}

type OutgoingEmail struct {
	ID   int    `json:"incoming_mail_id"`
	To   string `json:"to"`
	From string `json:"from"`
	Raw  string `json:"raw"`
}

type SentEmail struct {
	ID   int    `json:"incoming_mail_id"`
	To   string `json:"to"`
	From string `json:"from"`
	Raw  string `json:"raw"`
}

type TrackEmail struct {
	ID int `json:"incoming_mail_id"`
}

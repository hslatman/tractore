package events

type MercureMessage struct {
	ID    int    `json:"mail_id"`
	To    string `json:"to"`
	From  string `json:"from"`
	State string `json:"state"`
	Raw   string `json:"raw"`
}

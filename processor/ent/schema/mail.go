package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Mail holds the schema definition for the Mail entity.
type Mail struct {
	ent.Schema
}

// Fields of the Mail.
func (Mail) Fields() []ent.Field {
	return []ent.Field{
		field.String("to"),
		field.String("from"),
		field.Bytes("raw"),
		field.Int("incoming_mail_id"),
	}
}

// Edges of the Mail.
func (Mail) Edges() []ent.Edge {
	return nil
}

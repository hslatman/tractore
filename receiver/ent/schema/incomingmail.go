package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// IncomingMail holds the schema definition for the IncomingMail entity.
type IncomingMail struct {
	ent.Schema
}

// Fields of the IncomingMail.
func (IncomingMail) Fields() []ent.Field {
	return []ent.Field{
		field.String("to"),
		field.String("from"),
		field.Bytes("raw"),
	}
}

// Edges of the IncomingMail.
func (IncomingMail) Edges() []ent.Edge {
	return nil
}

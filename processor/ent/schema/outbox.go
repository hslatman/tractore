package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Outbox holds the schema definition for the Outbox entity.
type Outbox struct {
	ent.Schema
}

// Annotations of the User.
func (Outbox) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "outbox"},
	}
}

// Fields of the Outbox.
func (Outbox) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id"),
		field.Text("topic"),
		field.Bytes("data"),
		field.Time("inserted_at"),
	}
}

// Edges of the Outbox.
func (Outbox) Edges() []ent.Edge {
	return nil
}

func (Outbox) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("topic", "id"),
	}
}

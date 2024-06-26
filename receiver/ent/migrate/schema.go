// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// IncomingMailsColumns holds the columns for the "incoming_mails" table.
	IncomingMailsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt, Increment: true},
		{Name: "to", Type: field.TypeString},
		{Name: "from", Type: field.TypeString},
		{Name: "raw", Type: field.TypeBytes},
	}
	// IncomingMailsTable holds the schema information for the "incoming_mails" table.
	IncomingMailsTable = &schema.Table{
		Name:       "incoming_mails",
		Columns:    IncomingMailsColumns,
		PrimaryKey: []*schema.Column{IncomingMailsColumns[0]},
	}
	// OutboxColumns holds the columns for the "outbox" table.
	OutboxColumns = []*schema.Column{
		{Name: "id", Type: field.TypeInt64, Increment: true},
		{Name: "topic", Type: field.TypeString, Size: 2147483647},
		{Name: "data", Type: field.TypeBytes},
		{Name: "inserted_at", Type: field.TypeTime},
	}
	// OutboxTable holds the schema information for the "outbox" table.
	OutboxTable = &schema.Table{
		Name:       "outbox",
		Columns:    OutboxColumns,
		PrimaryKey: []*schema.Column{OutboxColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "outbox_topic_id",
				Unique:  false,
				Columns: []*schema.Column{OutboxColumns[1], OutboxColumns[0]},
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		IncomingMailsTable,
		OutboxTable,
	}
)

func init() {
	OutboxTable.Annotation = &entsql.Annotation{
		Table: "outbox",
	}
}

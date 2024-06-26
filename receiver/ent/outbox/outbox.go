// Code generated by ent, DO NOT EDIT.

package outbox

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the outbox type in the database.
	Label = "outbox"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldTopic holds the string denoting the topic field in the database.
	FieldTopic = "topic"
	// FieldData holds the string denoting the data field in the database.
	FieldData = "data"
	// FieldInsertedAt holds the string denoting the inserted_at field in the database.
	FieldInsertedAt = "inserted_at"
	// Table holds the table name of the outbox in the database.
	Table = "outbox"
)

// Columns holds all SQL columns for outbox fields.
var Columns = []string{
	FieldID,
	FieldTopic,
	FieldData,
	FieldInsertedAt,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the Outbox queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByTopic orders the results by the topic field.
func ByTopic(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTopic, opts...).ToFunc()
}

// ByInsertedAt orders the results by the inserted_at field.
func ByInsertedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldInsertedAt, opts...).ToFunc()
}

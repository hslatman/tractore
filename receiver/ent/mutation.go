// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"encore.app/receiver/ent/incomingmail"
	"encore.app/receiver/ent/outbox"
	"encore.app/receiver/ent/predicate"
	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeIncomingMail = "IncomingMail"
	TypeOutbox       = "Outbox"
)

// IncomingMailMutation represents an operation that mutates the IncomingMail nodes in the graph.
type IncomingMailMutation struct {
	config
	op            Op
	typ           string
	id            *int
	to            *string
	from          *string
	raw           *[]byte
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*IncomingMail, error)
	predicates    []predicate.IncomingMail
}

var _ ent.Mutation = (*IncomingMailMutation)(nil)

// incomingmailOption allows management of the mutation configuration using functional options.
type incomingmailOption func(*IncomingMailMutation)

// newIncomingMailMutation creates new mutation for the IncomingMail entity.
func newIncomingMailMutation(c config, op Op, opts ...incomingmailOption) *IncomingMailMutation {
	m := &IncomingMailMutation{
		config:        c,
		op:            op,
		typ:           TypeIncomingMail,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withIncomingMailID sets the ID field of the mutation.
func withIncomingMailID(id int) incomingmailOption {
	return func(m *IncomingMailMutation) {
		var (
			err   error
			once  sync.Once
			value *IncomingMail
		)
		m.oldValue = func(ctx context.Context) (*IncomingMail, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().IncomingMail.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withIncomingMail sets the old IncomingMail of the mutation.
func withIncomingMail(node *IncomingMail) incomingmailOption {
	return func(m *IncomingMailMutation) {
		m.oldValue = func(context.Context) (*IncomingMail, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m IncomingMailMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m IncomingMailMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *IncomingMailMutation) ID() (id int, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *IncomingMailMutation) IDs(ctx context.Context) ([]int, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().IncomingMail.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetTo sets the "to" field.
func (m *IncomingMailMutation) SetTo(s string) {
	m.to = &s
}

// To returns the value of the "to" field in the mutation.
func (m *IncomingMailMutation) To() (r string, exists bool) {
	v := m.to
	if v == nil {
		return
	}
	return *v, true
}

// OldTo returns the old "to" field's value of the IncomingMail entity.
// If the IncomingMail object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *IncomingMailMutation) OldTo(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTo is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTo requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTo: %w", err)
	}
	return oldValue.To, nil
}

// ResetTo resets all changes to the "to" field.
func (m *IncomingMailMutation) ResetTo() {
	m.to = nil
}

// SetFrom sets the "from" field.
func (m *IncomingMailMutation) SetFrom(s string) {
	m.from = &s
}

// From returns the value of the "from" field in the mutation.
func (m *IncomingMailMutation) From() (r string, exists bool) {
	v := m.from
	if v == nil {
		return
	}
	return *v, true
}

// OldFrom returns the old "from" field's value of the IncomingMail entity.
// If the IncomingMail object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *IncomingMailMutation) OldFrom(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldFrom is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldFrom requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldFrom: %w", err)
	}
	return oldValue.From, nil
}

// ResetFrom resets all changes to the "from" field.
func (m *IncomingMailMutation) ResetFrom() {
	m.from = nil
}

// SetRaw sets the "raw" field.
func (m *IncomingMailMutation) SetRaw(b []byte) {
	m.raw = &b
}

// Raw returns the value of the "raw" field in the mutation.
func (m *IncomingMailMutation) Raw() (r []byte, exists bool) {
	v := m.raw
	if v == nil {
		return
	}
	return *v, true
}

// OldRaw returns the old "raw" field's value of the IncomingMail entity.
// If the IncomingMail object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *IncomingMailMutation) OldRaw(ctx context.Context) (v []byte, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldRaw is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldRaw requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldRaw: %w", err)
	}
	return oldValue.Raw, nil
}

// ResetRaw resets all changes to the "raw" field.
func (m *IncomingMailMutation) ResetRaw() {
	m.raw = nil
}

// Where appends a list predicates to the IncomingMailMutation builder.
func (m *IncomingMailMutation) Where(ps ...predicate.IncomingMail) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the IncomingMailMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *IncomingMailMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.IncomingMail, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *IncomingMailMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *IncomingMailMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (IncomingMail).
func (m *IncomingMailMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *IncomingMailMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.to != nil {
		fields = append(fields, incomingmail.FieldTo)
	}
	if m.from != nil {
		fields = append(fields, incomingmail.FieldFrom)
	}
	if m.raw != nil {
		fields = append(fields, incomingmail.FieldRaw)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *IncomingMailMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case incomingmail.FieldTo:
		return m.To()
	case incomingmail.FieldFrom:
		return m.From()
	case incomingmail.FieldRaw:
		return m.Raw()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *IncomingMailMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case incomingmail.FieldTo:
		return m.OldTo(ctx)
	case incomingmail.FieldFrom:
		return m.OldFrom(ctx)
	case incomingmail.FieldRaw:
		return m.OldRaw(ctx)
	}
	return nil, fmt.Errorf("unknown IncomingMail field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *IncomingMailMutation) SetField(name string, value ent.Value) error {
	switch name {
	case incomingmail.FieldTo:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTo(v)
		return nil
	case incomingmail.FieldFrom:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetFrom(v)
		return nil
	case incomingmail.FieldRaw:
		v, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetRaw(v)
		return nil
	}
	return fmt.Errorf("unknown IncomingMail field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *IncomingMailMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *IncomingMailMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *IncomingMailMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown IncomingMail numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *IncomingMailMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *IncomingMailMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *IncomingMailMutation) ClearField(name string) error {
	return fmt.Errorf("unknown IncomingMail nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *IncomingMailMutation) ResetField(name string) error {
	switch name {
	case incomingmail.FieldTo:
		m.ResetTo()
		return nil
	case incomingmail.FieldFrom:
		m.ResetFrom()
		return nil
	case incomingmail.FieldRaw:
		m.ResetRaw()
		return nil
	}
	return fmt.Errorf("unknown IncomingMail field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *IncomingMailMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *IncomingMailMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *IncomingMailMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *IncomingMailMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *IncomingMailMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *IncomingMailMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *IncomingMailMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown IncomingMail unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *IncomingMailMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown IncomingMail edge %s", name)
}

// OutboxMutation represents an operation that mutates the Outbox nodes in the graph.
type OutboxMutation struct {
	config
	op            Op
	typ           string
	id            *int64
	topic         *string
	data          *[]byte
	inserted_at   *time.Time
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*Outbox, error)
	predicates    []predicate.Outbox
}

var _ ent.Mutation = (*OutboxMutation)(nil)

// outboxOption allows management of the mutation configuration using functional options.
type outboxOption func(*OutboxMutation)

// newOutboxMutation creates new mutation for the Outbox entity.
func newOutboxMutation(c config, op Op, opts ...outboxOption) *OutboxMutation {
	m := &OutboxMutation{
		config:        c,
		op:            op,
		typ:           TypeOutbox,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withOutboxID sets the ID field of the mutation.
func withOutboxID(id int64) outboxOption {
	return func(m *OutboxMutation) {
		var (
			err   error
			once  sync.Once
			value *Outbox
		)
		m.oldValue = func(ctx context.Context) (*Outbox, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().Outbox.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withOutbox sets the old Outbox of the mutation.
func withOutbox(node *Outbox) outboxOption {
	return func(m *OutboxMutation) {
		m.oldValue = func(context.Context) (*Outbox, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m OutboxMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m OutboxMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of Outbox entities.
func (m *OutboxMutation) SetID(id int64) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *OutboxMutation) ID() (id int64, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *OutboxMutation) IDs(ctx context.Context) ([]int64, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []int64{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().Outbox.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetTopic sets the "topic" field.
func (m *OutboxMutation) SetTopic(s string) {
	m.topic = &s
}

// Topic returns the value of the "topic" field in the mutation.
func (m *OutboxMutation) Topic() (r string, exists bool) {
	v := m.topic
	if v == nil {
		return
	}
	return *v, true
}

// OldTopic returns the old "topic" field's value of the Outbox entity.
// If the Outbox object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OutboxMutation) OldTopic(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTopic is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTopic requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTopic: %w", err)
	}
	return oldValue.Topic, nil
}

// ResetTopic resets all changes to the "topic" field.
func (m *OutboxMutation) ResetTopic() {
	m.topic = nil
}

// SetData sets the "data" field.
func (m *OutboxMutation) SetData(b []byte) {
	m.data = &b
}

// Data returns the value of the "data" field in the mutation.
func (m *OutboxMutation) Data() (r []byte, exists bool) {
	v := m.data
	if v == nil {
		return
	}
	return *v, true
}

// OldData returns the old "data" field's value of the Outbox entity.
// If the Outbox object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OutboxMutation) OldData(ctx context.Context) (v []byte, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldData is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldData requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldData: %w", err)
	}
	return oldValue.Data, nil
}

// ResetData resets all changes to the "data" field.
func (m *OutboxMutation) ResetData() {
	m.data = nil
}

// SetInsertedAt sets the "inserted_at" field.
func (m *OutboxMutation) SetInsertedAt(t time.Time) {
	m.inserted_at = &t
}

// InsertedAt returns the value of the "inserted_at" field in the mutation.
func (m *OutboxMutation) InsertedAt() (r time.Time, exists bool) {
	v := m.inserted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldInsertedAt returns the old "inserted_at" field's value of the Outbox entity.
// If the Outbox object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *OutboxMutation) OldInsertedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldInsertedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldInsertedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldInsertedAt: %w", err)
	}
	return oldValue.InsertedAt, nil
}

// ResetInsertedAt resets all changes to the "inserted_at" field.
func (m *OutboxMutation) ResetInsertedAt() {
	m.inserted_at = nil
}

// Where appends a list predicates to the OutboxMutation builder.
func (m *OutboxMutation) Where(ps ...predicate.Outbox) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the OutboxMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *OutboxMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.Outbox, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *OutboxMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *OutboxMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (Outbox).
func (m *OutboxMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *OutboxMutation) Fields() []string {
	fields := make([]string, 0, 3)
	if m.topic != nil {
		fields = append(fields, outbox.FieldTopic)
	}
	if m.data != nil {
		fields = append(fields, outbox.FieldData)
	}
	if m.inserted_at != nil {
		fields = append(fields, outbox.FieldInsertedAt)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *OutboxMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case outbox.FieldTopic:
		return m.Topic()
	case outbox.FieldData:
		return m.Data()
	case outbox.FieldInsertedAt:
		return m.InsertedAt()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *OutboxMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case outbox.FieldTopic:
		return m.OldTopic(ctx)
	case outbox.FieldData:
		return m.OldData(ctx)
	case outbox.FieldInsertedAt:
		return m.OldInsertedAt(ctx)
	}
	return nil, fmt.Errorf("unknown Outbox field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *OutboxMutation) SetField(name string, value ent.Value) error {
	switch name {
	case outbox.FieldTopic:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTopic(v)
		return nil
	case outbox.FieldData:
		v, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetData(v)
		return nil
	case outbox.FieldInsertedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetInsertedAt(v)
		return nil
	}
	return fmt.Errorf("unknown Outbox field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *OutboxMutation) AddedFields() []string {
	return nil
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *OutboxMutation) AddedField(name string) (ent.Value, bool) {
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *OutboxMutation) AddField(name string, value ent.Value) error {
	switch name {
	}
	return fmt.Errorf("unknown Outbox numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *OutboxMutation) ClearedFields() []string {
	return nil
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *OutboxMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *OutboxMutation) ClearField(name string) error {
	return fmt.Errorf("unknown Outbox nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *OutboxMutation) ResetField(name string) error {
	switch name {
	case outbox.FieldTopic:
		m.ResetTopic()
		return nil
	case outbox.FieldData:
		m.ResetData()
		return nil
	case outbox.FieldInsertedAt:
		m.ResetInsertedAt()
		return nil
	}
	return fmt.Errorf("unknown Outbox field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *OutboxMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *OutboxMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *OutboxMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *OutboxMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *OutboxMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *OutboxMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *OutboxMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown Outbox unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *OutboxMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown Outbox edge %s", name)
}

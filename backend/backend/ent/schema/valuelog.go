package schema

import "entgo.io/ent"

// ValueLog holds the schema definition for the ValueLog entity.
type ValueLog struct {
	ent.Schema
}

// Fields of the ValueLog.
func (ValueLog) Fields() []ent.Field {
	return nil
}

// Edges of the ValueLog.
func (ValueLog) Edges() []ent.Edge {
	return nil
}

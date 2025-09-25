package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ValueLog holds the schema definition for the ValueLog entity.
type ValueLog struct {
	ent.Schema
}

// Fields of the ValueLog.
func (ValueLog) Fields() []ent.Field {
	return []ent.Field{
		field.Int("value1"),
		field.Int("value2"),
		field.Time("created_at").
			Default(time.Now),
	}
}

// Edges of the ValueLog.
func (ValueLog) Edges() []ent.Edge {
	return nil
}
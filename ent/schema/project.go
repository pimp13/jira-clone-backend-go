package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Project holds the schema definition for the Project entity.
type Project struct {
	ent.Schema
}

// Fields of the Project.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New).Immutable(),

		field.String("name").NotEmpty().MinLen(3).MaxLen(195),
		field.String("image_url").Optional().Nillable(),
		field.Bool("is_active").Default(true),
		field.String("description").Optional().Nillable(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),

		// forign keys
		field.UUID("workspace_id", uuid.UUID{}),
	}
}

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("workspace", Workspace.Type).
			Ref("projects").
			Unique().
			Required().
			Field("workspace_id"),
	}
}

// Indexs of the Project.
func (Project) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("is_active", "workspace_id", "name"),
	}
}

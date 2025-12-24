package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New).Immutable(),
		field.String("email").NotEmpty().Unique(),
		field.String("name").NotEmpty().MinLen(3).MaxLen(195),
		field.String("password").MinLen(8),
		field.Bool("is_active").Default(true).Nillable().Optional(),
		field.String("avatar_url").Nillable().Optional(),
		field.Enum("role").Values("USER", "ADMIN").Default("USER"),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("workspaces", Workspace.Type),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("email").Unique(),
		index.Fields("is_active"),
	}
}

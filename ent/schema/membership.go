package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Membership holds the schema definition for the Membership entity.
type Membership struct {
	ent.Schema
}

// Fields of the Membership.
func (Membership) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).Default(uuid.New).Immutable(),

		field.Enum("role").
			Values("owner", "admin", "member", "viewer"),
		field.Enum("status").
			Values("active", "invited", "removed"),
		field.Time("joined_at").Default(time.Now).Immutable(),

		field.UUID("user_id", uuid.UUID{}),
		field.UUID("workspace_id", uuid.UUID{}),
	}
}

// Edges of the Membership.
func (Membership) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("memberships").
			Unique().
			Required().
			Field("user_id"),

		edge.From("workspace", Workspace.Type).
			Ref("memberships").
			Unique().
			Required().
			Field("workspace_id"),
	}
}

// Indexs of the Membership.
func (Membership) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("role", "status"),
		index.Fields("user_id", "workspace_id").Unique(),
	}
}

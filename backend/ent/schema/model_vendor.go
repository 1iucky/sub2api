package schema

import (
	"github.com/Wei-Shaw/sub2api/ent/schema/mixins"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// ModelVendor holds metadata for model providers/vendors.
type ModelVendor struct {
	ent.Schema
}

func (ModelVendor) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "model_vendors"},
	}
}

func (ModelVendor) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ModelVendor) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			NotEmpty().
			MaxLen(100).
			Unique(),
		field.String("provider_key").
			Default("").
			MaxLen(80),
		field.String("icon_key").
			Default("").
			MaxLen(80),
		field.String("description").
			Default(""),
		field.Int("sort_order").
			Default(0),
		field.Time("deleted_at").
			Optional().
			Nillable().
			SchemaType(map[string]string{"postgres": "timestamptz"}),
	}
}

func (ModelVendor) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("models", ModelCatalog.Type).
			Ref("vendor"),
	}
}

func (ModelVendor) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider_key"),
		index.Fields("deleted_at"),
	}
}

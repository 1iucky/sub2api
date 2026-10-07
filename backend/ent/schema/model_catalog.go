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

// ModelCatalog holds platform-wide model metadata used to enrich pricing and marketplace views.
type ModelCatalog struct {
	ent.Schema
}

func (ModelCatalog) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "model_catalogs"},
	}
}

func (ModelCatalog) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixins.TimeMixin{},
	}
}

func (ModelCatalog) Fields() []ent.Field {
	return []ent.Field{
		field.String("model_id").
			NotEmpty().
			MaxLen(200),
		field.String("normalized_model_id").
			NotEmpty().
			MaxLen(200),
		field.String("display_name").
			Default("").
			MaxLen(200),
		field.String("platform").
			Default("").
			MaxLen(50),
		field.String("provider").
			Default("").
			MaxLen(100),
		field.Int64("vendor_id").
			Optional().
			Nillable(),
		field.String("mode").
			Default("").
			MaxLen(50),
		field.String("description").
			Default(""),
		field.JSON("tags", []string{}).
			Default([]string{}),
		field.JSON("capabilities", map[string]any{}).
			Default(map[string]any{}),
		field.JSON("endpoints", []string{}).
			Default([]string{}),
		field.JSON("pricing", map[string]any{}).
			Default(map[string]any{}),
		field.JSON("metadata", map[string]any{}).
			Default(map[string]any{}),
		field.Enum("status").
			Values("active", "disabled").
			Default("active"),
		field.Enum("visibility").
			Values("public", "admin").
			Default("public"),
		field.String("source").
			Default("manual").
			MaxLen(50),
		field.String("icon_key").
			Default("").
			MaxLen(80),
		field.Time("last_synced_at").
			Optional().
			Nillable(),
	}
}

func (ModelCatalog) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("vendor", ModelVendor.Type).
			Field("vendor_id").
			Unique().
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}

func (ModelCatalog) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("platform", "normalized_model_id").
			Unique(),
		index.Fields("platform"),
		index.Fields("provider"),
		index.Fields("vendor_id"),
		index.Fields("status", "visibility"),
		index.Fields("source"),
	}
}

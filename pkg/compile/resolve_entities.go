package compile

import (
	"fmt"
	"strings"

	"github.com/kalo-build/morphe-go/pkg/yaml"
)

// ResolvedEntity contains the pre-resolved field path information for an entity.
// This enables downstream consumers to generate output (e.g., SQL views) without
// needing access to the full Morphe registry.
type ResolvedEntity struct {
	RootModel    string
	FieldSources map[string]ResolvedFieldSource
	Joins        []ResolvedJoin
}

// ResolvedFieldSource describes the resolved source of an entity field.
type ResolvedFieldSource struct {
	Model string
	Field string
	Type  string
}

// ResolvedJoin describes a single relationship traversal needed to reach entity fields.
type ResolvedJoin struct {
	FromModel        string
	Relationship     string
	RelationshipType string
	ToModel          string
}

// resolveEntityDefinition resolves all entity field paths against the model registry
// and returns the resolution metadata (root model, field sources, joins).
func resolveEntityDefinition(entity yaml.Entity, allModels map[string]yaml.Model) (*ResolvedEntity, error) {
	resolved := &ResolvedEntity{
		FieldSources: make(map[string]ResolvedFieldSource),
	}

	joinSet := make(map[string]ResolvedJoin)

	for fieldName, field := range entity.Fields {
		pathSegments := strings.Split(string(field.Type), ".")
		if len(pathSegments) < 2 {
			return nil, fmt.Errorf("entity %s field %s has invalid path: %s (need at least Model.Field)", entity.Name, fieldName, field.Type)
		}

		rootModelName := pathSegments[0]

		// Set root model (all fields should share the same root)
		if resolved.RootModel == "" {
			resolved.RootModel = rootModelName
		}

		rootModel, exists := allModels[rootModelName]
		if !exists {
			return nil, fmt.Errorf("entity %s field %s references unknown root model: %s", entity.Name, fieldName, rootModelName)
		}

		// Traverse relationship chain (middle segments)
		currentModel := rootModel
		for i := 1; i < len(pathSegments)-1; i++ {
			relName := pathSegments[i]

			relation, relExists := currentModel.Related[relName]
			if !relExists {
				return nil, fmt.Errorf("entity %s field %s: model %s has no relationship %s", entity.Name, fieldName, currentModel.Name, relName)
			}

			// Resolve target model name (handle aliases)
			targetModelName := relName
			if strings.TrimSpace(relation.Aliased) != "" {
				targetModelName = strings.TrimSpace(relation.Aliased)
			}

			// Record the join (deduplicated by key)
			joinKey := fmt.Sprintf("%s.%s", currentModel.Name, relName)
			if _, seen := joinSet[joinKey]; !seen {
				joinSet[joinKey] = ResolvedJoin{
					FromModel:        currentModel.Name,
					Relationship:     relName,
					RelationshipType: relation.Type,
					ToModel:          targetModelName,
				}
			}

			// Move to the next model
			nextModel, nextExists := allModels[targetModelName]
			if !nextExists {
				return nil, fmt.Errorf("entity %s field %s: relationship %s targets unknown model %s", entity.Name, fieldName, relName, targetModelName)
			}
			currentModel = nextModel
		}

		// Terminal field
		terminalFieldName := pathSegments[len(pathSegments)-1]
		terminalField, fieldExists := currentModel.Fields[terminalFieldName]
		if !fieldExists {
			return nil, fmt.Errorf("entity %s field %s: model %s has no field %s", entity.Name, fieldName, currentModel.Name, terminalFieldName)
		}

		resolved.FieldSources[fieldName] = ResolvedFieldSource{
			Model: currentModel.Name,
			Field: terminalFieldName,
			Type:  string(terminalField.Type),
		}
	}

	// Convert join set to sorted slice for deterministic output
	for _, join := range joinSet {
		resolved.Joins = append(resolved.Joins, join)
	}

	return resolved, nil
}

// serializeResolved converts a ResolvedEntity into a map suitable for YAML serialization.
func serializeResolved(resolved *ResolvedEntity) map[string]interface{} {
	fieldSources := make(map[string]interface{})
	for name, source := range resolved.FieldSources {
		fieldSources[name] = map[string]interface{}{
			"model": source.Model,
			"field": source.Field,
			"type":  source.Type,
		}
	}

	joins := make([]interface{}, 0, len(resolved.Joins))
	for _, join := range resolved.Joins {
		joins = append(joins, map[string]interface{}{
			"from_model":        join.FromModel,
			"relationship":      join.Relationship,
			"relationship_type": join.RelationshipType,
			"to_model":          join.ToModel,
		})
	}

	result := map[string]interface{}{
		"root_model":    resolved.RootModel,
		"field_sources": fieldSources,
	}

	if len(joins) > 0 {
		result["joins"] = joins
	}

	return result
}

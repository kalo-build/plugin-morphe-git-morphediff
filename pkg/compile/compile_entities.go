package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
)

// CompareEntities compares entities between base and head registries
func CompareEntities(baseReg, headReg *registry.Registry, diffDoc *diffdef.DiffDocument) error {
	baseEntities := baseReg.GetAllEntities()
	headEntities := headReg.GetAllEntities()

	// Find added entities
	for entityName, headEntity := range headEntities {
		if _, exists := baseEntities[entityName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeEntity,
				Target: map[string]string{
					"entity": entityName,
				},
				Definition:     serializeEntity(headEntity),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed entities
	for entityName := range baseEntities {
		if _, exists := headEntities[entityName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeEntity,
				Target: map[string]string{
					"entity": entityName,
				},
				Reason:         "Entity removed from schema",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified entities
	for entityName, baseEntity := range baseEntities {
		if headEntity, exists := headEntities[entityName]; exists {
			if err := compareEntityFields(entityName, baseEntity, headEntity, diffDoc); err != nil {
				return err
			}
			if err := compareEntityRelationships(entityName, baseEntity, headEntity, diffDoc); err != nil {
				return err
			}
		}
	}

	return nil
}

// compareEntityFields compares fields between two entity versions
func compareEntityFields(entityName string, baseEntity, headEntity yaml.Entity, diffDoc *diffdef.DiffDocument) error {
	// Find added fields
	for fieldName, headField := range headEntity.Fields {
		if _, exists := baseEntity.Fields[fieldName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeField,
				Target: map[string]string{
					"entity": entityName,
					"field":  fieldName,
				},
				Definition:     serializeEntityField(headField),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed fields
	for fieldName := range baseEntity.Fields {
		if _, exists := headEntity.Fields[fieldName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeField,
				Target: map[string]string{
					"entity": entityName,
					"field":  fieldName,
				},
				Reason:         "Field removed from entity",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified fields
	for fieldName, baseField := range baseEntity.Fields {
		if headField, exists := headEntity.Fields[fieldName]; exists {
			if !entityFieldsEqual(baseField, headField) {
				changes := make(map[string]interface{})

				// Check type change (field path)
				if baseField.Type != headField.Type {
					changes["type"] = map[string]interface{}{
						"before": string(baseField.Type),
						"after":  string(headField.Type),
					}
				}

				// Check attribute changes
				if !attributesEqual(baseField.Attributes, headField.Attributes) {
					changes["attributes"] = map[string]interface{}{
						"before": baseField.Attributes,
						"after":  headField.Attributes,
					}
				}

				if len(changes) > 0 {
					change := diffdef.Change{
						Operation: diffdef.OperationModify,
						Type:      diffdef.TypeField,
						Target: map[string]string{
							"entity": entityName,
							"field":  fieldName,
						},
						Changes:        changes,
						Classification: classifyEntityFieldModification(baseField, headField),
					}
					diffDoc.AddChange(change)
				}
			}
		}
	}

	return nil
}

// compareEntityRelationships compares relationships between two entity versions
func compareEntityRelationships(entityName string, baseEntity, headEntity yaml.Entity, diffDoc *diffdef.DiffDocument) error {
	// Find added relationships
	for relName, headRel := range headEntity.Related {
		if _, exists := baseEntity.Related[relName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeRelationship,
				Target: map[string]string{
					"entity":       entityName,
					"relationship": relName,
				},
				Definition:     serializeEntityRelationship(headRel),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed relationships
	for relName := range baseEntity.Related {
		if _, exists := headEntity.Related[relName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeRelationship,
				Target: map[string]string{
					"entity":       entityName,
					"relationship": relName,
				},
				Reason:         "Relationship removed from entity",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified relationships
	for relName, baseRel := range baseEntity.Related {
		if headRel, exists := headEntity.Related[relName]; exists {
			if !entityRelationshipsEqual(baseRel, headRel) {
				changes := make(map[string]interface{})

				if baseRel.Type != headRel.Type {
					changes["type"] = map[string]interface{}{
						"before": baseRel.Type,
						"after":  headRel.Type,
					}
				}

				if baseRel.Through != headRel.Through {
					changes["through"] = map[string]interface{}{
						"before": baseRel.Through,
						"after":  headRel.Through,
					}
				}

				if baseRel.Aliased != headRel.Aliased {
					changes["aliased"] = map[string]interface{}{
						"before": baseRel.Aliased,
						"after":  headRel.Aliased,
					}
				}

				if !stringSliceEqual(baseRel.For, headRel.For) {
					changes["for"] = map[string]interface{}{
						"before": baseRel.For,
						"after":  headRel.For,
					}
				}

				if len(changes) > 0 {
					change := diffdef.Change{
						Operation: diffdef.OperationModify,
						Type:      diffdef.TypeRelationship,
						Target: map[string]string{
							"entity":       entityName,
							"relationship": relName,
						},
						Changes:        changes,
						Classification: diffdef.ClassificationBreaking,
					}
					diffDoc.AddChange(change)
				}
			}
		}
	}

	return nil
}

// Helper functions

func serializeEntity(entity yaml.Entity) map[string]interface{} {
	fields := make(map[string]interface{})
	for name, field := range entity.Fields {
		fields[name] = serializeEntityField(field)
	}

	identifiers := make(map[string]interface{})
	for name, id := range entity.Identifiers {
		identifiers[name] = id.Fields
	}

	related := make(map[string]interface{})
	for name, rel := range entity.Related {
		related[name] = serializeEntityRelationship(rel)
	}

	return map[string]interface{}{
		"fields":      fields,
		"identifiers": identifiers,
		"related":     related,
	}
}

func serializeEntityField(field yaml.EntityField) map[string]interface{} {
	result := map[string]interface{}{
		"type": string(field.Type),
	}
	if len(field.Attributes) > 0 {
		result["attributes"] = field.Attributes
	}
	return result
}

func serializeEntityRelationship(rel yaml.EntityRelation) map[string]interface{} {
	result := map[string]interface{}{
		"type": rel.Type,
	}
	if rel.Through != "" {
		result["through"] = rel.Through
	}
	if rel.Aliased != "" {
		result["aliased"] = rel.Aliased
	}
	if len(rel.For) > 0 {
		result["for"] = rel.For
	}
	return result
}

func entityFieldsEqual(a, b yaml.EntityField) bool {
	return a.Type == b.Type && attributesEqual(a.Attributes, b.Attributes)
}

func entityRelationshipsEqual(a, b yaml.EntityRelation) bool {
	return a.Type == b.Type &&
		a.Through == b.Through &&
		a.Aliased == b.Aliased &&
		stringSliceEqual(a.For, b.For)
}

func classifyEntityFieldModification(baseField, headField yaml.EntityField) string {
	// Type (field path) changes are breaking
	if baseField.Type != headField.Type {
		return diffdef.ClassificationBreaking
	}

	// Attribute changes are generally breaking
	if !attributesEqual(baseField.Attributes, headField.Attributes) {
		return diffdef.ClassificationBreaking
	}

	return diffdef.ClassificationSafe
}

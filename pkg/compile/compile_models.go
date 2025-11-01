package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
)

// CompareModels compares models between base and head registries
func CompareModels(baseReg, headReg *registry.Registry, diffDoc *diffdef.DiffDocument) error {
	baseModels := baseReg.GetAllModels()
	headModels := headReg.GetAllModels()

	// Find added models
	for modelName, headModel := range headModels {
		if _, exists := baseModels[modelName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeModel,
				Target: map[string]string{
					"model": modelName,
				},
				Definition:     serializeModel(headModel),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed models
	for modelName := range baseModels {
		if _, exists := headModels[modelName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeModel,
				Target: map[string]string{
					"model": modelName,
				},
				Reason:         "Model removed from schema",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified models
	for modelName, baseModel := range baseModels {
		if headModel, exists := headModels[modelName]; exists {
			if err := compareModelFields(modelName, baseModel, headModel, diffDoc); err != nil {
				return err
			}
			if err := compareModelRelationships(modelName, baseModel, headModel, diffDoc); err != nil {
				return err
			}
		}
	}

	return nil
}

// compareModelFields compares fields between two model versions
func compareModelFields(modelName string, baseModel, headModel yaml.Model, diffDoc *diffdef.DiffDocument) error {
	// Find added fields
	for fieldName, headField := range headModel.Fields {
		if _, exists := baseModel.Fields[fieldName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeField,
				Target: map[string]string{
					"model": modelName,
					"field": fieldName,
				},
				Definition:     serializeField(headField),
				Classification: classifyFieldAddition(headField),
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed fields
	for fieldName := range baseModel.Fields {
		if _, exists := headModel.Fields[fieldName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeField,
				Target: map[string]string{
					"model": modelName,
					"field": fieldName,
				},
				Reason:         "Field removed from model",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified fields
	for fieldName, baseField := range baseModel.Fields {
		if headField, exists := headModel.Fields[fieldName]; exists {
			if !fieldsEqual(baseField, headField) {
				changes := make(map[string]interface{})

				// Check type change
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
							"model": modelName,
							"field": fieldName,
						},
						Changes:        changes,
						Classification: classifyFieldModification(baseField, headField),
					}
					diffDoc.AddChange(change)
				}
			}
		}
	}

	return nil
}

// compareModelRelationships compares relationships between two model versions
func compareModelRelationships(modelName string, baseModel, headModel yaml.Model, diffDoc *diffdef.DiffDocument) error {
	// Find added relationships
	for relName, headRel := range headModel.Related {
		if _, exists := baseModel.Related[relName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeRelationship,
				Target: map[string]string{
					"model":        modelName,
					"relationship": relName,
				},
				Definition:     serializeRelationship(headRel),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed relationships
	for relName := range baseModel.Related {
		if _, exists := headModel.Related[relName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeRelationship,
				Target: map[string]string{
					"model":        modelName,
					"relationship": relName,
				},
				Reason:         "Relationship removed from model",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified relationships
	for relName, baseRel := range baseModel.Related {
		if headRel, exists := headModel.Related[relName]; exists {
			if !relationshipsEqual(baseRel, headRel) {
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
							"model":        modelName,
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

func serializeModel(model yaml.Model) map[string]interface{} {
	fields := make(map[string]interface{})
	for name, field := range model.Fields {
		fields[name] = serializeField(field)
	}

	identifiers := make(map[string]interface{})
	for name, id := range model.Identifiers {
		identifiers[name] = id.Fields
	}

	related := make(map[string]interface{})
	for name, rel := range model.Related {
		related[name] = serializeRelationship(rel)
	}

	return map[string]interface{}{
		"fields":      fields,
		"identifiers": identifiers,
		"related":     related,
	}
}

func serializeField(field yaml.ModelField) map[string]interface{} {
	result := map[string]interface{}{
		"type": string(field.Type),
	}
	if len(field.Attributes) > 0 {
		result["attributes"] = field.Attributes
	}
	return result
}

func serializeRelationship(rel yaml.ModelRelation) map[string]interface{} {
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

func fieldsEqual(a, b yaml.ModelField) bool {
	return a.Type == b.Type && attributesEqual(a.Attributes, b.Attributes)
}

func relationshipsEqual(a, b yaml.ModelRelation) bool {
	return a.Type == b.Type &&
		a.Through == b.Through &&
		a.Aliased == b.Aliased &&
		stringSliceEqual(a.For, b.For)
}

func classifyFieldAddition(field yaml.ModelField) string {
	// If field has attributes indicating it's optional/nullable, it's additive
	// Otherwise it could be breaking (e.g., mandatory field)
	for _, attr := range field.Attributes {
		if attr == "nullable" || attr == "optional" {
			return diffdef.ClassificationAdditive
		}
		if attr == "mandatory" || attr == "required" {
			return diffdef.ClassificationBreaking
		}
	}
	// Default to additive for new fields
	return diffdef.ClassificationAdditive
}

func classifyFieldModification(baseField, headField yaml.ModelField) string {
	// Type changes are breaking
	if baseField.Type != headField.Type {
		return diffdef.ClassificationBreaking
	}

	// Check if attributes became more restrictive
	baseMandatory := containsAttribute(baseField.Attributes, "mandatory")
	headMandatory := containsAttribute(headField.Attributes, "mandatory")

	if !baseMandatory && headMandatory {
		return diffdef.ClassificationBreaking
	}

	// Other attribute changes are generally breaking
	if !attributesEqual(baseField.Attributes, headField.Attributes) {
		return diffdef.ClassificationBreaking
	}

	return diffdef.ClassificationSafe
}

func containsAttribute(attributes []string, attr string) bool {
	for _, a := range attributes {
		if a == attr {
			return true
		}
	}
	return false
}

package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
)

// CompareStructures compares structures between base and head registries
func CompareStructures(baseReg, headReg *registry.Registry, diffDoc *diffdef.DiffDocument) error {
	baseStructures := baseReg.GetAllStructures()
	headStructures := headReg.GetAllStructures()

	// Find added structures
	for structName, headStruct := range headStructures {
		if _, exists := baseStructures[structName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeStructure,
				Target: map[string]string{
					"structure": structName,
				},
				Definition:     serializeStructure(headStruct),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed structures
	for structName := range baseStructures {
		if _, exists := headStructures[structName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeStructure,
				Target: map[string]string{
					"structure": structName,
				},
				Reason:         "Structure removed from schema",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified structures
	for structName, baseStruct := range baseStructures {
		if headStruct, exists := headStructures[structName]; exists {
			if err := compareStructureFields(structName, baseStruct, headStruct, diffDoc); err != nil {
				return err
			}
		}
	}

	return nil
}

// compareStructureFields compares fields between two structure versions
func compareStructureFields(structName string, baseStruct, headStruct yaml.Structure, diffDoc *diffdef.DiffDocument) error {
	// Find added fields
	for fieldName, headField := range headStruct.Fields {
		if _, exists := baseStruct.Fields[fieldName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeField,
				Target: map[string]string{
					"structure": structName,
					"field":     fieldName,
				},
				Definition:     serializeStructureField(headField),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed fields
	for fieldName := range baseStruct.Fields {
		if _, exists := headStruct.Fields[fieldName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeField,
				Target: map[string]string{
					"structure": structName,
					"field":     fieldName,
				},
				Reason:         "Field removed from structure",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified fields
	for fieldName, baseField := range baseStruct.Fields {
		if headField, exists := headStruct.Fields[fieldName]; exists {
			if !structureFieldsEqual(baseField, headField) {
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
							"structure": structName,
							"field":     fieldName,
						},
						Changes:        changes,
						Classification: classifyStructureFieldModification(baseField, headField),
					}
					diffDoc.AddChange(change)
				}
			}
		}
	}

	return nil
}

// Helper functions

func serializeStructure(structure yaml.Structure) map[string]interface{} {
	fields := make(map[string]interface{})
	for name, field := range structure.Fields {
		fields[name] = serializeStructureField(field)
	}

	return map[string]interface{}{
		"fields": fields,
	}
}

func serializeStructureField(field yaml.StructureField) map[string]interface{} {
	result := map[string]interface{}{
		"type": string(field.Type),
	}
	if len(field.Attributes) > 0 {
		result["attributes"] = field.Attributes
	}
	return result
}

func structureFieldsEqual(a, b yaml.StructureField) bool {
	return a.Type == b.Type && attributesEqual(a.Attributes, b.Attributes)
}

func classifyStructureFieldModification(baseField, headField yaml.StructureField) string {
	// Type changes are breaking
	if baseField.Type != headField.Type {
		return diffdef.ClassificationBreaking
	}

	// Attribute changes are generally breaking for structures
	if !attributesEqual(baseField.Attributes, headField.Attributes) {
		return diffdef.ClassificationBreaking
	}

	return diffdef.ClassificationSafe
}

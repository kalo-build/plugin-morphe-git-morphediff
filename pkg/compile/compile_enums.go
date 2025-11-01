package compile

import (
	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
)

// CompareEnums compares enums between base and head registries
func CompareEnums(baseReg, headReg *registry.Registry, diffDoc *diffdef.DiffDocument) error {
	baseEnums := baseReg.GetAllEnums()
	headEnums := headReg.GetAllEnums()

	// Find added enums
	for enumName, headEnum := range headEnums {
		if _, exists := baseEnums[enumName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationAdd,
				Type:      diffdef.TypeEnum,
				Target: map[string]string{
					"enum": enumName,
				},
				Definition:     serializeEnum(headEnum),
				Classification: diffdef.ClassificationAdditive,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find removed enums
	for enumName := range baseEnums {
		if _, exists := headEnums[enumName]; !exists {
			change := diffdef.Change{
				Operation: diffdef.OperationRemove,
				Type:      diffdef.TypeEnum,
				Target: map[string]string{
					"enum": enumName,
				},
				Reason:         "Enum removed from schema",
				Classification: diffdef.ClassificationBreaking,
			}
			diffDoc.AddChange(change)
		}
	}

	// Find modified enums
	for enumName, baseEnum := range baseEnums {
		if headEnum, exists := headEnums[enumName]; exists {
			if err := compareEnumEntries(enumName, baseEnum, headEnum, diffDoc); err != nil {
				return err
			}
		}
	}

	return nil
}

// compareEnumEntries compares entries between two enum versions
func compareEnumEntries(enumName string, baseEnum, headEnum yaml.Enum, diffDoc *diffdef.DiffDocument) error {
	changes := make(map[string]interface{})
	hasChanges := false

	// Find added entries
	addedEntries := make(map[string]interface{})
	for entryName, entryValue := range headEnum.Entries {
		if _, exists := baseEnum.Entries[entryName]; !exists {
			addedEntries[entryName] = entryValue
		}
	}
	if len(addedEntries) > 0 {
		if !hasEntriesChanges(changes) {
			changes["entries"] = make(map[string]interface{})
		}
		changes["entries"].(map[string]interface{})["added"] = addedEntries
		hasChanges = true
	}

	// Find removed entries
	removedEntries := make(map[string]interface{})
	for entryName, entryValue := range baseEnum.Entries {
		if _, exists := headEnum.Entries[entryName]; !exists {
			removedEntries[entryName] = entryValue
		}
	}
	if len(removedEntries) > 0 {
		if !hasEntriesChanges(changes) {
			changes["entries"] = make(map[string]interface{})
		}
		changes["entries"].(map[string]interface{})["removed"] = removedEntries
		hasChanges = true
	}

	// Find modified entries
	modifiedEntries := make([]map[string]interface{}, 0)
	for entryName, baseValue := range baseEnum.Entries {
		if headValue, exists := headEnum.Entries[entryName]; exists {
			if baseValue != headValue {
				modifiedEntries = append(modifiedEntries, map[string]interface{}{
					"symbol": entryName,
					"before": baseValue,
					"after":  headValue,
				})
			}
		}
	}
	if len(modifiedEntries) > 0 {
		if !hasEntriesChanges(changes) {
			changes["entries"] = make(map[string]interface{})
		}
		changes["entries"].(map[string]interface{})["modified"] = modifiedEntries
		hasChanges = true
	}

	// Check type change
	if baseEnum.Type != headEnum.Type {
		changes["type"] = map[string]interface{}{
			"before": string(baseEnum.Type),
			"after":  string(headEnum.Type),
		}
		hasChanges = true
	}

	if hasChanges {
		change := diffdef.Change{
			Operation: diffdef.OperationModify,
			Type:      diffdef.TypeEnum,
			Target: map[string]string{
				"enum": enumName,
			},
			Changes:        changes,
			Classification: classifyEnumModification(len(addedEntries), len(removedEntries), len(modifiedEntries)),
		}
		diffDoc.AddChange(change)
	}

	return nil
}

// Helper functions

func serializeEnum(enum yaml.Enum) map[string]interface{} {
	return map[string]interface{}{
		"type":    string(enum.Type),
		"entries": enum.Entries,
	}
}

func hasEntriesChanges(changes map[string]interface{}) bool {
	_, exists := changes["entries"]
	return exists
}

func classifyEnumModification(added, removed, modified int) string {
	// Removing entries or modifying them is breaking
	if removed > 0 || modified > 0 {
		return diffdef.ClassificationBreaking
	}
	// Only adding entries is additive
	if added > 0 {
		return diffdef.ClassificationAdditive
	}
	return diffdef.ClassificationSafe
}

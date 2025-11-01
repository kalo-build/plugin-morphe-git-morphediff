package compile_test

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"github.com/stretchr/testify/suite"
)

type CompileEnumsTestSuite struct {
	suite.Suite
}

func TestCompileEnumsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEnumsTestSuite))
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_AddedEnum() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	headEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":  "ADMIN",
			"Editor": "EDITOR",
			"Viewer": "VIEWER",
		},
	}
	headReg.SetEnum("UserRole", headEnum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeEnum, change.Type)
	suite.Equal("UserRole", change.Target["enum"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_RemovedEnum() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEnum := yaml.Enum{
		Name: "LegacyStatus",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Active": "ACTIVE",
		},
	}
	baseReg.SetEnum("LegacyStatus", baseEnum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeEnum, change.Type)
	suite.Equal("LegacyStatus", change.Target["enum"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_AddedEntries() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin": "ADMIN",
		},
	}
	baseReg.SetEnum("UserRole", baseEnum)

	headEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":     "ADMIN",
			"Moderator": "MOD",
		},
	}
	headReg.SetEnum("UserRole", headEnum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeEnum, change.Type)
	suite.Equal("UserRole", change.Target["enum"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
	suite.NotNil(change.Changes["entries"])
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_RemovedEntries() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin":       "ADMIN",
			"LegacyAdmin": "LEGACY_ADMIN",
		},
	}
	baseReg.SetEnum("UserRole", baseEnum)

	headEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin": "ADMIN",
		},
	}
	headReg.SetEnum("UserRole", headEnum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeEnum, change.Type)
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_ModifiedEntries() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin": "ADMIN",
		},
	}
	baseReg.SetEnum("UserRole", baseEnum)

	headEnum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin": "ADMINISTRATOR",
		},
	}
	headReg.SetEnum("UserRole", headEnum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_TypeChange() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEnum := yaml.Enum{
		Name: "Priority",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"High": "HIGH",
		},
	}
	baseReg.SetEnum("Priority", baseEnum)

	headEnum := yaml.Enum{
		Name: "Priority",
		Type: yaml.EnumTypeInteger,
		Entries: map[string]any{
			"High": 1,
		},
	}
	headReg.SetEnum("Priority", headEnum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.NotNil(change.Changes["type"])
}

func (suite *CompileEnumsTestSuite) TestCompareEnums_NoChanges() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	enum := yaml.Enum{
		Name: "UserRole",
		Type: yaml.EnumTypeString,
		Entries: map[string]any{
			"Admin": "ADMIN",
		},
	}
	baseReg.SetEnum("UserRole", enum)
	headReg.SetEnum("UserRole", enum)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEnums(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 0)
}

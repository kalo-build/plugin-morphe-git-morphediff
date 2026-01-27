package compile_test

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"github.com/stretchr/testify/suite"
)

type CompileStructuresTestSuite struct {
	suite.Suite
}

func TestCompileStructuresTestSuite(t *testing.T) {
	suite.Run(t, new(CompileStructuresTestSuite))
}

func (suite *CompileStructuresTestSuite) TestCompareStructures_AddedStructure() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	headStruct := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {
				Type: yaml.StructureFieldTypeString,
			},
			"City": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}
	headReg.SetStructure("Address", headStruct)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareStructures(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeStructure, change.Type)
	suite.Equal("Address", change.Target["structure"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileStructuresTestSuite) TestCompareStructures_RemovedStructure() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseStruct := yaml.Structure{
		Name: "LegacyData",
		Fields: map[string]yaml.StructureField{
			"Value": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}
	baseReg.SetStructure("LegacyData", baseStruct)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareStructures(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeStructure, change.Type)
	suite.Equal("LegacyData", change.Target["structure"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileStructuresTestSuite) TestCompareStructures_AddedField() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseStruct := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}
	baseReg.SetStructure("Address", baseStruct)

	headStruct := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {
				Type: yaml.StructureFieldTypeString,
			},
			"Country": {
				Type: yaml.StructureFieldTypeString,
			},
		},
	}
	headReg.SetStructure("Address", headStruct)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareStructures(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("Address", change.Target["structure"])
	suite.Equal("Country", change.Target["field"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileStructuresTestSuite) TestCompareStructures_RemovedField() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseStruct := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street":   {Type: yaml.StructureFieldTypeString},
			"LegacyID": {Type: yaml.StructureFieldTypeString},
		},
	}
	baseReg.SetStructure("Address", baseStruct)

	headStruct := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {Type: yaml.StructureFieldTypeString},
		},
	}
	headReg.SetStructure("Address", headStruct)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareStructures(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("Address", change.Target["structure"])
	suite.Equal("LegacyID", change.Target["field"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileStructuresTestSuite) TestCompareStructures_ModifiedFieldType() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseStruct := yaml.Structure{
		Name: "Metadata",
		Fields: map[string]yaml.StructureField{
			"Count": {Type: yaml.StructureFieldTypeString},
		},
	}
	baseReg.SetStructure("Metadata", baseStruct)

	headStruct := yaml.Structure{
		Name: "Metadata",
		Fields: map[string]yaml.StructureField{
			"Count": {Type: yaml.StructureFieldTypeInteger},
		},
	}
	headReg.SetStructure("Metadata", headStruct)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareStructures(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("Metadata", change.Target["structure"])
	suite.Equal("Count", change.Target["field"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
	suite.NotNil(change.Changes["type"])
}

func (suite *CompileStructuresTestSuite) TestCompareStructures_NoChanges() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	structure := yaml.Structure{
		Name: "Address",
		Fields: map[string]yaml.StructureField{
			"Street": {Type: yaml.StructureFieldTypeString},
		},
	}
	baseReg.SetStructure("Address", structure)
	headReg.SetStructure("Address", structure)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareStructures(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 0)
}



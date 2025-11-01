package compile_test

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"github.com/stretchr/testify/suite"
)

type CompileModelsTestSuite struct {
	suite.Suite
}

func TestCompileModelsTestSuite(t *testing.T) {
	suite.Run(t, new(CompileModelsTestSuite))
}

func (suite *CompileModelsTestSuite) TestCompareModels_AddedModel() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	// Add a model only in head
	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
				Attributes: []string{
					"mandatory",
				},
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeModel, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
	suite.NotNil(change.Definition)
}

func (suite *CompileModelsTestSuite) TestCompareModels_RemovedModel() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	// Add a model only in base
	baseModel := yaml.Model{
		Name: "LegacyUser",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("LegacyUser", baseModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeModel, change.Type)
	suite.Equal("LegacyUser", change.Target["model"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
	suite.NotEmpty(change.Reason)
}

func (suite *CompileModelsTestSuite) TestCompareModels_AddedField() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"PhoneNumber": {
				Type: yaml.ModelFieldTypeString,
				Attributes: []string{
					"nullable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal("PhoneNumber", change.Target["field"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileModelsTestSuite) TestCompareModels_RemovedField() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"LegacyID": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal("LegacyID", change.Target["field"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileModelsTestSuite) TestCompareModels_ModifiedField_TypeChange() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Age": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Age": {
				Type: yaml.ModelFieldTypeInteger,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal("Age", change.Target["field"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
	suite.NotNil(change.Changes["type"])
}

func (suite *CompileModelsTestSuite) TestCompareModels_ModifiedField_AttributeChange() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Email": {
				Type: yaml.ModelFieldTypeString,
				Attributes: []string{
					"nullable",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Email": {
				Type: yaml.ModelFieldTypeString,
				Attributes: []string{
					"mandatory",
				},
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
	suite.NotNil(change.Changes["attributes"])
}

func (suite *CompileModelsTestSuite) TestCompareModels_AddedRelationship() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"ProfileImage": {
				Type: "HasOne",
			},
		},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeRelationship, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal("ProfileImage", change.Target["relationship"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileModelsTestSuite) TestCompareModels_RemovedRelationship() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"LegacyProfile": {
				Type: "HasOne",
			},
		},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeRelationship, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal("LegacyProfile", change.Target["relationship"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileModelsTestSuite) TestCompareModels_ModifiedRelationship_CardinalityChange() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Address": {
				Type: "HasOne",
			},
		},
	}
	baseReg.SetModel("User", baseModel)

	headModel := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Address": {
				Type: "HasMany",
			},
		},
	}
	headReg.SetModel("User", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeRelationship, change.Type)
	suite.Equal("User", change.Target["model"])
	suite.Equal("Address", change.Target["relationship"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
	suite.NotNil(change.Changes["type"])
}

func (suite *CompileModelsTestSuite) TestCompareModels_NoChanges() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	model := yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeUUID,
			},
			"Name": {
				Type: yaml.ModelFieldTypeString,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{},
	}
	baseReg.SetModel("User", model)
	headReg.SetModel("User", model)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 0)
}

func (suite *CompileModelsTestSuite) TestCompareModels_PolymorphicRelationship() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post"},
			},
		},
	}
	baseReg.SetModel("Comment", baseModel)

	headModel := yaml.Model{
		Name: "Comment",
		Fields: map[string]yaml.ModelField{
			"ID": {
				Type: yaml.ModelFieldTypeAutoIncrement,
			},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.ModelRelation{
			"Commentable": {
				Type: "ForOnePoly",
				For:  []string{"Post", "Article"},
			},
		},
	}
	headReg.SetModel("Comment", headModel)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareModels(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeRelationship, change.Type)
	suite.NotNil(change.Changes["for"])
}

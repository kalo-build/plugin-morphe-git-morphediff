package compile_test

import (
	"testing"

	"github.com/kalo-build/morphe-go/pkg/registry"
	"github.com/kalo-build/morphe-go/pkg/yaml"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"github.com/stretchr/testify/suite"
)

type CompileEntitiesTestSuite struct {
	suite.Suite
}

func TestCompileEntitiesTestSuite(t *testing.T) {
	suite.Run(t, new(CompileEntitiesTestSuite))
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_AddedEntity() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
				Attributes: []string{
					"immutable",
					"mandatory",
				},
			},
			"Email": {
				Type: "User.Email",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeEntity, change.Type)
	suite.Equal("UserProfile", change.Target["entity"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_RemovedEntity() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEntity := yaml.Entity{
		Name: "LegacyView",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("LegacyView", baseEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeEntity, change.Type)
	suite.Equal("LegacyView", change.Target["entity"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_AddedField() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("UserProfile", baseEntity)

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
			"CompanyName": {
				Type: "User.Company.Name",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("UserProfile", change.Target["entity"])
	suite.Equal("CompanyName", change.Target["field"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_RemovedField() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
			"LegacyField": {
				Type: "User.LegacyField",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("UserProfile", baseEntity)

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationRemove, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("UserProfile", change.Target["entity"])
	suite.Equal("LegacyField", change.Target["field"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_ModifiedFieldPath() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
			"Email": {
				Type: "User.Email",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("UserProfile", baseEntity)

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
			"Email": {
				Type: "User.ContactInfo.Email",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationModify, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("UserProfile", change.Target["entity"])
	suite.Equal("Email", change.Target["field"])
	suite.Equal(diffdef.ClassificationBreaking, change.Classification)
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_AddedRelationship() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	baseEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("UserProfile", baseEntity)

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{
			"Company": {
				Type: "ForOne",
			},
		},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeRelationship, change.Type)
	suite.Equal("UserProfile", change.Target["entity"])
	suite.Equal("Company", change.Target["relationship"])
	suite.Equal(diffdef.ClassificationAdditive, change.Classification)
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_AddedEntity_WithResolution() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	// Set up models in the head registry for resolution
	headReg.SetModel("User", yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: "AutoIncrement"},
			"Email": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{
			"ContactInfo": {Type: "HasOne"},
		},
	})
	headReg.SetModel("ContactInfo", yaml.Model{
		Name: "ContactInfo",
		Fields: map[string]yaml.ModelField{
			"ID":    {Type: "AutoIncrement"},
			"Phone": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	})

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
			"Phone": {
				Type: "User.ContactInfo.Phone",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeEntity, change.Type)

	// Verify the resolved block is present in the definition
	resolved, hasResolved := change.Definition["resolved"].(map[string]interface{})
	suite.True(hasResolved, "definition should contain 'resolved' block")
	suite.Equal("User", resolved["root_model"])

	fieldSources, hasFS := resolved["field_sources"].(map[string]interface{})
	suite.True(hasFS, "resolved should contain 'field_sources'")

	idSource, hasID := fieldSources["ID"].(map[string]interface{})
	suite.True(hasID)
	suite.Equal("User", idSource["model"])
	suite.Equal("ID", idSource["field"])
	suite.Equal("AutoIncrement", idSource["type"])

	phoneSource, hasPhone := fieldSources["Phone"].(map[string]interface{})
	suite.True(hasPhone)
	suite.Equal("ContactInfo", phoneSource["model"])
	suite.Equal("Phone", phoneSource["field"])
	suite.Equal("String", phoneSource["type"])

	joins, hasJoins := resolved["joins"].([]interface{})
	suite.True(hasJoins, "resolved should contain 'joins'")
	suite.Len(joins, 1)

	join := joins[0].(map[string]interface{})
	suite.Equal("User", join["from_model"])
	suite.Equal("ContactInfo", join["relationship"])
	suite.Equal("HasOne", join["relationship_type"])
	suite.Equal("ContactInfo", join["to_model"])
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_AddedField_WithEntitySnapshot() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	// Set up models
	headReg.SetModel("User", yaml.Model{
		Name: "User",
		Fields: map[string]yaml.ModelField{
			"ID":   {Type: "AutoIncrement"},
			"Name": {Type: "String"},
		},
		Identifiers: map[string]yaml.ModelIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.ModelRelation{},
	})

	baseEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {Type: "User.ID"},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("UserProfile", baseEntity)

	headEntity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID":       {Type: "User.ID"},
			"FullName": {Type: "User.Name"},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {Fields: []string{"ID"}},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	headReg.SetEntity("UserProfile", headEntity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 1)

	change := diffDoc.Changes[0]
	suite.Equal(diffdef.OperationAdd, change.Operation)
	suite.Equal(diffdef.TypeField, change.Type)
	suite.Equal("FullName", change.Target["field"])

	// Verify entity_snapshot is present with full post-change entity
	suite.NotNil(change.EntitySnapshot, "change should have entity_snapshot")

	resolved, hasResolved := change.EntitySnapshot["resolved"].(map[string]interface{})
	suite.True(hasResolved, "entity_snapshot should contain 'resolved' block")
	suite.Equal("User", resolved["root_model"])

	fieldSources, hasFS := resolved["field_sources"].(map[string]interface{})
	suite.True(hasFS)
	suite.Len(fieldSources, 2) // ID and FullName

	nameSource := fieldSources["FullName"].(map[string]interface{})
	suite.Equal("User", nameSource["model"])
	suite.Equal("Name", nameSource["field"])
	suite.Equal("String", nameSource["type"])
}

func (suite *CompileEntitiesTestSuite) TestCompareEntities_NoChanges() {
	baseReg := registry.NewRegistry()
	headReg := registry.NewRegistry()

	entity := yaml.Entity{
		Name: "UserProfile",
		Fields: map[string]yaml.EntityField{
			"ID": {
				Type: "User.ID",
			},
		},
		Identifiers: map[string]yaml.EntityIdentifier{
			"primary": {
				Fields: []string{"ID"},
			},
		},
		Related: map[string]yaml.EntityRelation{},
	}
	baseReg.SetEntity("UserProfile", entity)
	headReg.SetEntity("UserProfile", entity)

	diffDoc := diffdef.NewDiffDocument("base", "head")
	err := compile.CompareEntities(baseReg, headReg, diffDoc)

	suite.Nil(err)
	suite.Len(diffDoc.Changes, 0)
}

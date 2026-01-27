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



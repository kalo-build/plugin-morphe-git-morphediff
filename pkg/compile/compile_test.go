package compile_test

import (
	"os"
	"path/filepath"
	"testing"

	rcfg "github.com/kalo-build/morphe-go/pkg/registry/cfg"
	"github.com/kalo-build/plugin-morphe-git-morphediff/internal/testutils"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/compile"
	"github.com/kalo-build/plugin-morphe-git-morphediff/pkg/diffdef"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type CompileTestSuite struct {
	suite.Suite
	TestDirPath string
}

func TestCompileTestSuite(t *testing.T) {
	suite.Run(t, new(CompileTestSuite))
}

func (suite *CompileTestSuite) SetupTest() {
	suite.TestDirPath = testutils.GetTestDirPath()
}

func (suite *CompileTestSuite) TestMorpheToMorpheDiff() {
	basePath := filepath.Join(suite.TestDirPath, "base")
	headPath := filepath.Join(suite.TestDirPath, "head")
	outputPath := filepath.Join(suite.TestDirPath, "working", "output.yaml")
	expectedPath := filepath.Join(suite.TestDirPath, "expected", "morphe-diff.yaml")

	// Ensure working directory exists and is clean
	workingDir := filepath.Dir(outputPath)
	os.RemoveAll(workingDir)
	err := os.MkdirAll(workingDir, 0755)
	suite.Nil(err)
	defer os.RemoveAll(workingDir)

	config := compile.MorpheDiffConfig{
		BaseRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(basePath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(basePath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(basePath, "enums"),
			RegistryStructuresDirPath: filepath.Join(basePath, "structures"),
		},
		HeadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath:     filepath.Join(headPath, "models"),
			RegistryEntitiesDirPath:   filepath.Join(headPath, "entities"),
			RegistryEnumsDirPath:      filepath.Join(headPath, "enums"),
			RegistryStructuresDirPath: filepath.Join(headPath, "structures"),
		},
		OutputPath:  outputPath,
		BaseVersion: "base",
		HeadVersion: "head",
	}

	err = compile.MorpheToMorpheDiff(config)
	suite.Nil(err)

	// Verify output file exists
	suite.FileExists(outputPath)

	// Read generated output
	generatedData, err := os.ReadFile(outputPath)
	suite.Nil(err)

	var generatedDiff diffdef.DiffDocument
	err = yaml.Unmarshal(generatedData, &generatedDiff)
	suite.Nil(err)

	// Read expected output
	expectedData, err := os.ReadFile(expectedPath)
	suite.Nil(err)

	var expectedDiff diffdef.DiffDocument
	err = yaml.Unmarshal(expectedData, &expectedDiff)
	suite.Nil(err)

	// Verify metadata
	suite.Equal(expectedDiff.Metadata.SpecVersion, generatedDiff.Metadata.SpecVersion)
	// Compare Ref (new field) - the test config sets BaseVersion/HeadVersion which maps to Ref
	suite.Equal("base", generatedDiff.Metadata.Source.Ref)
	suite.Equal("head", generatedDiff.Metadata.Target.Ref)

	// Verify summary counts
	suite.Equal(8, generatedDiff.Metadata.Summary.TotalChanges)
	suite.Equal(1, generatedDiff.Metadata.Summary.Breaking)
	suite.Equal(7, generatedDiff.Metadata.Summary.Additive)

	// Verify changes count
	suite.Len(generatedDiff.Changes, 8)

	// Verify key changes exist
	foundPhoneNumberAdd := false
	foundEmailModify := false
	foundOrgAdd := false
	foundEntityFieldAdd := false
	foundEntityRelAdd := false

	for _, change := range generatedDiff.Changes {
		if change.Operation == diffdef.OperationAdd &&
			change.Type == diffdef.TypeField &&
			change.Target["model"] == "User" &&
			change.Target["field"] == "PhoneNumber" {
			foundPhoneNumberAdd = true
			suite.Equal(diffdef.ClassificationAdditive, change.Classification)
		}

		if change.Operation == diffdef.OperationModify &&
			change.Type == diffdef.TypeField &&
			change.Target["model"] == "User" &&
			change.Target["field"] == "Email" {
			foundEmailModify = true
			suite.Equal(diffdef.ClassificationBreaking, change.Classification)
		}

		if change.Operation == diffdef.OperationAdd &&
			change.Type == diffdef.TypeModel &&
			change.Target["model"] == "Organization" {
			foundOrgAdd = true
			suite.Equal(diffdef.ClassificationAdditive, change.Classification)
		}

		// Entity field addition: FullName added to UserProfile
		if change.Operation == diffdef.OperationAdd &&
			change.Type == diffdef.TypeField &&
			change.Target["entity"] == "UserProfile" &&
			change.Target["field"] == "FullName" {
			foundEntityFieldAdd = true
			suite.Equal(diffdef.ClassificationAdditive, change.Classification)

			// Verify entity_snapshot is present
			suite.NotNil(change.EntitySnapshot, "Entity field change should include entity_snapshot")

			// Verify resolved block inside entity_snapshot
			resolved, hasResolved := change.EntitySnapshot["resolved"].(map[string]interface{})
			suite.True(hasResolved, "entity_snapshot should contain 'resolved' block")
			suite.Equal("User", resolved["root_model"], "resolved root_model should be User")

			// Verify field_sources contains all post-change fields
			fieldSources, hasFS := resolved["field_sources"].(map[string]interface{})
			suite.True(hasFS, "resolved should contain field_sources")
			suite.Contains(fieldSources, "ID", "field_sources should include ID")
			suite.Contains(fieldSources, "Email", "field_sources should include Email")
			suite.Contains(fieldSources, "FullName", "field_sources should include FullName")

			// Verify the FullName field source resolves to User.Name
			fullNameSource, ok := fieldSources["FullName"].(map[string]interface{})
			suite.True(ok)
			suite.Equal("User", fullNameSource["model"])
			suite.Equal("Name", fullNameSource["field"])
			suite.Equal("String", fullNameSource["type"])
		}

		// Entity relationship addition: Organization added to UserProfile
		if change.Operation == diffdef.OperationAdd &&
			change.Type == diffdef.TypeRelationship &&
			change.Target["entity"] == "UserProfile" &&
			change.Target["relationship"] == "Organization" {
			foundEntityRelAdd = true
			suite.Equal(diffdef.ClassificationAdditive, change.Classification)

			// Verify entity_snapshot is present
			suite.NotNil(change.EntitySnapshot, "Entity relationship change should include entity_snapshot")

			// Verify snapshot contains the full post-change entity (including the new relationship)
			related, hasRelated := change.EntitySnapshot["related"].(map[string]interface{})
			suite.True(hasRelated, "entity_snapshot should contain 'related'")
			suite.Contains(related, "Organization", "snapshot related should include Organization")

			// Verify snapshot has the resolved block
			resolved, hasResolved := change.EntitySnapshot["resolved"].(map[string]interface{})
			suite.True(hasResolved, "entity_snapshot should contain 'resolved' block")
			suite.Equal("User", resolved["root_model"])
		}
	}

	suite.True(foundPhoneNumberAdd, "Should find PhoneNumber field addition")
	suite.True(foundEmailModify, "Should find Email field modification")
	suite.True(foundOrgAdd, "Should find Organization model addition")
	suite.True(foundEntityFieldAdd, "Should find FullName entity field addition with entity_snapshot")
	suite.True(foundEntityRelAdd, "Should find Organization entity relationship addition with entity_snapshot")
}

func (suite *CompileTestSuite) TestMorpheToMorpheDiff_EmptyRegistries() {
	basePath := filepath.Join(suite.TestDirPath, "empty-base")
	headPath := filepath.Join(suite.TestDirPath, "empty-head")
	outputPath := filepath.Join(suite.TestDirPath, "working-empty", "output.yaml")

	// Create empty directories
	os.MkdirAll(filepath.Join(basePath, "models"), 0755)
	os.MkdirAll(filepath.Join(headPath, "models"), 0755)
	defer os.RemoveAll(basePath)
	defer os.RemoveAll(headPath)

	workingDir := filepath.Dir(outputPath)
	os.MkdirAll(workingDir, 0755)
	defer os.RemoveAll(workingDir)

	config := compile.MorpheDiffConfig{
		BaseRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath: filepath.Join(basePath, "models"),
		},
		HeadRegistryConfig: rcfg.MorpheLoadRegistryConfig{
			RegistryModelsDirPath: filepath.Join(headPath, "models"),
		},
		OutputPath:  outputPath,
		BaseVersion: "base",
		HeadVersion: "head",
	}

	err := compile.MorpheToMorpheDiff(config)
	suite.Nil(err)

	// Read output
	data, err := os.ReadFile(outputPath)
	suite.Nil(err)

	var diffDoc diffdef.DiffDocument
	err = yaml.Unmarshal(data, &diffDoc)
	suite.Nil(err)

	// Should have no changes
	suite.Len(diffDoc.Changes, 0)
	suite.Equal(0, diffDoc.Metadata.Summary.TotalChanges)
}

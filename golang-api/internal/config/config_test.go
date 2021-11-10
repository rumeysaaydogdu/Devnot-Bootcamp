package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {

	t.Run("Config not found", func(t *testing.T) {

		_, err := GetAllValues(".", "NoFile")
		assert.Contains(t, err.Error(), "Not Found")
	})
	t.Run("Invalid Config file", func(t *testing.T) {

		p, f, cleanup := testFile(t, `
		mongoSettings:
		  databaseName: "Db"
		  uri:"mongodb://guest:guest@127.0.0.1:27017/admin?readPreference=primary"
		  arr: [1,2,
		`)
		defer cleanup()
		_, err := GetAllValues(p, f)
		assert.Contains(t, err.Error(), "While parsing config: yaml")
	})
	t.Run("Invalid Config file Should return unmarshal error", func(t *testing.T) {

		p, f, cleanup := testFile(t, `
mongoSettings:
  databaseName: "Db"
  uri: "mongodb://guest:guest@127.0.0.1:27017/admin?readPreference=primary"
  timeout: "abc"
`)
		defer cleanup()
		_, err := GetAllValues(p, f)
		assert.Contains(t, err.Error(), "Failed to unmarshal yaml file to configuration object")
	})
	t.Run("Valid Config found", func(t *testing.T) {
		p, f, cleanup := testFile(t, `
mongoSettings:
  databaseName: "Db"
  uri: "mongodb://guest:guest@127.0.0.1:27017/admin?readPreference=primary"
`)
		defer cleanup()
		cfg, err := GetAllValues(p, f)
		assert.NotNil(t, cfg)
		assert.Nil(t, err)
	})
}

func testFile(t *testing.T, contents string) (filePath, fileName string, cleanup func()) {

	t.Helper()
	filePath = os.TempDir()
	fileName = fmt.Sprintf("test-config-file-%v.yaml", time.Now().UnixNano())

	err := ioutil.WriteFile(path.Join(filePath, fileName), []byte(contents), 0644)

	if err != nil {
		t.Fatalf("failed to write test file %v", err)
	}

	return filePath, fileName, func() {
		os.Remove(path.Join(filePath, fileName))
	}
}

package catalog

import (
	"bytes"
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/tinyiot/thing-directory/wot"
)

const (
	envTestSchemaPath = "TEST_SCHEMA_PATH"
	defaultSchemaPath = "../wot/wot_td_schema.json"
)

type any = interface{}

var (
	TestSupportedBackends = map[string]bool{
		BackendMemory:  false,
		BackendLevelDB: true,
	}
	TestStorageType string
)

func loadSchema() error {
	if wot.LoadedJSONSchemas() {
		return nil
	}
	path := os.Getenv(envTestSchemaPath)
	if path == "" {
		path = defaultSchemaPath
	}
	return wot.LoadJSONSchemas([]string{path})
}

func serializedEqual(td1 ThingDescription, td2 ThingDescription) bool {
	// serialize to ease comparison of interfaces and concrete types
	tdBytes, _ := json.Marshal(td1)
	storedTDBytes, _ := json.Marshal(td2)

	return reflect.DeepEqual(tdBytes, storedTDBytes)
}

func copyMap(in, out interface{}) {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(in)
	json.NewDecoder(buf).Decode(out)
}

func TestMain(m *testing.M) {
	// run tests for each storage backend
	for b, supported := range TestSupportedBackends {
		if supported {
			TestStorageType = b
			if m.Run() == 1 {
				os.Exit(1)
			}
		}
	}
	os.Exit(0)
}

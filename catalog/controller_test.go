package catalog

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	uuid "github.com/satori/go.uuid"
)

func setup(t *testing.T) CatalogController {
	var (
		storage Storage
		tempDir = fmt.Sprintf("%s/thing-directory/test-%s-ldb",
			strings.Replace(os.TempDir(), "\\", "/", -1), uuid.NewV4())
	)

	err := loadSchema()
	if err != nil {
		t.Fatalf("error loading WoT Thing Description schema: %s", err)
	}

	switch TestStorageType {
	case BackendLevelDB:
		storage, err = NewLevelDBStorage(tempDir, nil)
		if err != nil {
			t.Fatalf("error creating leveldb storage: %s", err)
		}
	}

	controller, err := NewController(storage)
	if err != nil {
		storage.Close()
		t.Fatalf("error creating controller: %s", err)
	}

	t.Cleanup(func() {
		// t.Logf("Cleaning up...")
		controller.Stop()
		storage.Close()
		err = os.RemoveAll(tempDir) // Remove temp files
		if err != nil {
			t.Fatalf("error removing test files: %s", err)
		}
	})

	return controller
}

func TestControllerAdd(t *testing.T) {
	controller := setup(t)

	t.Run("user-defined ID", func(t *testing.T) {

		var td = map[string]any{
			"@context": "https://www.w3.org/2019/wot/td/v1",
			"id":       "urn:example:test/thing1",
			"title":    "example thing",
			"security": []string{"basic_sc"},
			"securityDefinitions": map[string]any{
				"basic_sc": map[string]string{
					"in":     "header",
					"scheme": "basic",
				},
			},
		}

		id, err := controller.add(td)
		if err != nil {
			t.Fatalf("Unexpected error on add: %s", err)
		}
		if id != td["id"] {
			t.Fatalf("User defined ID is not returned. Getting %s instead of %s\n", id, td["id"])
		}

		// add it again
		_, err = controller.add(td)
		if err == nil {
			t.Error("Didn't get any error when adding a service with non-unique id.")
		}
	})

	t.Run("system-generated ID", func(t *testing.T) {
		// System-generated id
		var td = map[string]any{
			"@context": "https://www.w3.org/2019/wot/td/v1",
			"title":    "example thing",
			"security": []string{"basic_sc"},
			"securityDefinitions": map[string]any{
				"basic_sc": map[string]string{
					"in":     "header",
					"scheme": "basic",
				},
			},
		}

		id, err := controller.add(td)
		if err != nil {
			t.Fatalf("Unexpected error on add: %s", err)
		}
		if !strings.HasPrefix(id, "urn:") {
			t.Fatalf("System-generated ID is not a URN. Got: %s\n", id)
		}
		_, err = uuid.FromString(strings.TrimPrefix(id, "urn:"))
		if err == nil {
			t.Fatalf("System-generated ID is not a uuid. Got: %s\n", id)
		}
	})
}

func TestControllerGet(t *testing.T) {
	controller := setup(t)

	var td = map[string]any{
		"@context": "https://www.w3.org/2019/wot/td/v1",
		"id":       "urn:example:test/thing1",
		"title":    "example thing",
		"security": []string{"basic_sc"},
		"securityDefinitions": map[string]any{
			"basic_sc": map[string]string{
				"in":     "header",
				"scheme": "basic",
			},
		},
	}

	id, err := controller.add(td)
	if err != nil {
		t.Fatalf("Unexpected error on add: %s", err)
	}

	t.Run("retrieve", func(t *testing.T) {
		storedTD, err := controller.get(id)
		if err != nil {
			t.Fatalf("Error retrieving: %s", err)
		}

		// set system-generated attributes
		storedTD["registration"] = td["registration"]

		if !serializedEqual(td, storedTD) {
			t.Fatalf("Added and retrieved TDs are not equal:\n Added:\n%v\n Retrieved:\n%v\n", td, storedTD)
		}
	})

	t.Run("retrieve non-existed", func(t *testing.T) {
		_, err := controller.get("some_id")
		if err != nil {
			switch err.(type) {
			case *NotFoundError:
			// good
			default:
				t.Fatalf("TD doesn't exist. Expected NotFoundError but got %s", err)
			}
		} else {
			t.Fatal("No error when retrieving a non-existed TD.")
		}
	})

}

func TestControllerUpdate(t *testing.T) {
	controller := setup(t)

	var td = map[string]any{
		"@context": "https://www.w3.org/2019/wot/td/v1",
		"id":       "urn:example:test/thing1",
		"title":    "example thing",
		"security": []string{"basic_sc"},
		"securityDefinitions": map[string]any{
			"basic_sc": map[string]string{
				"in":     "header",
				"scheme": "basic",
			},
		},
	}

	id, err := controller.add(td)
	if err != nil {
		t.Fatalf("Unexpected error on add: %s", err)
	}

	t.Run("update attributes", func(t *testing.T) {
		// Change
		td["title"] = "new title"
		td["description"] = "description of the thing"

		err = controller.update(id, td)
		if err != nil {
			t.Fatal("Error updating TD:", err.Error())
		}

		storedTD, err := controller.get(id)
		if err != nil {
			t.Fatal("Error retrieving TD:", err.Error())
		}

		// set system-generated attributes
		storedTD["registration"] = td["registration"]

		if !serializedEqual(td, storedTD) {
			t.Fatalf("Updates were not applied or returned:\n Expected:\n%v\n Retrieved:\n%v\n", td, storedTD)
		}
	})
}

func TestControllerDelete(t *testing.T) {
	controller := setup(t)

	var td = map[string]any{
		"@context": "https://www.w3.org/2019/wot/td/v1",
		"id":       "urn:example:test/thing1",
		"title":    "example thing",
		"security": []string{"basic_sc"},
		"securityDefinitions": map[string]any{
			"basic_sc": map[string]string{
				"in":     "header",
				"scheme": "basic",
			},
		},
	}

	id, err := controller.add(td)
	if err != nil {
		t.Fatalf("Error adding a TD: %s", err)
	}

	t.Run("delete", func(t *testing.T) {
		err = controller.delete(id)
		if err != nil {
			t.Fatalf("Error deleting TD: %s", err)
		}
	})

	t.Run("delete a deleted TD", func(t *testing.T) {
		err = controller.delete(id)
		if err != nil {
			switch err.(type) {
			case *NotFoundError:
			// good
			default:
				t.Fatalf("TD was deleted. Expected NotFoundError but got %s", err)
			}
		} else {
			t.Fatalf("No error when deleting a deleted TD: %s", err)
		}
	})

	t.Run("retrieve a deleted TD", func(t *testing.T) {
		_, err = controller.get(id)
		if err != nil {
			switch err.(type) {
			case *NotFoundError:
				// good
			default:
				t.Fatalf("TD was deleted. Expected NotFoundError but got %s", err)
			}
		} else {
			t.Fatal("No error when retrieving a deleted TD")
		}
	})
}

func TestControllerListPaginate(t *testing.T) {
	controller := setup(t)

	// add several entries
	var addedTDs []ThingDescription
	for i := 0; i < 5; i++ {
		var td = map[string]any{
			"@context": "https://www.w3.org/2019/wot/td/v1",
			"id":       "urn:example:test/thing_" + strconv.Itoa(i),
			"title":    "example thing",
			"security": []string{"basic_sc"},
			"securityDefinitions": map[string]any{
				"basic_sc": map[string]string{
					"in":     "header",
					"scheme": "basic",
				},
			},
		}

		tdCopy := make(map[string]any)
		copyMap(&td, &tdCopy)
		_, err := controller.add(tdCopy)
		if err != nil {
			t.Fatal("Error adding a TD:", err.Error())
		}

		addedTDs = append(addedTDs, td)
	}

	var list []ThingDescription

	// [0-3)
	TDs, err := controller.listPaginate(0, 3)
	if err != nil {
		t.Fatal("Error getting list of TDs:", err.Error())
	}
	if len(TDs) != 3 {
		t.Fatalf("Page has %d entries instead of 3", len(TDs))
	}
	list = append(list, TDs...)

	// [3-end)
	TDs, err = controller.listPaginate(3, 10)
	if err != nil {
		t.Fatal("Error getting list of TDs:", err.Error())
	}
	if len(TDs) != 2 {
		t.Fatalf("Page has %d entries instead of 2", len(TDs))
	}
	list = append(list, TDs...)

	if len(list) != 5 {
		t.Fatalf("Catalog contains %d entries instead of 5", len(list))
	}

	// compare added and collection
	for i, li := range list {
		// remove the whole registration object because it contains dynamic values
		delete(sd, "registration")
		if !serializedEqual(addedTDs[i], li) {
			t.Fatalf("TD added in catalog is different with the one listed:\n Added:\n%v\n Listed:\n%v\n",
				addedTDs[i], li)
		}
	}
}

func TestControllerFilter(t *testing.T) {
	controller := setup(t)

	for i := 0; i < 5; i++ {
		var td = map[string]any{
			"@context": "https://www.w3.org/2019/wot/td/v1",
			"id":       "urn:example:test/thing_" + strconv.Itoa(i),
			"title":    "example thing",
			"security": []string{"basic_sc"},
			"securityDefinitions": map[string]any{
				"basic_sc": map[string]string{
					"in":     "header",
					"scheme": "basic",
				},
			},
		}

		_, err := controller.add(td)
		if err != nil {
			t.Fatal("Error adding a TD:", err.Error())
		}
	}

	_, err := controller.add(map[string]any{
		"@context": "https://www.w3.org/2019/wot/td/v1",
		"id":       "urn:example:test/thing_x",
		"title":    "interesting thing",
		"security": []string{"basic_sc"},
		"securityDefinitions": map[string]any{
			"basic_sc": map[string]string{
				"in":     "header",
				"scheme": "basic",
			},
		},
	})
	if err != nil {
		t.Fatal("Error adding a TD:", err.Error())
	}

	_, err = controller.add(map[string]any{
		"@context": "https://www.w3.org/2019/wot/td/v1",
		"id":       "urn:example:test/thing_y",
		"title":    "interesting thing",
		"security": []string{"basic_sc"},
		"securityDefinitions": map[string]any{
			"basic_sc": map[string]string{
				"in":     "header",
				"scheme": "basic",
			},
		},
	})
	if err != nil {
		t.Fatal("Error adding a TD:", err.Error())
	}

	t.Run("JSONPath filter", func(t *testing.T) {
		b, err := controller.filterJSONPathBytes("$[?(@.title=='interesting thing')]")
		if err != nil {
			t.Fatal("Error filtering:", err.Error())
		}
		var TDs []ThingDescription
		err = json.Unmarshal(b, &TDs)
		if err != nil {
			t.Fatal("Error unmarshalling output:", err.Error())
		}
		if len(TDs) != 2 {
			t.Fatalf("Returned %d instead of 2 TDs when filtering based on title: \n%v", len(TDs), TDs)
		}
		for _, td := range TDs {
			if td["title"].(string) != "interesting thing" {
				t.Fatal("Wrong results when filtering based on title:\n", td)
			}
		}
	})

}

func TestControllerCleanExpired(t *testing.T) {

	// shorten controller's cleanup interval to test quickly
	controllerExpiryCleanupInterval = 1 * time.Second
	const wait = 3 * time.Second

	controller := setup(t)

	var td = ThingDescription{
		"@context": "https://www.w3.org/2019/wot/td/v1",
		"id":       "urn:example:test/thing1",
		"title":    "example thing",
		"security": []string{"basic_sc"},
		"securityDefinitions": map[string]any{
			"basic_sc": map[string]string{
				"in":     "header",
				"scheme": "basic",
			},
		},
		"registration": map[string]any{
			"ttl": 0.1,
		},
	}

	id, err := controller.add(td)
	if err != nil {
		t.Fatal("Error adding a TD:", err.Error())
	}

	time.Sleep(wait)

	_, err = controller.get(id)
	if err != nil {
		switch err.(type) {
		case *NotFoundError:
		// good
		default:
			t.Fatalf("Got an error other than NotFoundError when getting an expired TD: %s\n", err)
		}
	} else {
		t.Fatalf("Expired TD was not removed")
	}
}

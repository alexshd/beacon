package configmerge

import "testing"

func TestMerge(t *testing.T) {
	a := Config{"foo": "bar", "x": 1}
	b := Config{"baz": "qux", "y": 2}

	result := Merge(a, b)

	if result["foo"] != "bar" {
		t.Errorf("Expected foo=bar, got %v", result["foo"])
	}
	if result["baz"] != "qux" {
		t.Errorf("Expected baz=qux, got %v", result["baz"])
	}
}

func TestMergeOverride(t *testing.T) {
	a := Config{"key": "value1"}
	b := Config{"key": "value2"}

	result := Merge(a, b)

	if result["key"] != "value2" {
		t.Errorf("Expected key=value2, got %v", result["key"])
	}
}

func TestDeepMerge(t *testing.T) {
	a := Config{
		"db": map[string]interface{}{
			"host": "localhost",
			"port": 5432,
		},
	}
	b := Config{
		"db": map[string]interface{}{
			"port": 3306,
			"user": "admin",
		},
	}

	result := DeepMerge(a, b)

	dbRaw := result["db"]
	var db map[string]interface{}
	switch v := dbRaw.(type) {
	case Config:
		db = map[string]interface{}(v)
	case map[string]interface{}:
		db = v
	}

	if db["host"] != "localhost" {
		t.Errorf("Expected host=localhost, got %v", db["host"])
	}
	if db["port"] != 3306 {
		t.Errorf("Expected port=3306, got %v", db["port"])
	}
	if db["user"] != "admin" {
		t.Errorf("Expected user=admin, got %v", db["user"])
	}
}

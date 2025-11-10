package configmerge

import "testing"

func FuzzMerge(f *testing.F) {
	f.Add("key1", "val1", "key2", "val2")
	f.Add("", "", "key", "value")
	f.Add("x", "1", "x", "2")

	f.Fuzz(func(t *testing.T, k1, v1, k2, v2 string) {
		a := Config{k1: v1}
		b := Config{k2: v2}

		result := Merge(a, b)

		if len(result) > 2 {
			t.Errorf("Merge created too many keys\nInputs: a=%v, b=%v\nResult: %v", a, b, result)
		}

		// Check immutability
		if len(a) != 1 {
			t.Errorf("Merge mutated input a\nOriginal should have 1 key, now has %d\nInputs: k1=%q, v1=%q", len(a), k1, v1)
		}
		if len(b) != 1 {
			t.Errorf("Merge mutated input b\nOriginal should have 1 key, now has %d\nInputs: k2=%q, v2=%q", len(b), k2, v2)
		}
	})
}

func FuzzDeepMerge(f *testing.F) {
	f.Add("db", "host", "localhost", "db", "port", "5432")

	f.Fuzz(func(t *testing.T, k1, k2, v1, k3, k4, v2 string) {
		a := Config{
			k1: map[string]any{k2: v1},
		}
		b := Config{
			k3: map[string]any{k4: v2},
		}

		result := DeepMerge(a, b)

		if result == nil {
			t.Errorf("DeepMerge returned nil\nInputs: a[%q][%q]=%q, b[%q][%q]=%q", k1, k2, v1, k3, k4, v2)
		}

		// Check original wasn't mutated
		if a[k1] == nil {
			t.Errorf("DeepMerge corrupted input a\nInputs: k1=%q, k2=%q, v1=%q", k1, k2, v1)
		}
	})
}

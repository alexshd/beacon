package configmerge

import "maps"

// Config represents a simple configuration with nested values
type Config map[string]any

// ConfigWrapper wraps Config to make it work with lawtest
//
// lawtest requires types to be comparable (for equality checks in property tests).
// Go maps are NOT comparable - you cannot use == or != with them.
// However, pointers ARE comparable in Go.
//
// Solution: Wrap the map in a struct and use pointers (*ConfigWrapper).
// This allows lawtest's property-based testing to work with map-based configs.
type ConfigWrapper struct {
	config Config
}

// NewConfigWrapper creates a new wrapper
func NewConfigWrapper(c Config) *ConfigWrapper {
	return &ConfigWrapper{config: c}
}

// Unwrap returns the underlying Config
func (w *ConfigWrapper) Unwrap() Config {
	return w.config
}

// WrapMerge wraps Merge for lawtest compatibility
// Takes wrapped configs, unwraps them, merges, and re-wraps the result
func WrapMerge(a, b *ConfigWrapper) *ConfigWrapper {
	return NewConfigWrapper(Merge(a.config, b.config))
}

// WrapDeepMerge wraps DeepMerge for lawtest compatibility
// Takes wrapped configs, unwraps them, deep merges, and re-wraps the result
func WrapDeepMerge(a, b *ConfigWrapper) *ConfigWrapper {
	return NewConfigWrapper(DeepMerge(a.config, b.config))
}

// Merge combines two configs, with the second config's values taking precedence
func Merge(a, b Config) Config {
	result := make(Config)

	// Copy all from a
	maps.Copy(result, a)

	// Override with b
	maps.Copy(result, b)

	return result
}

// DeepMerge combines configs recursively
func DeepMerge(a, b Config) Config {
	result := make(Config)

	// Copy all from a
	maps.Copy(result, a)

	// Merge with b
	for k, v := range b {
		if existing, ok := result[k]; ok {
			// If both are maps, merge recursively
			if existingMap, ok := existing.(map[string]any); ok {
				if vMap, ok := v.(map[string]any); ok {
					result[k] = DeepMerge(Config(existingMap), Config(vMap))
					continue
				}
			}
		}
		// Otherwise, b wins
		result[k] = v
	}

	return result
}

package configmerge

import (
	"testing"

	"github.com/alexshd/lawtest"
)

// TestMergeImmutability uses lawtest to verify Merge doesn't mutate inputs
//
// Note: We use ConfigWrapper because lawtest requires comparable types.
// Config is a map, which is NOT comparable in Go.
// ConfigWrapper uses pointers, which ARE comparable.
func TestMergeImmutability(t *testing.T) {
	gen := func() *ConfigWrapper {
		return NewConfigWrapper(Config{
			lawtest.StringGen(5)(): lawtest.StringGen(10)(),
		})
	}

	lawtest.ImmutableOp(t, WrapMerge, gen)
}

// TestMergeAssociativity uses lawtest to verify (a+b)+c = a+(b+c)
func TestMergeAssociativity(t *testing.T) {
	gen := func() *ConfigWrapper {
		return NewConfigWrapper(Config{
			lawtest.StringGen(5)(): lawtest.StringGen(10)(),
		})
	}

	lawtest.Associative(t, WrapMerge, gen)
}

func TestDeepMergeImmutability(t *testing.T) {
	gen := func() *ConfigWrapper {
		return NewConfigWrapper(Config{
			"nested": map[string]interface{}{
				lawtest.StringGen(5)(): lawtest.StringGen(10)(),
			},
		})
	}

	lawtest.ImmutableOp(t, WrapDeepMerge, gen)
}

func TestDeepMergeAssociativity(t *testing.T) {
	gen := func() *ConfigWrapper {
		return NewConfigWrapper(Config{
			"nested": map[string]interface{}{
				lawtest.StringGen(5)(): lawtest.StringGen(10)(),
			},
		})
	}

	lawtest.Associative(t, WrapDeepMerge, gen)
}

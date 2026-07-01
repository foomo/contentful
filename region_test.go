package contentful

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseRegion(t *testing.T) {
	t.Run("empty string is valid", func(t *testing.T) {
		r, err := ParseRegion("")
		require.NoError(t, err)
		assert.Equal(t, Region(""), r)
	})
	t.Run("us is valid", func(t *testing.T) {
		r, err := ParseRegion("us")
		require.NoError(t, err)
		assert.Equal(t, RegionUS, r)
	})
	t.Run("eu is valid", func(t *testing.T) {
		r, err := ParseRegion("eu")
		require.NoError(t, err)
		assert.Equal(t, RegionEU, r)
	})
	t.Run("uppercase EU is rejected", func(t *testing.T) {
		_, err := ParseRegion("EU")
		require.Error(t, err)
	})
	t.Run("unknown region is rejected", func(t *testing.T) {
		_, err := ParseRegion("ap")
		require.Error(t, err)
	})
}

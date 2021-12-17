package helpr

import (
	"testing"

	"github.com/skker/helpr/str"
	"github.com/stretchr/testify/assert"
)

func TestStr(t *testing.T) {
	t.Log("\n\nTest str package...")

	s := "This is My String to Search"

	substr := "sea"
	result := str.CaseInsensitiveSearch(s, substr)
	assert.True(t, result, "String ("+s+") should contain "+substr)

	substr = "land"
	result = str.CaseInsensitiveSearch(s, substr)
	assert.False(t, result, "String ("+s+") should not contain "+substr)
}

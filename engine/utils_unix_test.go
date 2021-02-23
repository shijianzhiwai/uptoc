package engine

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var excludes = []string{"test1", "test2", "test3/sub/", "/test4/sub/"}
var dirPath = []string{"./a/", "a/"}

func TestShouldExclude(t *testing.T) {
	for _, d := range dirPath {
		// true
		assert.Equal(t, shouldExclude(d, "a/test1/abc", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test1/abc.log", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test2/abc", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test2/abc.log", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test3/sub/abc", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test3/sub/abc.log", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test4/sub/abc", excludes), true)
		assert.Equal(t, shouldExclude(d, "a/test4/sub/abc.log", excludes), true)

		// false
		assert.Equal(t, shouldExclude(d, "a/test1_tail/abc.log", excludes), false)
		assert.Equal(t, shouldExclude(d, "a/test2_tail/abc.log", excludes), false)
	}
}

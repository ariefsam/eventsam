package idgenerator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	ids := map[string]bool{}
	for i := 0; i < 1000; i++ {
		t.Run("Test generate ke "+fmt.Sprint(i), func(t *testing.T) {
			gotId := Generate()
			assert.False(t, ids[gotId])
		})
	}
}

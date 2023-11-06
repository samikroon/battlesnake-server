package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GetEnvString(t *testing.T) {
	val := getEnvString("NON_EXISTING", "default")
	assert.Equal(t, "default", val)

	os.Setenv("TEST_ENV", "exists")
	val = getEnvString("TEST_ENV", "default")
	assert.Equal(t, "exists", val)
}

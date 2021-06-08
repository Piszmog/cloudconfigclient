package cloudconfigclient_test

import (
	"github.com/Piszmog/cloudconfigclient/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotFoundError_Error(t *testing.T) {
	err := cloudconfigclient.NotFoundError{}
	assert.Equal(t, "failed to find resource", err.Error())
}

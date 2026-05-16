package ctxlib

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	testRequestID := "test-request-id"
	ctx := context.Background()
	ctx = SetRequestID(ctx, testRequestID)

	reqID := RequestID(ctx)
	assert.Equal(t, testRequestID, reqID)
}

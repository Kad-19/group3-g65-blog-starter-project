package database

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitMongoDB(t *testing.T) {
	// This is an integration test and requires a running MongoDB instance.
	// Skip this test if you don't have a local MongoDB instance.
	t.Skip("Skipping MongoDB integration test")

	client := InitMongoDB()
	assert.NotNil(t, client)

	err := client.Ping(context.Background(), nil)
	assert.NoError(t, err)
}

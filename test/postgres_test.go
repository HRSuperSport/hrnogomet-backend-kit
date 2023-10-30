package test_test

import (
	"context"
	"github.com/hrsupersport/hrnogomet-backend-kit/test"
	"testing"
)

func TestCreatePostgresContainer(t *testing.T) {
	ctx := context.Background()
	c := test.SetupPostgresDB(ctx, "user", "test", "testdb")
	defer c.TeardownPostgresDB()
}

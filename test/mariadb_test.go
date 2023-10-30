package test_test

import (
	"context"
	"github.com/hrsupersport/hrnogomet-backend-kit/test"
	"testing"
)

func TestCreateMariadbContainer(t *testing.T) {
	ctx := context.Background()
	c := test.SetupMariaDB(ctx, "user", "test", "test", "testdb")
	defer c.TeardownMariaDB()
}

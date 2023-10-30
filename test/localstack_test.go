package test_test

import (
	"context"
	"github.com/hrsupersport/hrnogomet-backend-kit/test"
	"testing"
)

func TestCreateLocalstackContainer(t *testing.T) {
	ctx := context.Background()
	c := test.SetupLocalstack(ctx)
	defer c.TeardownLocalstack()
}

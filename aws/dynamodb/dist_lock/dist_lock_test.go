package dist_lock_test

import (
	"cirello.io/dynamolock/v2"
	"context"
	"errors"
	"github.com/hrsupersport/hrnogomet-backend-kit/aws"
	"github.com/hrsupersport/hrnogomet-backend-kit/aws/dynamodb/dist_lock"
	"github.com/hrsupersport/hrnogomet-backend-kit/logging"
	"github.com/hrsupersport/hrnogomet-backend-kit/test"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const (
	dynamoDistLockTableName = "distTable"
)

func TstDistributedLockWithDynamoDbHappyPath(t *testing.T) {
	t.Helper()
	logging.ConfigureDefaultLoggingSetup("")
	ctx := context.Background()
	c := test.SetupLocalstack(ctx)
	ctx = aws.SetCustomAwsEndpoint(ctx, c.URI)
	defer c.TeardownLocalstack()

	err := dist_lock.CreateLockTable(ctx, dynamoDistLockTableName)
	assert.Nil(t, err)

	chnl := make(chan struct{}, 2)

	go func() {
		lock, err := dist_lock.NewDistributedDynamodbLock(ctx, "foo", 3*time.Second, 1*time.Second, dynamoDistLockTableName)
		//Acquired lock! If I die, my lock will expire in 3 seconds.
		// Otherwise, I will hold it until I stop heart beating.
		log.Info().Str("goroutine", "1").Msg("NewDistributedDynamodbLock done")
		assert.Nil(t, err)
		assert.NotNil(t, lock)
		if lock == nil {
			chnl <- struct{}{}
			return
		}

		log.Info().Str("goroutine", "1").Msg("acquired lock")
		time.Sleep(time.Second * 8)

		if err := lock.ReleaseLock(); err != nil {
			log.Error().Err(err).Msg("lock.ReleaseLock() error")
		}
		log.Info().Str("goroutine", "1").Msg("released lock")

		chnl <- struct{}{}
	}()

	time.Sleep(1000 * time.Millisecond) // wait to make sure goroutine 1 acquires lock as first
	testWasHere := 0

	go func() {
		// will wait for lock 3 sec + 2 * 3 sec = at most 9 sec before timing out
		// so we should get lock because go routine 1 releases at 8 sec after retrieving it
		lock, err := dist_lock.NewDistributedDynamodbLock(ctx, "foo", 3*time.Second, 1*time.Second, dynamoDistLockTableName)
		log.Info().Str("goroutine", "2").Msg("NewDistributedDynamodbLock done")
		assert.Nil(t, err)
		assert.NotNil(t, lock)
		if lock == nil {
			chnl <- struct{}{}
			return
		}

		log.Info().Str("goroutine", "2").Msg("acquired lock")
		time.Sleep(time.Second * 1)

		if err := lock.ReleaseLock(); err != nil {
			log.Error().Err(err).Msg("lock.ReleaseLock() error")
		}
		log.Info().Str("goroutine", "2").Msg("released lock")

		testWasHere++

		chnl <- struct{}{}
	}()

	<-chnl
	<-chnl
	assert.Equal(t, 1, testWasHere)
}

func TstDistributedLockWithDynamoDbUnhappyPathTimeout(t *testing.T) {
	t.Helper()
	logging.ConfigureDefaultLoggingSetup("")
	ctx := context.Background()
	c := test.SetupLocalstack(ctx)
	ctx = aws.SetCustomAwsEndpoint(ctx, c.URI)
	defer c.TeardownLocalstack()

	err := dist_lock.CreateLockTable(ctx, dynamoDistLockTableName)
	assert.Nil(t, err)

	chnl := make(chan struct{}, 2)

	go func() {
		lock, err := dist_lock.NewDistributedDynamodbLock(ctx, "foo", 3*time.Second, 1*time.Second, dynamoDistLockTableName)
		log.Info().Str("goroutine", "1").Msg("NewDistributedDynamodbLock done")
		assert.Nil(t, err)
		assert.NotNil(t, lock)
		if lock == nil {
			chnl <- struct{}{}
			return
		}

		log.Info().Str("goroutine", "1").Msg("acquired lock")
		time.Sleep(time.Second * 12) // HOLD LOCK FOR 12 SECONDS!

		if err := lock.ReleaseLock(); err != nil {
			log.Error().Err(err).Msg("lock.ReleaseLock() error")
		}
		log.Info().Str("goroutine", "1").Msg("released lock")

		chnl <- struct{}{}
	}()

	time.Sleep(1000 * time.Millisecond) // wait to make sure goroutine 1 acquires lock as first
	testWasHere := 0

	go func() {
		var lockNotGrantedErr *dynamolock.LockNotGrantedError

		// lock retrieval must always fail we are waiting up to 9 seconds while go routine 1 will hold the lock for 12 seconds!
		lock, err := dist_lock.NewDistributedDynamodbLock(ctx, "foo", 3*time.Second, 1*time.Second, dynamoDistLockTableName)
		log.Info().Str("goroutine", "2").Msg("NewDistributedDynamodbLock done")
		if err != nil {
			log.Error().Err(err).Str("goroutine", "2").Msg("")
		}
		assert.NotNil(t, err)
		assert.Nil(t, lock)

		if errors.As(err, &lockNotGrantedErr) {
			testWasHere++
		}

		chnl <- struct{}{}
	}()

	<-chnl
	<-chnl
	assert.Equal(t, 1, testWasHere)
}

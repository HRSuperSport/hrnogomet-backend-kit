package dist_lock

/*
Sample implementation of distributed lock using dynamodb table
So far not used, rather POC and demonstration of using AWS SDK for GO V2
*/

import (
	"cirello.io/dynamolock/v2"
	"context"
	"errors"
	aws_sdk "github.com/aws/aws-sdk-go-v2/aws"
	dynamodb_types "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/hrsupersport/hrnogomet-backend-kit/aws"
	"github.com/hrsupersport/hrnogomet-backend-kit/constants"
	"time"
)

type distributedDynamodbLock struct {
	dynamoTableName string
	client          *dynamolock.Client
	lock            *dynamolock.Lock
}

// NewLockClient creates new client for distributed lock operations
func NewLockClient(ctx context.Context, leaseDuration time.Duration, leaseExtensionHeartbeatInterval time.Duration, dynamoTableName string) (*dynamolock.Client, error) {
	dynamodbClient, err := aws.CreateDynamodbClient(ctx, constants.AwsDefaultRegion)
	if err != nil {
		return nil, err
	}

	return dynamolock.New(dynamodbClient,
		dynamoTableName,
		dynamolock.WithLeaseDuration(leaseDuration),
		dynamolock.WithHeartbeatPeriod(leaseExtensionHeartbeatInterval),
	)
}

// AcquireLock retrieves new lock utilizing provided client
func AcquireLock(lockKey string, client *dynamolock.Client, lockData []byte, additionalTimeToWaitForLock time.Duration, failIfLocked bool) (*dynamolock.Lock, error) {
	// https://github.com/cirello-io/dynamolock/issues/199
	// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/dynamodb#Client.UpdateTimeToLive
	// TTL can be used automatically expire lock items in database and let dynamodb delete them after some time
	// we are deleting lock rows  (see below WithDeleteLockOnRelease) but still it is good to have it here
	//attrs := make(map[string]*dynamodb.AttributeValue, 0)
	//attrTTL := &dynamodb.AttributeValue{N: aws_sdk.String(fmt.Sprintf("%d", time.Now().Add(24*time.Hour).Unix()))}
	//attrs["TTL"] = attrTTL
	// causing false positive during linting (need to wait for fixed version of golangci/golangci-lint-action):
	// Error: internal/distributed_lock/distributed_lock.go:57:2: G602: Potentially accessing slice out of bounds (gosec)
	// attrs["TTL"] = attrTTL

	if failIfLocked {
		return client.AcquireLock(lockKey,
			dynamolock.FailIfLocked(),
			dynamolock.WithData(lockData),
			dynamolock.WithAdditionalTimeToWaitForLock(additionalTimeToWaitForLock),
			///dynamolock.WithAdditionalAttributes(attrs),
			dynamolock.WithDeleteLockOnRelease(), // delete row from dynamodb table once lock is released
		)
	} else {
		return client.AcquireLock(lockKey,
			dynamolock.WithData(lockData),
			dynamolock.WithAdditionalTimeToWaitForLock(additionalTimeToWaitForLock),
			//dynamolock.WithAdditionalAttributes(attrs),
			dynamolock.WithDeleteLockOnRelease(), // delete row from dynamodb table once lock is released
		)
	}
}

// NewDistributedDynamodbLock is convenience function combining creation of locking client and lock into single step.
// Created lock holds no data because this feature will be mostly not needed.
// If needed call NewLockClient and AcquireLock explicitly instead.
// lock created this call will wait at most (2*leaseDuration time  + leaseDuration) to get the lock before timing out
func NewDistributedDynamodbLock(ctx context.Context, lockKey string, leaseDuration time.Duration, leaseExtensionHeartbeatInterval time.Duration, dynamoTableName string) (*distributedDynamodbLock, error) {
	if client, errNewClient := NewLockClient(ctx, leaseDuration, leaseExtensionHeartbeatInterval, dynamoTableName); errNewClient != nil {
		return nil, errNewClient
	} else {
		if lock, errLock := AcquireLock(lockKey, client, nil, leaseDuration*2, false); errLock != nil {
			return nil, errLock
		} else {
			return &distributedDynamodbLock{
				dynamoTableName,
				client,
				lock,
			}, nil
		}
	}
}

// ReleaseLock will release distributed lock and also close locking client
func (l *distributedDynamodbLock) ReleaseLock() error {
	if l.client != nil && l.lock != nil {
		if releasedOk, err := l.client.ReleaseLock(l.lock); err != nil {
			return err
		} else {
			if !releasedOk {
				return errors.New("ReleaseLock: unable to release lock")
			}
		}
	}

	if l.client != nil {
		if err := l.client.Close(); err != nil {
			return err
		}
	}
	return nil
}

// CreateLockTable should be used only within unit tests testing the lock
// in real environments locking table is provisioned in advance using terraform scripts
func CreateLockTable(ctx context.Context, dynamoTableName string) error {
	if client, errNewClient := NewLockClient(ctx, 3*time.Second, 1*time.Second, dynamoTableName); errNewClient != nil {
		return errNewClient
	} else {
		defer client.Close()
		_, errCreateTable := client.CreateTable(dynamoTableName,
			dynamolock.WithProvisionedThroughput(&dynamodb_types.ProvisionedThroughput{
				ReadCapacityUnits:  aws_sdk.Int64(5),
				WriteCapacityUnits: aws_sdk.Int64(5),
			}),
			dynamolock.WithCustomPartitionKeyName("key"))
		if errCreateTable != nil {
			return errCreateTable
		}

		return nil
	}
}

package dynamodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/bitcoin-sv/arc/metamorph/metamorph_api"
	"github.com/bitcoin-sv/arc/metamorph/store"
	"github.com/libsv/go-p2p/chaincfg/chainhash"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
	"github.com/ordishs/gocore"
)

type DynamoDB struct {
	dynamoCli *dynamodb.Client
}

const ISO8601 = "2006-01-02T15:04:05.999Z"

func New() (store.MetamorphStore, error) {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return &DynamoDB{}, err
	}

	// create dynamodb client
	dynamodbClient := &DynamoDB{dynamoCli: dynamodb.NewFromConfig(cfg)}

	// create table if not exists
	exists, err := dynamodbClient.TableExists("transactions")
	if err != nil {
		return dynamodbClient, err
	} else if !exists {
		dynamodbClient.CreateTransactionsTable()
	}

	// create table if not exists
	exists, err = dynamodbClient.TableExists("blocks")
	if err != nil {
		fmt.Println(err)
		return dynamodbClient, err
	} else if !exists {
		fmt.Println("creating ....")
		dynamodbClient.CreateBlocksTable()
	}

	data := []byte(fmt.Sprintf("Hello world %d-%d-%d", 23, 2, time.Now().UnixMilli()))
	hash := chainhash.DoubleHashH(data)
	err = dynamodbClient.Set(context.TODO(), nil, &store.StoreData{Hash: &hash})
	if err != nil {
		fmt.Println(err)
	}

	// Using the Config value, create the DynamoDB client
	return dynamodbClient, nil
}

// CREATE TABLE IF NOT EXISTS transactions (
// 	hash BLOB PRIMARY KEY,
// 	stored_at TEXT,
// 	announced_at TEXT,
// 	mined_at TEXT,
// 	status INTEGER,
// 	block_height BIGINT,
// 	block_hash BLOB,
// 	callback_url TEXT,
// 	callback_token TEXT,
// 	merkle_proof TEXT,
// 	reject_reason TEXT,
// 	raw_tx BLOB
// 	);

func (ddb *DynamoDB) TableExists(tableName string) (bool, error) {
	// Build the request with its input parameters
	resp, err := ddb.dynamoCli.ListTables(context.TODO(), &dynamodb.ListTablesInput{})
	if err != nil {
		return false, err
	}

	for _, tn := range resp.TableNames {
		if tn == tableName {
			return true, nil
		}
	}

	return false, nil
}

// CreateTransactionsTable creates a DynamoDB table for storing metamorph user transactions
func (ddb *DynamoDB) CreateTransactionsTable() error {
	// construct primary key and create table
	if _, err := ddb.dynamoCli.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("hash"),
				AttributeType: types.ScalarAttributeTypeB,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("hash"),
				KeyType:       types.KeyTypeHash,
			}},
		TableName: aws.String("transactions"),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}); err != nil {
		return err
	}

	return nil
}

// CreateBlocksTable creates a DynamoDB table for storing blocks
func (ddb *DynamoDB) CreateBlocksTable() error {
	// construct primary key and create table
	if _, err := ddb.dynamoCli.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("hash"),
				AttributeType: types.ScalarAttributeTypeB,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("hash"),
				KeyType:       types.KeyTypeHash,
			}},
		TableName: aws.String("blocks"),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}); err != nil {
		return err
	}

	return nil
}

func (ddb *DynamoDB) Get(ctx context.Context, key []byte) (*store.StoreData, error) {
	// config log and tracing
	startNanos := time.Now().UnixNano()
	defer func() {
		gocore.NewStat("mtm_store_sql").NewStat("Get").AddTime(startNanos)
	}()
	span, _ := opentracing.StartSpanFromContext(ctx, "dynamodb:Get")
	defer span.Finish()

	// delete the item
	val, err := attributevalue.MarshalMap(key)
	if err != nil {
		return nil, err
	}

	response, err := ddb.dynamoCli.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: val, TableName: aws.String("transactions"),
	})
	if err != nil {
		return nil, err
	}

	var transaction store.StoreData
	err = attributevalue.UnmarshalMap(response.Item, &transaction)
	if err != nil {
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return nil, err
	}

	return &transaction, nil
}

func (ddb *DynamoDB) Set(ctx context.Context, key []byte, value *store.StoreData) error {
	// setup log and tracing
	startNanos := time.Now().UnixNano()
	defer func() {
		gocore.NewStat("mtm_store_dynamodb").NewStat("Set").AddTime(startNanos)
	}()
	span, _ := opentracing.StartSpanFromContext(ctx, "dynamodb:Set")
	defer span.Finish()

	// marshal input value for new entry
	item, err := attributevalue.MarshalMap(value)
	if err != nil {
		return err
	}

	// put item into table
	_, err = ddb.dynamoCli.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("transactions"), Item: item,
	})

	if err != nil {
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return nil
	}

	return nil
}

func (ddb *DynamoDB) Del(ctx context.Context, key []byte) error {
	// setup log and tracing
	startNanos := time.Now().UnixNano()
	defer func() {
		gocore.NewStat("mtm_store_dynamodb").NewStat("Del").AddTime(startNanos)
	}()
	span, _ := opentracing.StartSpanFromContext(ctx, "dynamodb:Del")
	defer span.Finish()

	// delete the item
	val, err := attributevalue.MarshalMap(key)
	if err != nil {
		return err
	}

	_, err = ddb.dynamoCli.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String("transactions"), Key: val,
	})

	if err != nil {
		span.SetTag(string(ext.Error), true)
		span.LogFields(log.Error(err))
		return err
	}

	return nil

}

func (dynamodb *DynamoDB) GetUnmined(_ context.Context, callback func(s *store.StoreData)) error {
	return nil
}

func (dynamodb *DynamoDB) UpdateStatus(ctx context.Context, hash *chainhash.Hash, status metamorph_api.Status, rejectReason string) error {
	return nil
}

func (dynamodb *DynamoDB) UpdateMined(ctx context.Context, hash *chainhash.Hash, blockHash *chainhash.Hash, blockHeight uint64) error {
	return nil
}

func (dynamodb *DynamoDB) GetBlockProcessed(ctx context.Context, blockHash *chainhash.Hash) (*time.Time, error) {
	return nil, nil
}

func (dynamodb *DynamoDB) SetBlockProcessed(ctx context.Context, blockHash *chainhash.Hash) error {
	return nil
}

func (ddb *DynamoDB) Close(ctx context.Context) error {
	startNanos := time.Now().UnixNano()
	defer func() {
		gocore.NewStat("mtm_store_sql").NewStat("Close").AddTime(startNanos)
	}()
	span, _ := opentracing.StartSpanFromContext(ctx, "dynamodb:Close")
	defer span.Finish()

	ctx.Done()
	return nil
}

package mongodb

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
)

type MongodbConnection struct {
	Uri            string
	DbName         string
	CollectionName string
}

type MongodbConnectionResult struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

type MongodbConnectionInterface interface {
	InsertOne(data interface{}, insertOptions *options.InsertOneOptions) (*mongo.InsertOneResult, error)
	InsertMany(data []interface{}, insertOptions *options.InsertManyOptions) (*mongo.InsertManyResult, error)
	UpdateOne(filter interface{}, data interface{}, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error)
	UpsertOne(filter interface{}, data interface{}) (*mongo.UpdateResult, error)
	UpdateMany(filter interface{}, data []interface{}, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error)
	UpdateByID(id interface{}, data []interface{}, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error)
	DeleteOne(filter interface{}, deleteOptions *options.DeleteOptions) (*mongo.DeleteResult, error)
	DeleteMany(filter interface{}, deleteOptions *options.DeleteOptions) (*mongo.DeleteResult, error)
	FindOne(filter interface{}, result interface{}, findOptions *options.FindOneOptions) error
	Find(filter interface{}, result interface{}, findOptions *options.FindOptions) error
}

var ctx context.Context
var ctxCancel context.CancelFunc

func setContext() {
	// ctx will be used to set deadline for process, here deadline will be of 30 seconds.
	ctx, ctxCancel = context.WithTimeout(context.Background(), 30*time.Second)
}

func connectRetrieveCollection(mongodbConnection MongodbConnection) (*MongodbConnectionResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	setContext()
	defer ctxCancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbConnection.Uri))
	if err != nil {
		logger.Error("Error in Connecting to MongoDB: ", zap.Error(err))
		return nil, err
	}

	collection := client.Database(mongodbConnection.DbName, options.Database().SetReadPreference(readpref.SecondaryPreferred())).Collection(mongodbConnection.CollectionName)
	if err != nil {
		logger.Error("error in get collection "+mongodbConnection.CollectionName, zap.Error(err))
		return nil, err
	}

	return &MongodbConnectionResult{Client: client, Collection: collection}, nil
}

func disconnectClient(client *mongo.Client) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	if err := client.Disconnect(ctx); err != nil {
		logger.Error("Error in disconnecting MongoClient: ", zap.Error(err))
	}
}

func (mongoconnection MongodbConnection) InsertOne(data interface{}, insertOptions *options.InsertOneOptions) (*mongo.InsertOneResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	result, err := mongodbConnectionResult.Collection.InsertOne(ctx, data, insertOptions)
	if err != nil {
		logger.Error("Error in Insert(One) data: ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) InsertMany(data []interface{}, insertOptions *options.InsertManyOptions) (*mongo.InsertManyResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	result, err := mongodbConnectionResult.Collection.InsertMany(ctx, data, insertOptions)
	if err != nil {
		logger.Error("Error in Insert(Many) data: ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) UpdateOne(filter interface{}, data interface{}, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	bsonData := bson.M{"$set": data}

	result, err := mongodbConnectionResult.Collection.UpdateOne(ctx, filter, bsonData, updateOptions)
	if err != nil {
		logger.Error("Error in Update(One) data: ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) UpsertOne(filter interface{}, data interface{}) (*mongo.UpdateResult, error) {
	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	setContext()
	defer ctxCancel()

	bsonData := bson.M{"$set": data}

	result, err := mongodbConnectionResult.Collection.UpdateOne(ctx, filter, bsonData, options.Update().SetUpsert(true))
	if err != nil {
		logger.Error("Error in upsert data: ", zap.Error(err))
		return result, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) UpdateMany(filter interface{}, data []interface{}, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	bsonData := bson.M{"$set": data}

	result, err := mongodbConnectionResult.Collection.UpdateMany(ctx, filter, bsonData, updateOptions)
	if err != nil {
		logger.Error("Error in Update(Many) data: ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) UpdateByID(id interface{}, data []interface{}, updateOptions *options.UpdateOptions) (*mongo.UpdateResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	bsonData := bson.M{"$set": data}
	result, err := mongodbConnectionResult.Collection.UpdateByID(ctx, id, bsonData, updateOptions)
	if err != nil {
		logger.Error("Error in Update(ById) ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) DeleteOne(filter interface{}, deleteOptions *options.DeleteOptions) (*mongo.DeleteResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	result, err := mongodbConnectionResult.Collection.DeleteOne(ctx, filter, deleteOptions)
	if err != nil {
		logger.Error("Error in Delete(One) data: ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) DeleteMany(filter interface{}, deleteOptions *options.DeleteOptions) (*mongo.DeleteResult, error) {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return nil, err
	}

	setContext()
	defer ctxCancel()

	result, err := mongodbConnectionResult.Collection.DeleteMany(ctx, filter, deleteOptions)
	if err != nil {
		logger.Error("Error in Delete(Many) data: ", zap.Error(err))
		return nil, err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return result, nil
}

func (mongoconnection MongodbConnection) FindOne(filter interface{}, result interface{}, findOptions *options.FindOneOptions) error {

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return err
	}

	setContext()
	defer ctxCancel()

	singleResult := mongodbConnectionResult.Collection.FindOne(ctx, filter, findOptions)
	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	err = singleResult.Decode(result)
	if err != nil {
		return err
	}

	return nil
}

func (mongoconnection MongodbConnection) Find(filter interface{}, result interface{}, findOptions *options.FindOptions) error {

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	mongodbConnectionResult, err := connectRetrieveCollection(mongoconnection)
	if err != nil {
		return err
	}

	setContext()
	defer ctxCancel()

	cursor, err := mongodbConnectionResult.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		logger.Error("Error in find data: ", zap.Error(err))
		return err
	}
	// Defer closing the cursor
	defer cursor.Close(ctx)

	err = cursor.All(ctx, result)

	// Check for any errors during cursor iteration
	if err != nil {
		return err
	}

	defer func() {
		disconnectClient(mongodbConnectionResult.Client)
	}()

	return nil
}

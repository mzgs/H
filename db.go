package H

import (
	"context"
	"errors"
	"log"
	"reflect"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

var DEBUG_DATABASE = false

func InitMongoDB(dbName string) MongoDBHelper {

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal("MONGODB init failed!", err)
	}
	MongoDB = client.Database(dbName)

	return MongoDBHelper{}

}

type MongoDBHelper struct {
	Model             interface{}
	CurrentCollection *mongo.Collection
	Filters           interface{}
	limit             int
	skip              int
	sort              []string
}

func checkDBError(title string, err error) {
	if DEBUG_DATABASE {
		if err != nil {
			log.Println("MongoDB Error", title, err)
		}
	}
}

func (m MongoDBHelper) Database() *mongo.Database {
	return MongoDB
}

func (m MongoDBHelper) InsertOne(i interface{}) (primitive.ObjectID, error) {
	r, err := m.Col(i).InsertOne(context.TODO(), i)

	checkDBError("InsertOne", err)

	return r.InsertedID.(primitive.ObjectID), err
}

func (m MongoDBHelper) InsertMany(documents []interface{}) error {
	if len(documents) == 0 {
		return errors.New("documents size 0")
	}
	_, err := m.Col(documents[0]).InsertMany(context.TODO(), documents)
	checkDBError("InsertMany", err)
	return err
}

func (m MongoDBHelper) UpdateOne(i interface{}) error {

	r := reflect.ValueOf(i)
	id := reflect.Indirect(r).FieldByName("ID").Interface()

	_, err := m.Col(i).UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": i})
	checkDBError("UpdateOne", err)
	return err
}

func (m MongoDBHelper) UpdateMany(i interface{}) {

}

func (m MongoDBHelper) DeleteMany(i interface{}, filters interface{}) error {

	_, err := m.Col(i).DeleteMany(context.TODO(), filters)
	checkDBError("DeleteMany", err)
	return err
}

func (m MongoDBHelper) DeleteOne(i interface{}, filters interface{}) error {
	_, err := m.Col(i).DeleteOne(context.TODO(), filters)
	checkDBError("DeleteOne", err)
	return err
}

func (m MongoDBHelper) Find(i interface{}, filters ...interface{}) MongoDBHelper {

	var f interface{} = bson.M{}

	if len(filters) != 0 {
		f = filters[0]
	}

	m.CurrentCollection = m.Col(i)
	m.Filters = f
	m.Model = i

	return m
}

func (m MongoDBHelper) Limit(limit int) MongoDBHelper {
	m.limit = limit
	return m
}

func (m MongoDBHelper) Skip(skip int) MongoDBHelper {
	m.skip = skip
	return m
}

func (m MongoDBHelper) Sort(sort ...string) MongoDBHelper {
	m.sort = sort
	return m
}

func (m MongoDBHelper) One() error {
	return m.CurrentCollection.FindOne(context.Background(), m.Filters, getFindOneOptions(m)).Decode(m.Model)
}

func (m MongoDBHelper) All() error {

	cur, err := m.CurrentCollection.Find(context.TODO(), m.Filters, getFindOptions(m))
	if err != nil {
		return err
	}
	return cur.All(context.TODO(), m.Model)

}

func getFindOptions(m MongoDBHelper) *options.FindOptions {
	options := options.Find()

	sorts := bson.D{}
	for _, s := range m.sort {

		direction := 1
		if strings.HasPrefix(s, "-") {
			direction = -1
		}
		s = strings.TrimLeft(s, "-")
		sorts = append(sorts, bson.E{Key: s, Value: direction})

	}

	options.SetSort(sorts)
	options.SetSkip(int64(m.skip))
	options.SetLimit(int64(m.limit))
	return options
}

func getFindOneOptions(m MongoDBHelper) *options.FindOneOptions {
	options := options.FindOne()

	sorts := bson.D{}
	for _, s := range m.sort {

		direction := 1
		if strings.HasPrefix(s, "-") {
			direction = -1
		}
		s = strings.TrimLeft(s, "-")
		sorts = append(sorts, bson.E{Key: s, Value: direction})

	}

	options.SetSort(sorts)
	options.SetSkip(int64(m.skip))

	return options
}

func (m MongoDBHelper) FindAll(i interface{}, opts ...*options.FindOptions) error {
	return nil
}

func (m MongoDBHelper) FindOne(i interface{}, filters interface{}) error {

	err := m.Col(i).FindOne(context.TODO(), filters).Decode(i)
	checkDBError("FindOne", err)
	return err

}

func (m MongoDBHelper) CreateSearchIndex(i interface{}, keys []string, weights bson.M) {

	keysM := bson.M{}
	for _, value := range keys {
		keysM[value] = "text"
	}

	opts := options.Index().SetDefaultLanguage("turkish").SetWeights(weights)

	_, err := m.Col(i).Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    keysM,
			Options: opts,
		},
	})
	checkDBError("CreateSearchIndex", err)
}

func (m MongoDBHelper) DropAllIndex(i interface{}) {

	_, err := m.Col(i).Indexes().DropAll(context.Background())
	if err != nil {
		checkDBError("DropAllIndex", err)

	}
}

func (m MongoDBHelper) Col(i interface{}) *mongo.Collection {
	return MongoDB.Collection(typeName(i))
}

func (m MongoDBHelper) FindByID(i interface{}, id primitive.ObjectID) error {

	return m.Col(i).FindOne(context.TODO(), bson.M{"_id": id}).Decode(i)
}

func (m MongoDBHelper) Count(i interface{}, filters interface{}, opts ...*options.CountOptions) int {
	count, err := m.Col(i).CountDocuments(context.TODO(), filters, opts...)
	checkDBError("count", err)
	return int(count)
}

func (m MongoDBHelper) Aggregate(colType interface{}, pipeline mongo.Pipeline, opts ...*options.AggregateOptions) []bson.M {

	var showsWithInfo []bson.M

	cur, err := m.Col(colType).Aggregate(context.TODO(), pipeline, opts...)
	cur.All(context.TODO(), &showsWithInfo)

	checkDBError("aggragate", err)

	return showsWithInfo

}

func (m MongoDBHelper) AggregateJson(results interface{}, json string, opts ...*options.AggregateOptions) error {

	var pipeline mongo.Pipeline
	err := bson.UnmarshalExtJSON([]byte(json), false, &pipeline)

	if err != nil {
		return err
	}

	cur, err := m.Col(results).Aggregate(context.TODO(), pipeline, opts...)
	if err != nil {
		return err
	}
	err = cur.All(context.TODO(), results)

	return err
}

func (m MongoDBHelper) ObjectId(id string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(id)
	return oid
}

func typeName(i interface{}) string {
	tp := reflect.TypeOf(i)

	if tp.Kind() == reflect.Ptr {
		tp = tp.Elem()
	}

	if tp.Kind() == reflect.Slice {
		tp = tp.Elem()
	}
	return tp.Name()
}

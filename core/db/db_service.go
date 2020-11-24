package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

func connect() (*mongo.Client, context.Context) {
	credential := options.Credential{
		Username: os.Getenv("user"),
		Password: os.Getenv("password"),
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("base_url")).SetAuth(credential))

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("Connect db ok")

	return client, ctx
}

func GetRow(find bson.M, coll string, sel bson.M) map[string]interface{} {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	collection := db.Database(os.Getenv("database")).Collection(coll)

	cur, err := collection.Find(ctx, find, options.Find().SetProjection(sel))

	if err != nil {
		log.Panic(err)
	}
	defer cur.Close(ctx)

	result := make(map[string]interface{})

	var json bson.M
	for cur.Next(ctx) {
		if err = cur.Decode(&json); err != nil {
			log.Panic(err)
		}
		result["data"] = json["data"]
	}

	return result
}

func GetRows(find bson.M, coll string, sel bson.M) []map[string]interface{} {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	collection := db.Database(os.Getenv("database")).Collection(coll)

	cur, err := collection.Find(ctx, find, options.Find().SetProjection(sel))

	if err != nil {
		log.Panic(err)
	}
	defer cur.Close(ctx)

	jsonArray := make([]map[string]interface{}, 0, 0)
	for cur.Next(ctx) {
		jsonObject := make(map[string]interface{})

		if err = cur.Decode(&jsonObject); err != nil {
			log.Panic(err)
		}

		fmt.Println(jsonObject)
		jsonArray = append(jsonArray, jsonObject)
	}

	fmt.Println(jsonArray)

	return jsonArray
}

func GetRowAgtegate(find bson.M, coll string, coll_join string, where bson.A, show bson.M) bson.M {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	collection := db.Database(os.Getenv("database")).Collection(coll)

	pipeline := []bson.M{
		{"$lookup": bson.M{
			"from": coll_join,
			"let":  find,
			"pipeline": bson.A{
				bson.M{"$match": bson.M{"$expr": bson.M{"$and": where}}},
				bson.M{"$project": show}},
			"as": "data"}}}

	cur, err := collection.Aggregate(ctx, pipeline)

	if err != nil {
		log.Panic(err)
	}
	defer cur.Close(ctx)

	var json bson.M
	for cur.Next(ctx) {
		if err = cur.Decode(&json); err != nil {
			log.Panic(err)
		}

		fmt.Println(json)

		jArray := json["data"].(bson.A)

		if len(jArray) != 0 {
			break
		}
	}

	return json
}

func InserRow(json bson.M, coll string) int {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	collection := db.Database(os.Getenv("database")).Collection(coll)

	result, err := collection.InsertOne(ctx, json)

	if err != nil {
		log.Panic(err)
		return -1
	}
	fmt.Printf("Inserted %v documents into episode collection!\n", result)

	return 1
}

func InserRows(json []interface{}) int {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	collection := db.Database(os.Getenv("database")).Collection("")

	_, err := collection.InsertMany(ctx, json)

	if err != nil {
		log.Panic(err)
		return -1
	}

	return 1
}

func UpdateRow(id bson.M, set bson.D, coll string) int {
	db, ctx := connect()
	defer db.Disconnect(ctx)

	collection := db.Database(os.Getenv("database")).Collection(coll)

	update, err := collection.UpdateOne(ctx, id, set)

	if err != nil {
		log.Panic(err)
		return -1
	}

	fmt.Printf("Updated %v Documents!\n", update.ModifiedCount)

	if update.ModifiedCount > 0 {
		return 1
	} else {
		return 0
	}

}

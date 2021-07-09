// Package corrlinker provides ...
package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"runtime"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	MongoDBhost string `json:"mongoDBhost"`
	MongoDBport string `json:"mongoDBport"`
	MongoDB     string `json:"mongoDB"`
	MongoDBuser string `json:"mongoDBuser"`
	MongoDBpass string `json:"mongoDBpass"`
}

type Person struct {
	Name string
	Age  int
	City string
}

func readConfig() (*config, error) {
	_, filePath, _, _ := runtime.Caller(0)
	pwd := filePath[:len(filePath)-12]
	txt, err := ioutil.ReadFile(pwd + "/config.json")
	if err != nil {
		return nil, err
	}
	var cfg = new(config)
	if err := json.Unmarshal(txt, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Database(data []interface{}) {
	cfg, err := readConfig()
	if err != nil {
		log.Fatalf("Could not read config file: %v", err)
	}

	credential := options.Credential{
		Username: cfg.MongoDBuser,
		Password: cfg.MongoDBpass,
	}
	clientOpts := options.Client().ApplyURI("mongodb://" + cfg.MongoDBhost + ":" + cfg.MongoDBport).SetAuth(credential)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB!")

	collection := client.Database("corrlinker").Collection("messages")
	entries := data
	opts := options.InsertMany().SetOrdered(false)
	insertManyResult, err := collection.InsertMany(context.TODO(), entries, opts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)
}

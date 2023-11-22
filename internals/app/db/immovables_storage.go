package db

import (
	"cian-parse/internals/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

type ImmovablesStorage struct {
	collection *mongo.Collection
	log        *slog.Logger
}

func NewImmovablesStorage(database *mongo.Database, collection string, logger *slog.Logger) *ImmovablesStorage {
	storage := new(ImmovablesStorage)
	storage.collection = database.Collection(collection)
	storage.log = logger
	return storage
}

func (s *ImmovablesStorage) Create(ctx context.Context, immovable models.Immovable) (string, error) {
	s.log.Info("create immovable")
	result, err := s.collection.InsertOne(ctx, immovable)
	if err != nil {
		return "", fmt.Errorf("failed to create user due to error: %v", err)
	}

	oid, ok := result.InsertedID.(primitive.ObjectID)
	if ok {
		return oid.Hex(), nil
	}

	return "", fmt.Errorf("failed to convert object to hex %s", oid)
}

func (s *ImmovablesStorage) FindOne(ctx context.Context, id string) (immovable models.Immovable, err error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return immovable, fmt.Errorf("failed to convert hex: %s to objectid", id)
	}

	filter := bson.M{"_id": oid}

	result := s.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return immovable, fmt.Errorf("error to find one immovable by id: %s", id)
	}
	if err = result.Decode(&immovable); err != nil {
		return immovable, fmt.Errorf("error to find one immovable by id from DB: %s", id)
	}

	return immovable, nil
}

func (s *ImmovablesStorage) FindAll(ctx context.Context) (immovables []models.Immovable, err error) {
	result, err := s.collection.Find(ctx, bson.M{})
	if result.Err() != nil {
		return immovables, fmt.Errorf("failed to find all immovables: %v", err)
	}

	if err = result.All(ctx, &immovables); err != nil {
		return immovables, fmt.Errorf("failed to read all document: %v", err)
	}

	return immovables, nil
}

func (s *ImmovablesStorage) Update(ctx context.Context, id string, immovable models.Immovable) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex: %s to objectid", id)
	}

	filter := bson.M{"_id": oid}

	immovablesBytes, err := bson.Marshal(immovable)
	if err != nil {
		return fmt.Errorf("failed to marshal immovable: %v", err)
	}

	var updateImmovableObj bson.M
	err = bson.Unmarshal(immovablesBytes, &updateImmovableObj)
	if err != nil {
		return fmt.Errorf("failed to unmarshal immovable bytes: %v", err)
	}

	delete(updateImmovableObj, "_id")

	update := bson.M{
		"$set": updateImmovableObj,
	}

	result, err := s.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to execute update immovable: %v", err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("immovable not found")
	}
	s.log.Info(fmt.Sprintf("Matched %d documents, Modified %d documents", result.MatchedCount, result.ModifiedCount))

	return nil
}

func (s *ImmovablesStorage) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("failed to convert hex: %s to objectid", id)
	}

	filter := bson.M{"_id": oid}

	result, err := s.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to execute query")
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("immovable not found")
	}
	s.log.Info(fmt.Sprintf("Delete %d documents"), result.DeletedCount)

	return nil
}

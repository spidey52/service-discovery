package repository

import (
	"context"
	"time"

	"github.com/spidey52/service-discovery/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoRepo struct {
	coll *mongo.Collection
}

// NewMongoRepo creates a new repository
func NewMongoRepo(coll *mongo.Collection) *MongoRepo {
	return &MongoRepo{coll: coll}
}

func (r *MongoRepo) Register(ctx context.Context, inst models.Instance) error {
	inst.LastHeartbeat = time.Now().UTC()
	inst.Health = "UP"
	filter := bson.M{"serviceName": inst.ServiceName, "id": inst.ID}
	update := bson.M{"$set": inst}
	_, err := r.coll.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *MongoRepo) UpdateHeartbeat(ctx context.Context, serviceName, id string) error {
	filter := bson.M{"serviceName": serviceName, "id": id}
	update := bson.M{"$set": bson.M{"lastHeartbeat": time.Now().UTC(), "health": "UP"}}
	res, err := r.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

func (r *MongoRepo) Find(ctx context.Context, serviceName, mode string, metadata map[string]interface{}, aliveOnly bool, ttl time.Duration) ([]models.Instance, error) {
	filter := bson.M{}
	if serviceName != "" {
		filter["serviceName"] = serviceName
	}
	if mode != "" {
		filter["mode"] = mode
	}
	for k, v := range metadata {
		filter["metadata."+k] = v
	}
	if aliveOnly {
		cutoff := time.Now().Add(-ttl)
		filter["lastHeartbeat"] = bson.M{"$gte": cutoff}
	}

	cur, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var instances []models.Instance
	if err := cur.All(ctx, &instances); err != nil {
		return nil, err
	}
	return instances, nil
}

func (r *MongoRepo) CleanupDead(ctx context.Context, ttl time.Duration) error {
	cutoff := time.Now().Add(-ttl)
	_, err := r.coll.DeleteMany(ctx, bson.M{"lastHeartbeat": bson.M{"$lt": cutoff}})
	return err
}

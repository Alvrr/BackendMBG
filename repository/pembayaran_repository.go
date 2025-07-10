package repository

import (
	"backend/config"
	"backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func pembayaranCol() *mongo.Collection {
	return config.PembayaranCollection
}

func GetAllPembayaran() ([]models.Pembayaran, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := pembayaranCol().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var list []models.Pembayaran
	for cursor.Next(ctx) {
		var p models.Pembayaran
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		list = append(list, p)
	}

	return list, nil
}

func GetPembayaranByID(id string) (*models.Pembayaran, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var pembayaran models.Pembayaran
	err := pembayaranCol().FindOne(ctx, bson.M{"_id": id}).Decode(&pembayaran)
	if err != nil {
		return nil, err
	}
	return &pembayaran, nil
}

func CreatePembayaran(p models.Pembayaran) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return pembayaranCol().InsertOne(ctx, p)
}

func UpdatePembayaran(id string, p models.Pembayaran) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"$set": p}
	return pembayaranCol().UpdateByID(ctx, id, update)
}

func DeletePembayaran(id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return pembayaranCol().DeleteOne(ctx, bson.M{"_id": id})
}

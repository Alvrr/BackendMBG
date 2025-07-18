package repository

import (
	"backend/config"
	"backend/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func userCol() *mongo.Collection {
	return config.UserCollection
}

// List all drivers
func GetAllDrivers() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"role": "driver"}
	cursor, err := userCol().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var drivers []models.User
	for cursor.Next(ctx) {
		var u models.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		drivers = append(drivers, u)
	}
	return drivers, nil
}

// CRUD Karyawan (User, kecuali role admin)
func GetAllKaryawan() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"role": bson.M{"$ne": "admin"}}
	cursor, err := userCol().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []models.User
	for cursor.Next(ctx) {
		var u models.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func GetKaryawanByID(id string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User
	err := userCol().FindOne(ctx, bson.M{"id": id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func CreateKaryawan(user models.User) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return userCol().InsertOne(ctx, user)
}

func UpdateKaryawan(id string, user models.User) (*mongo.UpdateResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{"$set": user}
	return userCol().UpdateByID(ctx, id, update)
}

func DeleteKaryawan(id string) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return userCol().DeleteOne(ctx, bson.M{"id": id})
}

// Update status aktif/nonaktif karyawan
func UpdateKaryawanStatus(id string, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	update := bson.M{"$set": bson.M{"status": status}}
	_, err := userCol().UpdateOne(ctx, bson.M{"id": id}, update)
	return err
}

// Get all active karyawan (selain admin)
func GetActiveKaryawan() ([]models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	filter := bson.M{"status": "aktif", "role": bson.M{"$ne": "admin"}}
	cursor, err := userCol().Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var users []models.User
	for cursor.Next(ctx) {
		var u models.User
		if err := cursor.Decode(&u); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

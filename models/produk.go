package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Produk struct {
	ID         string             `json:"id" bson:"_id"`
	NamaProduk string             `json:"nama_produk" bson:"nama_produk" validate:"required"`
	Kategori   string             `json:"kategori" bson:"kategori" validate:"required"`
	Harga      int                `json:"harga" bson:"harga" validate:"required"`
	Stok       int                `json:"stok" bson:"stok" validate:"required"`
	Deskripsi  string             `json:"deskripsi" bson:"deskripsi"`
	CreatedAt  primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`
}

package models

type Pembayaran struct {
	ID          string       `json:"id" bson:"_id"`
	IDPelanggan string       `json:"id_pelanggan" bson:"id_pelanggan" validate:"required"`
	Tanggal     string       `json:"tanggal,omitempty" bson:"tanggal"`
	Produk      []ItemProduk `json:"produk" bson:"produk" validate:"required"`
	TotalBayar  int          `json:"total_bayar" bson:"total_bayar" validate:"required"`
}

type ItemProduk struct {
	IDProduk   string `json:"id_produk" bson:"id_produk"`
	NamaProduk string `json:"nama_produk" bson:"nama_produk" validate:"required"`
	Jumlah     int    `json:"jumlah" bson:"jumlah" validate:"required"`
	Subtotal   int    `json:"subtotal" bson:"subtotal" validate:"required"`
}

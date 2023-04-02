package models

// model
type Product struct {
	ID       int    `json:"-"`
	Kategori string `json:"kategori"`
	Urun     string `json:"urun"`
	Aciklama string `json:"aciklama"`
	Fiyat    string `json:"fiyat"`
}

type Response struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    []Product
}

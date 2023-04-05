package controller

import (
	Config "RestApiGo/Config"
	models "RestApiGo/model"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
)

// localhost:9000/happmoons
// tüm veriler
func Handler(w http.ResponseWriter, r *http.Request) {
	db := Config.Connect()
	result, err := db.Query("SELECT * FROM happymoons3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	products := []models.Product{}
	for result.Next() {
		var product models.Product
		err = result.Scan(&product.ID, &product.Kategori, &product.Aciklama, &product.Urun, &product.Fiyat)
		if err != nil {
			panic(err)
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")

	// Marshal the products to JSON format
	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	// Write the JSON response
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
	//http.HandleFunc("/happymoons", handler)
	//http.ListenAndServe(":9000", nil)
	//fmt.Println("web server")

}

// http://localhost:9000/happymoons/csv
// csv olarak indirmek
func CsvHandler(w http.ResponseWriter, r *http.Request) {
	db := Config.Connect()
	result, err := db.Query("SELECT * FROM happymoons3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment;filename=happymoons3.csv")

	// Create a CSV writer and write the header row
	writer := csv.NewWriter(w)
	columnNames := []string{"ID", "Kategori", "Aciklama", "Urun", "Fiyat"}
	writer.Write(columnNames)

	for result.Next() {
		var product models.Product // models paketindeki Product yapısını kullanın
		err = result.Scan(&product.ID, &product.Kategori, &product.Aciklama, &product.Urun, &product.Fiyat)
		if err != nil {
			panic(err)
		}

		row := []string{strconv.Itoa(product.ID), product.Kategori, product.Aciklama, product.Urun, product.Fiyat}
		writer.Write(row)
	}

	writer.Flush()
}

func ExHandler(w http.ResponseWriter, r *http.Request) {
	// Parse query string for excluded columns
	queryParams := r.URL.Query()
	excludedColumns := map[string]bool{}
	if excluded, ok := queryParams["ex"]; ok {
		for _, column := range excluded {
			excludedColumns[column] = true
		}
	}

	db := Config.Connect()
	result, err := db.Query("SELECT * FROM happymoons3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	columns, err := result.Columns()
	if err != nil {
		panic(err)
	}

	products := []map[string]interface{}{}
	for result.Next() {
		row := make(map[string]interface{})
		values := make([]interface{}, len(columns))
		for i := range values {
			values[i] = new(interface{})
		}

		err = result.Scan(values...)
		if err != nil {
			panic(err)
		}

		for i, col := range columns {
			if !excludedColumns[col] {
				row[col] = *(values[i].(*interface{}))
			}
		}

		products = append(products, row)
	}

	// Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the products to JSON format
	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	// Write the JSON response
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}

func InHandler(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()
	includeCols := q["in"]

	db := Config.Connect()
	result, err := db.Query("SELECT * FROM happymoons3")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	products := []models.Product{}
	for result.Next() {
		var product models.Product
		err = result.Scan(&product.ID, &product.Kategori, &product.Aciklama, &product.Urun, &product.Fiyat)
		if err != nil {
			panic(err)
		}

		// Filter out the columns not included in the query parameter
		filteredProduct := models.Product{}
		for _, col := range includeCols {
			switch col {
			case "kategori":
				filteredProduct.Kategori = product.Kategori
			case "urun":
				filteredProduct.Urun = product.Urun
			case "aciklama":
				filteredProduct.Aciklama = product.Aciklama
			case "fiyat":
				filteredProduct.Fiyat = product.Fiyat
			}
		}

		products = append(products, filteredProduct)
	}

	// Set the content type header to JSON
	w.Header().Set("Content-Type", "application/json")

	// Marshal the products to JSON format
	jsonBytes, err := json.Marshal(products)
	if err != nil {
		panic(err)
	}

	// Write the JSON response
	_, err = w.Write(jsonBytes)
	if err != nil {
		panic(err)
	}
}

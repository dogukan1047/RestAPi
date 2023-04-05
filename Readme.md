# GOLANG REST API

<p>This code is a REST API in Go language. The API uses a MySQL database to store data about products. The API has four endpoints:
/happymoons: returns all products in JSON format
/happymoons/csv: returns all products in CSV format for download
/happymoons?in={column_name}: returns all products excluding the specified columns in JSON format
/happymoons?ex={column_name}: returns only the specified columns in JSON format.
The main function sets up a new router using the Gorilla mux library and registers the handlers for each endpoint. It then listens on a random port on the local machine and starts the server. The server will continue to run until the program is terminated.

The Config package defines a function for connecting to the MySQL database. It uses the "github.com/go-sql-driver/mysql" driver to establish a connection to the database using the connection string stored in the "DB_CONNECTION_STRING" environment variable.

The controller package defines the four handlers for the API endpoints. The Handler function returns all products in JSON format. The CsvHandler function returns all products in CSV format for download. The ExHandler function returns all products excluding the specified columns in JSON format. The InHandler function returns only the specified columns in JSON format. <p>

```main.go```
<p>This code is a basic REST API written in Go language. It sets up a server using the Gorilla mux router and listens to incoming requests on certain endpoints, namely "/happymoons" and "/happymoons/csv". The API has two handlers, one for handling GET requests on "/happymoons" and another for handling GET requests on "/happymoons/csv".

Additionally, the code creates two subrouters for the "/happymoons" endpoint, one for handling GET requests with the "ex" query parameter and another for handling GET requests with the "in" query parameter. These subrouters point to separate handler functions in the controller package.

The code uses the net package to randomly select an available port to listen on and outputs this port to the console. The server runs in a separate goroutine to allow the program to continue running until manually stopped.

Lastly, the code includes a defer statement to gracefully shut down the server when the program terminates.<p>

```Config.go```

<p>This code defines a function named Connect() that returns a pointer to a SQL database connection. The function reads the database connection string from an environment variable named DB_CONNECTION_STRING. If the variable is not set, the function prints an error message and exits the program. If the variable is set, the function uses the go-sql-driver/mysql package to open a MySQL database connection using the connection string, and returns the connection.<p>


```contoller.go```

<p>The code includes several endpoints in a RESTful API.<p>

# Handler()

<p>The Handler() function responds to GET requests coming from "localhost:9000/happymoons". This request retrieves all the data from the database and returns a response in JSON format.<p>

```

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
```

# CSV HANDLER()

<p>The CsvHandler() function responds to GET requests coming from "localhost:9000/happymoons/csv". This request retrieves all the data from the database and returns a response in CSV format.<p>

```
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
```

# EXHANDLER()
<p>The ExHandler() function responds to GET requests coming from "localhost:9000/happymoons/ex". This request retrieves all the data from the database, excluding certain columns, and returns a response in JSON format.<p>

```
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
```
# INHANDLER()
<p>The InHandler() function responds to GET requests coming from "localhost:9000/happymoons/in". This request retrieves all the data from the database, including certain columns, and returns a response in JSON format.<p>

```
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
```

<p>The code is contained within a package called RestApiGo, which includes the Config and model packages. The Config package is used to configure the database connection, while the model package represents rows in the database using a structure called Product.<p>


# Base image that our image will inherit

```
FROM golang:1.2.0

WORKDIR /app

RUN apk update &&
apk add libc-dev &&
apk add gcc &&
apk add make

COPY go.mod ./
COPY go.sum ./
``` 

# Download dependencies using go mod files

```
RUN go mod download && go mod verify
RUN go get github.com/githubnemo/CompileDaemon
```

# Copy all application files

```
COPY . .
COPY ./entrypoint.sh /entrypoint.sh

ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for /entrypoint.sh

ENTRYPOINT [ "sh", "/entrypoint.sh" ]
```

# MySQL image

```
FROM mysql:8.0.2

ENV MYSQL_ROOT_PASSWORD=104725
ENV MYSQL_DATABASE=MENU

COPY migration.sql /docker-entrypoint-initdb.d/

EXPOSE 3306

version: '3.7'

services:
db:
container_name: "test_db"
platform: linux/x86_64
image: mysql:8.0.2
environment:
MYSQL_ROOT_PASSWORD: "104725"
MYSQL_DATABASE: "MENU"
volumes:
- mysql_data:/var/lib/mysql
ports:
- "3306:3306"
restart: always
command: --default-authentication-plugin=mysql_native_password

api:
container_name: "test_api"
build:
context: .
dockerfile: api.Dockerfile
ports:
- "8080:8080"
depends_on:
- db
environment:
DATABASE_HOST: db
DATABASE_PORT: 3306
DATABASE_USER: root
DATABASE_PASSWORD: 104725
DATABASE_NAME: MENU
command: >
sh -c "wait-for ${DATABASE_HOST}:${DATABASE_PORT} --
CompileDaemon --build='go build -o main main.go'
--command=./main"

volumes:
mysql_data:
 
```
<p>The code defines a Dockerfile for a Go-based API and a MySQL database. The Dockerfile is based on the golang:1.2.0 image, and sets up the working directory, installs dependencies and compiles the API code.

It also copies the entrypoint.sh script that runs at container startup, sets environment variables, and exposes the necessary port for the database.

The file also includes a MySQL database image and configuration options for the database service. Finally, it sets up the API service with configuration options, including the database host and port, and a command that uses a wait-for script to ensure that the database is <p>
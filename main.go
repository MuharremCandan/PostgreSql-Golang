package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// use the package to use .env module ın your code automaticlly
	"LibraryProject/services"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	// Loading evironment variables
	dialect := os.Getenv("DIALECT")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	//Database connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s  port=%s", host, user, dbName, password, dbPort)

	//Openning connection to database

	db, err := gorm.Open(dialect, dbURI)

	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database !")
	}

	//close connection to database when the  main function finishes
	defer db.Close()

	services.Db = db

	// Make migrations  to the database  if they have not already  been created

	db.AutoMigrate(&services.Person{})
	db.AutoMigrate(&services.Book{})

	// will make table at database

	db.Create(&services.Personn)
	for idx := range services.Books {
		db.Create(&services.Books[idx])
	}

	// API routes

	router := mux.NewRouter()

	// Controller for People
	router.HandleFunc("/people", services.GetPeople).Methods("GET")
	router.HandleFunc("/person/{id}", services.GetPerson).Methods("GET")
	router.HandleFunc("/create/person", services.CreatePerson).Methods("POST")
	router.HandleFunc("/delete/person/{id}", services.DeletePerson).Methods("DELETE")

	//Controller for Books

	router.HandleFunc("/books", services.GetBooks).Methods("GET")
	router.HandleFunc("/book/{id}", services.GetBook).Methods("GET")
	router.HandleFunc("/create/book", services.CreateBook).Methods("POST")
	router.HandleFunc("delete/book/{id}", services.DeleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8090", router))

}

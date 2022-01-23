package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	// use the package to use .env module ın your code automaticlly
	_ "github.com/joho/godotenv/autoload"
)

type Person struct {
	gorm.Model

	Name  string
	Email string `gorm:"typevarchar(100);uniqe_index"`
	Book  []Book
	Age   int
}

type Book struct {
	gorm.Model

	Title      string
	Author     string
	CallNumber int `gorm:"uniqe_index"`
	PersonID   int
}

var (
	person = &Person{
		Name:  "Muharrem",
		Email: "1muharremcandan@gmail.com",
		Age:   21,
		Book:  []Book{},
	}

	books = []Book{
		{
			Title:      "Dan Vinci Code",
			Author:     "Dawn Brown",
			CallNumber: 1,
			PersonID:   1,
		},
		{
			Title:      "Kısa Kes",
			Author:     "Leigh Russell",
			CallNumber: 2,
			PersonID:   1,
		},
	}
)

var Db *gorm.DB
var err error

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

	// Make migrations  to the database  if they have not already  been created

	db.AutoMigrate(&Person{})
	db.AutoMigrate(&Book{})

	// will make table at database

	db.Create(&person)
	for idx := range books {
		db.Create(&books[idx])
	}

	// API routes

	Db = db
	router := mux.NewRouter()

	// Controller for People
	router.HandleFunc("/people", getPeople).Methods("GET")
	router.HandleFunc("/person/{id}", getPerson).Methods("GET")
	router.HandleFunc("/create/person", createPerson).Methods("POST")
	router.HandleFunc("/delete/person/{id}", deletePerson).Methods("DELETE")

	//Controller for Books

	router.HandleFunc("/books", getBooks).Methods("GET")
	router.HandleFunc("/book/{id}", getBook).Methods("GET")
	router.HandleFunc("/create/book", createBook).Methods("POST")
	router.HandleFunc("delete/book/{id}", deleteBook).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8090", router))

}

//API Controllers for People
func getPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	Db.Find(&people)

	json.NewEncoder(w).Encode(people)
}
func getPerson(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	var books []Book

	Db.First(&person, params["id"])
	Db.Model(&person).Related(&books)
	person.Book = books

	json.NewEncoder(rw).Encode(person)
}

func createPerson(rw http.ResponseWriter, r *http.Request) {
	var person Person

	json.NewDecoder(r.Body).Decode(&person)

	createdPerson := Db.Create(&person)
	err = createdPerson.Error

	if err != nil {
		json.NewEncoder(rw).Encode(err)
	} else {
		json.NewEncoder(rw).Encode(&person)
	}

}

func deletePerson(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	Db.First(&person, params["id"])

	Db.Delete(&person)

	json.NewEncoder(rw).Encode(&person)

}

// Cotroller for Books
func getBooks(rw http.ResponseWriter, r *http.Request) {
	var books []Book

	Db.Find(&books)

	json.NewEncoder(rw).Encode(&books)
}

func getBook(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book
	Db.First(&book, params["id"])

	json.NewEncoder(rw).Encode(&book)

}

func createBook(rw http.ResponseWriter, r *http.Request) {
	var book Book

	json.NewDecoder(r.Body).Decode(&book)

	createdBook := Db.Create(&book)
	err = createdBook.Error
	if err != nil {
		json.NewEncoder(rw).Encode(err)
	} else {
		json.NewEncoder(rw).Encode(&book)
	}

}

func deleteBook(rw http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var book Book

	Db.First(&book, params["id"])
	Db.Delete(&person)

	json.NewEncoder(rw).Encode(&book)

}

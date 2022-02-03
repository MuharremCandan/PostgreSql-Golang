package services

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
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
	Personn = &Person{
		Name:  "Muharrem",
		Email: "1muharremcandan@gmail.com",
		Age:   21,
		Book:  []Book{},
	}

	Books = []Book{
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

//API Controllers for People
func GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []Person
	Db.Find(&people)

	json.NewEncoder(w).Encode(people)
}
func GetPerson(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	var books []Book

	Db.First(&person, params["id"])
	Db.Model(&person).Related(&books)
	person.Book = books

	json.NewEncoder(rw).Encode(person)
}

func CreatePerson(rw http.ResponseWriter, r *http.Request) {
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

func DeletePerson(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var person Person
	Db.First(&person, params["id"])

	Db.Delete(&person)

	json.NewEncoder(rw).Encode(&person)

}

// Cotroller for Books
func GetBooks(rw http.ResponseWriter, r *http.Request) {
	var books []Book

	Db.Find(&books)

	json.NewEncoder(rw).Encode(&books)
}

func GetBook(rw http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	var book Book
	Db.First(&book, params["id"])

	json.NewEncoder(rw).Encode(&book)

}

func CreateBook(rw http.ResponseWriter, r *http.Request) {
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

func DeleteBook(rw http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	var book Book

	Db.First(&book, params["id"])
	Db.Delete(&Personn)

	json.NewEncoder(rw).Encode(&book)

}

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/krzkowalczyk/GoSzkolenie/MVC/model"
)

func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	persons, iloscr := getAll()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	log.Printf("Ilość rekordów: %d\n", iloscr)
	json.NewEncoder(w).Encode(persons)
}

func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal("ID is not an integer !")
	}

	json.NewEncoder(w).Encode(getPersonByID(id))
}

func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	var person model.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		log.Fatal("Request body does not match model !")
	}
	createNewPerson(person.Name, person.Location)
}

func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatal("ID is not an integer !")
	}

	deletePersonByID(id)
}

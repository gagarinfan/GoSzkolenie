package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"fmt"
	"os"

	"github.com/gorilla/mux"

	_ "github.com/denisenkom/go-mssqldb"
)

var people []Person

var server = getEnv("SQL_SERVER", "localhost")
var user = getEnv("SQL_USER", "sa")
var port = getEnv("SQL_PORT", "1433")
var password = getEnv("SQL_PASSWORD", "P@ssw0rd")
var database = getEnv("SQL_DATABASE", "testdb")

var cnstr = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database)

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	_, err := getUserById(cnstr, params["id"])
	if err != nil {
		log.Fatal(err)
	}
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	iloscr, err := getAllUsers(cnstr)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Ilość rekordów: %d\n", iloscr)
	json.NewEncoder(w).Encode(people)
}
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	// params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	// person.ID = params["id"]
	// people = append(people, person)
	// json.NewEncoder(w).Encode(people)
	dodaj(cnstr, person.Name, person.Location)
}
func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// for index, item := range people {
	// 	if item.ID == params["id"] {
	// 		people = append(people[:index], people[index+1:]...)
	// 		break
	// 	}
	// 	json.NewEncoder(w).Encode(people)
	// }
	tid, _ := strconv.Atoi(params["id"])
	usun(cnstr, tid)
}

func main() {
	fmt.Println(cnstr)

	router := mux.NewRouter()
	// people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	// people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

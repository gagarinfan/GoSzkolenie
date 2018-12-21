package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gorilla/mux"
)

type Person struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Location string `json:"location,omitempty"`
}

var server = getEnv("SQL_SERVER", "localhost")
var user = getEnv("SQL_USER", "sa")
var port = getEnv("SQL_PORT", "1433")
var password = getEnv("SQL_PASSWORD", "P@ssword")
var database = getEnv("SQL_DATABASE", "testdb")

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

var people []Person
var cnstr = fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;", server, user, password, port, database)

func main() {

	fmt.Println(cnstr)
	var err error

	err = czytajdb(cnstr)
	if err != nil {
		log.Fatal("No nie mogę odczytać z bazy, bo: ", err)
	}
	router := mux.NewRouter()
	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}+{name}+{location}", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))

}

func czytajdb(cnstr string) error {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	sqlq := "select * from TestSchema.Employees"
	rows, errq := db.QueryContext(ctx, sqlq)
	if errq != nil {
		log.Fatal("Problem zapytania: ", errq.Error())
	}
	for rows.Next() {
		var id int
		var name, loc string

		errqw := rows.Scan(&id, &name, &loc)
		if errqw != nil {
			log.Fatal("Błąd w wyniku zapytania", errqw.Error())
			return errqw
		}
		fmt.Printf("id: %d, name: %s, location: %s\n", id, name, loc)

		people = append(people, Person{ID: strconv.Itoa(id), Name: name, Location: loc})
	}

	defer rows.Close()
	connclose(db)
	return errq
}

func dodaj(cnstr string, name string, location string, tableName string) (int64, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Println("Problem z db prepare: ", err.Error())
	}
	var zapins = "insert into TestSchema.EmployeesTMP (name, location) values (@name,@location); select @@identity;"

	skladnia, err := db.Prepare(zapins)
	if err != nil {
		log.Println("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()
	//var name, location string
	row := skladnia.QueryRowContext(ctx,
		//sql.Named("tableName", tableName),
		sql.Named("name", name),
		sql.Named("location", location))

	var noweid int64
	err = row.Scan(&noweid)
	if err != nil {
		log.Println("Błąd w wyniku zapytania: ", err.Error())
	}
	connclose(db)
	return noweid, err
}

func aktualizuj(cnstr string, id int, name string, location string) (int, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	var aktualizacja = "update Testschema.Employees set name=@name, location=@location where id=@id"

	skladnia, err := db.Prepare(aktualizacja)
	if err != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()

	newakt := fmt.Sprintf("update Testschema.Employees set name='@name', location='@location' where id=@id;")
	result, errctx := skladnia.ExecContext(ctx, newakt,
		sql.Named("name", name),
		sql.Named("location", location),
		sql.Named("id", id),
	)

	if errctx != nil {
		log.Fatal("Problem z execcontext: ", errctx.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatal("Panika ", err.Error())
	}

	log.Println("Gicior")
	connclose(db)
	return id, err

}

func usun(cnstr string, id int) (int, error) {
	var err error
	db, ctx, errdb := connopen(cnstr)
	if errdb != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}

	statement := fmt.Sprintf("delete from Testschema.Employees where id=@id")

	skladnia, err := db.Prepare(statement)
	if err != nil {
		log.Fatal("Problem z db prepare: ", err.Error())
	}
	defer skladnia.Close()

	result, errctx := skladnia.ExecContext(ctx, statement,
		sql.Named("id", id),
	)

	if errctx != nil {
		log.Fatal("Problem z execcontext: ", errctx.Error())
	}

	rows, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	if rows != 1 {
		log.Fatal("Panika ", err.Error())
	}

	log.Println("Gicior")
	connclose(db)
	return id, err

}

func connopen(cnstr string) (*sql.DB, context.Context, error) {
	var err error
	db, err := sql.Open("sqlserver", cnstr)

	if err != nil {
		log.Println("Błąd sterownika: ", err.Error())
	}
	ctx := context.Background()
	err = db.PingContext(ctx)

	if err != nil {
		log.Println("Problem puli połączenia: ", err.Error())
	}

	return db, ctx, err
}

func connclose(db *sql.DB) {
	defer db.Close()
}

func wypelnijTempTable() {
	//var sessionid string
	//sessionid = "tg1"
	nazwatb := "TestSchema.EmployeesTMP"

	for _, item := range people {
		{
			name := item.Name
			location := item.Location
			_, err := dodaj(cnstr, name, location, nazwatb)
			if err != nil {
				fmt.Println("Nie dodałem, bo: ", err.Error())
			}
		}
	}
}

//<import json>
func GetPersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}
func GetPeopleEndpoint(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}
func CreatePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	person.Name = params["name"]
	person.Location = params["location"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)

	//wywołanie dodania do temptable
	wypelnijTempTable()

	_, _ = dodaj(cnstr, person.Name, person.Location, "TestSchema.EmployeesTMP")
}
func DeletePersonEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

//</import json>

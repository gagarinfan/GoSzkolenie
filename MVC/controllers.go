package main

import (
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"github.com/krzkowalczyk/GoSzkolenie/MVC/model"
)

func getAll() (persons []model.Person, count int) {
	//	db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("mssql", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=testdb2;")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Find(&persons).Count(&count)

	log.Println("Number of records: ", count)
	for _, i := range persons {
		log.Println(i.ID, i.Name, i.Location)
	}
	return
}

func getPersonByID(id int) (person model.Person) {
	//	db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("mssql", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=testdb2;")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Where("ID = ?", id).First(&person)
	if person.ID == 0 {
		log.Println("Record does not exist")
		person = model.Person{}
	} else {
		log.Println("Selected: ", person.ID, person.Name, person.Location)

	}
	//debug

	return
}

func createNewPerson(name, location string) {
	//	db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("mssql", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=testdb2;")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Create(&model.Person{Name: name, Location: location})
}

func updatePersonByID(id int, name string, location string) (person model.Person) {
	//	db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("mssql", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=testdb2;")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	db.Where("ID = ?", id).First(&person)
	person.Name = name
	person.Location = location
	db.Save(&person)
	//debug
	log.Println("Updated: ", person.ID, person.Name, person.Location)

	return
}

func deletePersonByID(id int) {
	//	db, err := gorm.Open("sqlite3", "test.db")
	db, err := gorm.Open("mssql", "server=localhost;user id=sa;password=P@ssw0rd;port=1433;database=testdb2;")

	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	db.Delete(model.Person{}, "ID = ?", id)

	// db.Delete(&model.Person{ID: 1})
}

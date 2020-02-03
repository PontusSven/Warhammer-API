package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Table struct {
	Name string
}

type Unit struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Type           string `json:"type"`
	Movement       int    `json:"M"`
	WeaponsSkill   int    `json:"WS"`
	BallisticSkill int    `json:"BS"`
	Strength       int    `json:"S"`
	Toughness      int    `json:"T"`
	Wounds         int    `json:"W"`
	Attacks        int    `json:"A"`
	Leadership     int    `json:"Ld"`
	Initiate       int    `json:"I"`
	Points         int    `json:"points"`
}

var db *sql.DB
var err error

func welcome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func getArmies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var tables []Table
	result, err := db.Query("SHOW TABLES;")

	if err != nil {
		panic(err.Error)
	}

	defer result.Close()

	for result.Next() {
		var table Table
		err := result.Scan(&table.Name)
		if err != nil {
			panic(err.Error)
		}
		tables = append(tables, table)
	}
	json.NewEncoder(w).Encode(tables)

}

func getTyranids(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var units []Unit

	result, err := db.Query("SELECT * FROM tyranids")

	if err != nil {
		panic(err.Error)
	}

	defer result.Close()

	for result.Next() {
		var unit Unit
		err := result.Scan(&unit.ID, &unit.Name, &unit.Type, &unit.Movement, &unit.WeaponsSkill, &unit.BallisticSkill, &unit.Strength, &unit.Toughness, &unit.Wounds, &unit.Attacks, &unit.Leadership, &unit.Initiate, &unit.Points)
		if err != nil {
			panic(err.Error)
		}
		units = append(units, unit)
	}
	json.NewEncoder(w).Encode(units)
}

func getUnit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT * FROM tyranids WHERE name = ?", params["name"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var unit Unit
	for result.Next() {
		err := result.Scan(&unit.ID, &unit.Name, &unit.Type, &unit.Movement, &unit.WeaponsSkill, &unit.BallisticSkill, &unit.Strength, &unit.Toughness, &unit.Wounds, &unit.Attacks, &unit.Leadership, &unit.Initiate, &unit.Points)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(unit)
}

func main() {

	db, err = sql.Open("mysql", "root:jqjipotv12@tcp(127.0.0.1:3306)/Warhammer")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/tyranids", getTyranids).Methods("GET")
	router.HandleFunc("/tyranids/{name}", getUnit).Methods("GET")
	router.HandleFunc("/armies", getArmies).Methods("GET")

	http.ListenAndServe(":8001", router)
}

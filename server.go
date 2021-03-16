package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type PageVariables struct {
	Date string
	Time string
}

type Item struct {
	ID          string `json:"ID"`
	Item        string `json:"Item"`
	Description string `json:"Description"`
	Price       string `json:"Price"`
}

type allItems []Item

var Items = allItems{
	{
		ID:          "1",
		Item:        "Iphone",
		Description: "bla bla bla",
		Price:       "1000$",
	},
	{
		ID:          "3",
		Item:        "Drip",
		Description: "To be cool",
		Price:       "20000$",
	},
	{
		ID:          "4",
		Item:        "Laptop",
		Description: "Normal Laptop",
		Price:       "300$",
	},
	{
		ID:          "5",
		Item:        "Pc gaming",
		Description: "powerful pc",
		Price:       "2000$",
	},
}

func HomePage(w http.ResponseWriter, r *http.Request) {

	now := time.Now()              // find the time right now
	HomePageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("home.html") //parse the html file homepage.html
	if err != nil {                            // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func createItem(w http.ResponseWriter, r *http.Request) {
	var newItem Item
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "BRUH enter data with the event id, title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newItem)

	// Add the newly created event to the array of events
	Items = append(Items, newItem)

	// Return the 201 created status code
	w.WriteHeader(http.StatusCreated)
	// Return the newly created event
	json.NewEncoder(w).Encode(newItem)
}

func getOneItem(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	ItemID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for _, singleItem := range Items {
		if singleItem.ID == ItemID {
			json.NewEncoder(w).Encode(singleItem)
		}
	}
}

func getAllItems(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Items)
}

func updateItem(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	ItemID := mux.Vars(r)["id"]
	var updatedItem Item
	// Convert r.Body into a readable formart
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "BRUH enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &updatedItem)

	for i, singleItem := range Items {
		if singleItem.ID == ItemID {
			singleItem.Item = updatedItem.Item
			singleItem.Description = updatedItem.Description
			Items[i] = singleItem
			json.NewEncoder(w).Encode(singleItem)
		}
	}
}

func deleteItem(w http.ResponseWriter, r *http.Request) {
	// Get the ID from the url
	ItemID := mux.Vars(r)["id"]

	// Get the details from an existing event
	// Use the blank identifier to avoid creating a value that will not be used
	for i, singleItem := range Items {
		if singleItem.ID == ItemID {
			Items = append(Items[:i], Items[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", ItemID)
		}
	}
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", HomePage)
	router.HandleFunc("/add-item", createItem).Methods("POST")
	router.HandleFunc("/store", getAllItems).Methods("GET")
	router.HandleFunc("/item/{id}", getOneItem).Methods("GET")
	router.HandleFunc("/item/{id}", updateItem).Methods("PATCH")
	router.HandleFunc("/item/{id}", deleteItem).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}

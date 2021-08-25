package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

func main() {
	fmt.Println("Venusian Industries Rest API - vHydrogen - Mux Routers")
	vs = venus{
		MultiverseVersion: "Hydrogen",
		VenusianId: 1,
		Atmosphere: venusianAtmosphere{
			Pressure: 92.0,
		},
		Surface: venusianSurface{
			Temperature: 462.0,
		},
	}
	Drones = []drone{
		{Id: "1", Name: "HMS Iceslinger", Model: "Mining", CurrentFuel: 100},
		{Id: "2", Name: "Grinds-a-lot", Model: "Crushing", CurrentFuel: 100},
	}
	handleRequests()
}

func handleRequests() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	muxRouter.HandleFunc("/", homepage)
	muxRouter.HandleFunc("/venus", returnVenusStatus)
	muxRouter.HandleFunc("/drones", returnAllDrones).Methods("GET")
	muxRouter.HandleFunc("/drones", createNewDrone).Methods("POST")
	muxRouter.HandleFunc("/drones/{id}", returnSingleDrone).Methods("GET")
	muxRouter.HandleFunc("/drones/{id}", deleteDrone).Methods("DELETE")
	muxRouter.HandleFunc("/drones/{id}", updateDrone).Methods("PUT")
	log.Fatal(http.ListenAndServe(":50236", muxRouter))
}

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func prettyJSON(input interface{}) string {
	res, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(res)
}

var Drones []drone
type drone struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Model string `json:"model"`
	CurrentFuel uint32 `json:"currentFuel"`
}

func returnAllDrones(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllDrones")
	fmt.Fprint(w, prettyJSON(Drones))
}

func returnSingleDrone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	fmt.Println("Endpoint Hit: returnSingleDrone, key: " + key)

	for _, drone := range Drones {
		if drone.Id == key {
			fmt.Fprint(w, prettyJSON(drone))
		}
	}
}

var vs venus
func returnVenusStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnVenusStatus")
	// [dep] Serve json without indent
	// json.NewEncoder(w).Encode(vs)
	w.Header().Set("Content-Type", "application/json")
	// Serve json with indent
	fmt.Fprint(w, prettyJSON(vs))
}

type venus struct {
	MultiverseVersion string `json:"MultiverseVersion"`
	VenusianId uint8 `json:"VenusianId"`
	Atmosphere venusianAtmosphere `json:"Atmosphere"`
	Surface venusianSurface `json:"Surface"`
}

type venusianAtmosphere struct {
	Pressure float64 `json:"Pressure"`
}

type venusianSurface struct {
	Temperature float64 `json:"Temperature"`
}

func createNewDrone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: createNewDrone")
	// get the body of the POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newDrone drone
	json.Unmarshal(reqBody, &newDrone)
	index := getStructByFieldValue(Drones, "Id", newDrone.Id)
	if (index != -1) {
		fmt.Fprintf(w, "{\"error\": \"Could not execute CREATE as drone with id %s already exists\"}", newDrone.Id)
	} else {
		Drones = append(Drones, newDrone)
		fmt.Fprint(w, prettyJSON(Drones))
	}
}

func deleteDrone(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Endpoint Hit: deleteDrone, id: " + id)

	index := getStructByFieldValue(Drones, "Id", id)
	if (index == -1) {
		fmt.Fprintf(w, "{\"error\": \"Could not execute DELETE as drone with id %s does not exist\"}", id)
	} else {
		Drones = append(Drones[:index], Drones[index+1:]...)
		fmt.Fprint(w, prettyJSON(Drones))
	}
}

func updateDrone(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var newDrone drone
	json.Unmarshal(reqBody, &newDrone)
	
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Println("Endpoint Hit: updateDrone, id: " + id)

	index := getStructByFieldValue(Drones, "Id", id)
	if (index == -1) {
		fmt.Fprintf(w, "{\"error\": \"Could not execute UPDATE as drone with id %s does not exist\"}", id)
	} else {
		Drones[index] = newDrone
		fmt.Fprint(w, prettyJSON(Drones))
	}
}

func getStructByFieldValue(slice interface{} ,fieldName string,fieldValueToCheck interface {}) int {
	// Check for value of a given field in a slice of structs
	rangeOnMe := reflect.ValueOf(slice)
	for i := 0; i < rangeOnMe.Len(); i++ {
		s := rangeOnMe.Index(i)
		f := s.FieldByName(fieldName)
		if f.IsValid(){
			if f.Interface() == fieldValueToCheck {
				return i
			}
		}
	}
	return -1
}
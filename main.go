package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
		drone{Id: "1", Name: "HMS Iceslinger", Model: "Mining", CurrentFuel: 100},
		drone{Id: "2", Name: "Grinds-a-lot", Model: "Crushing", CurrentFuel: 100},
	}
	handleRequests()
}

func handleRequests() {
	muxRouter := mux.NewRouter().StrictSlash(true)
	muxRouter.HandleFunc("/", homepage)
	// use /hypixel/* prefix for testing
	muxRouter.HandleFunc("/venus", returnVenusStatus)
	muxRouter.HandleFunc("/drones", returnAllDrones)
	muxRouter.HandleFunc("/drones/{id}", returnSingleDrone)
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
	Pressure float64 `json: "Pressure"`
}

type venusianSurface struct {
	Temperature float64 `json: "Temperature"`
}
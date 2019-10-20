package main

// TODO: Add Open Flights as a better data source than Kiwi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"

	"strconv"

	yaml "gopkg.in/yaml.v2"
)

type Continent struct {
	Id   string
	Code string
	Name string
}

type Country struct {
	Id   string
	Name string
	Code string
	//Continent Continent
}

type City struct {
	Id        string
	Name      string
	Code      string
	Country   Country
	Continent Continent
}

type Airport struct {
	Id              string
	Name            string
	Icao            string
	Code            string
	Rank            int32
	PopularityScore float32 `json:"dst_popularity_score"`
	City            City
}

type LocationInfo interface {
	getTopAirports(client *http.Client, numCountries, numAirportsPerCountry int) map[string]Airport
}

// Kiwi provider for Location Info
// Using the Locations API https://docs.kiwi.com/locations/
type KiwiLocation struct {
	Countries []Country
	Aiports   []Airport
}

func (k KiwiLocation) getTopCountries(client *http.Client, numCountries int) map[string]Country {

	// https://api.skypicker.com/locations
	// ?type=dump&locale=en-US&location_types=country&limit=100&sort=rank&active_only=true
	req, _ := http.NewRequest("GET", "https://api.skypicker.com/locations", nil)

	q := req.URL.Query()
	q.Add("type", "dump")
	q.Add("locale", "en-US")
	q.Add("location_types", "country")
	q.Add("limit", strconv.Itoa(numCountries))
	q.Add("sort", "rank")
	//q.Add("partner", "picky")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		var countries []Country
		return createCountryMap(countries)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var countryMap struct {
		Locations []Country
	}
	json.Unmarshal(body, &countryMap)

	return createCountryMap(countryMap.Locations)
}

func (k KiwiLocation) getTopNAirports(client *http.Client, numAirportsPerCountry int) map[string]Airport {

	// https://api.skypicker.com/locations?type=dump&locale=en-US&location_types=airport
	// &limit=300&sort=rank&active_only=true

	req, _ := http.NewRequest("GET", "https://api.skypicker.com/locations", nil) // Top 300 Airports

	q := req.URL.Query()
	q.Add("type", "dump")
	q.Add("locale", "en-US")
	q.Add("location_types", "airport")
	q.Add("limit", strconv.Itoa(numAirportsPerCountry))
	q.Add("sort", "rank")
	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)

	if err != nil {
		var airports map[string]Airport
		return airports
	}

	defer resp.Body.Close()

	var airportObject struct {
		Locations []Airport
	}
	body, _ := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &airportObject)

	// var test struct {
	// 	Locations []interface{}
	// }

	// json.Unmarshal(body, &test)

	// fmt.Println(test)
	airportMap := make(map[string]Airport)
	for _, airport := range airportObject.Locations {
		airportMap[airport.Id] = airport
	}
	return airportMap
}

func (k KiwiLocation) getTopAirports(client *http.Client, numCountry, numAirportsPerCountry int) map[string]Airport {

	// https://api.skypicker.com/locations?type=dump&locale=en-US&location_types=airport
	// &limit=300&sort=rank&active_only=true

	countries := k.getTopCountries(client, numCountry)

	c := make(chan []Airport, runtime.NumCPU())

	for _, countryObj := range countries {
		country := countryObj.Id
		req, _ := http.NewRequest("GET", "https://api.skypicker.com/locations", nil) // Top 300 Airports

		q := req.URL.Query()
		q.Add("type", "subentity")
		q.Add("term", country)
		q.Add("active_only", "true")
		q.Add("locale", "en-US")
		q.Add("location_types", "airport")
		q.Add("limit", strconv.Itoa(numAirportsPerCountry))
		q.Add("sort", "rank")
		req.URL.RawQuery = q.Encode()

		go func(client *http.Client, req *http.Request) {
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()

				var airportMap struct {
					Locations []Airport
				}
				body, _ := ioutil.ReadAll(resp.Body)

				json.Unmarshal(body, &airportMap)

				c <- airportMap.Locations
			} else {
				fmt.Println("Error in fetching data for request: {method:", req.Method, ", URL:", req.URL, "}")
				var airport []Airport
				c <- airport
			}
		}(client, req)

	}

	airports := make(map[string]Airport, len(countries)*numAirportsPerCountry)
	for i := 0; i < len(countries); i++ {
		tempAirports := <-c
		for _, airport := range tempAirports {
			airports[airport.Id] = airport
		}
	}

	return airports
}

func createCountryMap(countries []Country) map[string]Country {
	countryMap := make(map[string]Country)
	for _, country := range countries {
		countryMap[country.Id] = country
	}
	return countryMap
}

func Locations(locFile *os.File, numCountries int, numAirportsPerCountry int, debug bool) {
	client := createHttpClient()
	var k KiwiLocation
	var airports map[string]Airport

	if numCountries == 0 {
		airports = k.getTopNAirports(client, numAirportsPerCountry)
	} else {
		airports = k.getTopAirports(client, numCountries, numAirportsPerCountry)
	}

	if debug {
		fmt.Println("Fetching Locations with arguments: {numLocations:", numCountries, "}")
		fmt.Println("Got airports:\n", airports)
	}

	if len(airports) > 0 {
		countriesYaml, _ := yaml.Marshal(airports)
		if debug {
			fmt.Println(string(countriesYaml))
		}

		fmt.Println("Writing", len(airports), "airports to file", locFile.Name())
		locFile.Write(countriesYaml)
	}
}

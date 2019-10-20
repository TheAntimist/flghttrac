package main

import (
	"fmt"
	"os"

	"github.com/akamensky/argparse"
)

func main() {

	parser := argparse.NewParser("flghttrac", "Flight Tracker, which based on the home airport can send a notifications for when there are low fares on specific routes")

	debug := parser.Flag("d", "debug", &argparse.Options{Required: false, Help: "Enables debug mode."})

	locations := parser.NewCommand("locations", "Exports top locations to specific file")

	numCountries := locations.Int("c", "countries", &argparse.Options{Required: false, Default: 100, Help: "Number of countries to use for airports output"})

	numAirportsPerCountry := locations.Int("a", "aiports", &argparse.Options{Required: false, Default: 2, Help: "Number of airports per country to fetch"})

	locFile := locations.File("l", "loc-file", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644, &argparse.Options{Required: false, Default: "locations.yml", Help: "File where to export to"})

	startCmd := parser.NewCommand("start", "Start tracking for new changes")

	configFile := startCmd.File("c", "config-file", os.O_RDONLY, 0600, &argparse.Options{Required: false, Default: "config.yml", Help: "Config file to read the enabled notifications and data from"})

	err := parser.Parse(os.Args)

	if err != nil {
		fmt.Println(parser.Usage(err))
		os.Exit(1)
	}

	switch true {
	case locations.Happened():
		Locations(locFile, *numCountries, *numAirportsPerCountry, *debug)
		break
	case startCmd.Happened():
		Start(configFile)
		break
	}

	defer configFile.Close()
	defer locFile.Close()
}

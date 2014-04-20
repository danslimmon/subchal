package main

import (
    "fmt"
    "log"
    "strings"
    "flag"
    "github.com/danslimmon/subchal/lib"
)

func main() {
    dataDir := flag.String("datadir", "", "directory containing the GTFS data")
    flag.Parse()

    stopsByID, stationsByID, err := subchal.LoadStops(strings.Join([]string{*dataDir, "stops.txt"}, "/"))
    if err != nil {
        log.Println("Error loading stop data: ", err)
        return
    }

    log.Println(fmt.Sprintf("Loaded %d stations and %d stops.", len(stationsByID), len(stopsByID)))
}

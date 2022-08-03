package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	dataPath              string = "data.json"
	layout                string = "15:04:05"
	priceCriteria         string = "price"
	arrivalTimeCriteria   string = "arrival-time"
	departureTimeCriteria string = "departure-time"
	numOfReturnTrains     int    = 3
	minIDNumber           int    = 1
	minLen                int    = 1
)

var (
	errCriteria        = errors.New("unsupported criteria")
	errEmptyDepStation = errors.New("empty departure station")
	errEmptyArrStation = errors.New("empty arrival station")
	errBadDepStation   = errors.New("bad departure station input")
	errBadArrStation   = errors.New("bad arrival station input")
	criteriaOfSort     = map[string]string{
		"price":          "price",
		"arrival-time":   "arrival-time",
		"departure-time": "departure-time",
	}
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

func (t *Train) UnmarshalJSON(b []byte) error {

	var alias struct {
		TrainID            int        `json:"trainId"`
		DepartureStationID int        `json:"departureStationId"`
		ArrivalStationID   int        `json:"arrivalStationId"`
		Price              float32    `json:"price"`
		ArrivalTime        CustomTime `json:"arrivalTime"`
		DepartureTime      CustomTime `json:"departureTime"`
	}

	err := json.Unmarshal(b, &alias)
	if err != nil {
		return err
	}

	t.TrainID = alias.TrainID
	t.DepartureStationID = alias.DepartureStationID
	t.ArrivalStationID = alias.ArrivalStationID
	t.Price = alias.Price
	t.ArrivalTime = time.Time(alias.ArrivalTime)
	t.DepartureTime = time.Time(alias.DepartureTime)
	return nil
}

func (t *Train) printTrain() {
	fmt.Printf("%+v\n", t)
}

type CustomTime time.Time

func (c *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `\"`)
	if s == "null" || s == "" {
		return nil
	}
	t, err := time.Parse(layout, s)
	if err != nil {
		return err
	}
	*c = CustomTime(t)
	return nil
}

type Query struct {
	DepartureStationID string
	ArrivalStationID   string
	Criteria           string
}

func (q *Query) readUserParamsFromTerminal() {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter departure station ID: ")
	q.DepartureStationID, _ = reader.ReadString('\n')
	q.DepartureStationID = q.DepartureStationID[0 : len(q.DepartureStationID)-1]

	fmt.Print("Enter arrival station ID: ")
	q.ArrivalStationID, _ = reader.ReadString('\n')
	q.ArrivalStationID = q.ArrivalStationID[0 : len(q.ArrivalStationID)-1]

	fmt.Print("Enter sorting criteria: ")
	q.Criteria, _ = reader.ReadString('\n')
	q.Criteria = q.Criteria[0 : len(q.Criteria)-1]
}

func main() {

	//query data from user
	query := new(Query)
	query.readUserParamsFromTerminal()

	//find trains by query params
	result, err := FindTrains(query.DepartureStationID, query.ArrivalStationID, query.Criteria)

	//handle error
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	//print result
	printFindingResult(result)
}

func readDataFromJSON(path string) Trains {

	d := make(Trains, 0)

	jsonFile, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &d)
	return d
}

func printFindingResult(t Trains) {
	if len(t) >= minLen {
		for _, v := range t {
			v.printTrain()
		}
	} else {
		fmt.Println("Can't find trains on request. Please, try again.")
	}

}

func sortTrainsByAscending(trains Trains, criteria string) {
	fmt.Println("Sorting by", criteria)
	switch criteria {
	case priceCriteria:
		sort.SliceStable(trains, func(i, j int) bool {
			return trains[i].Price < trains[j].Price
		})
	case departureTimeCriteria:
		sort.SliceStable(trains, func(i, j int) bool {
			return trains[i].DepartureTime.Before(trains[j].DepartureTime)
		})
	case arrivalTimeCriteria:
		sort.SliceStable(trains, func(i, j int) bool {
			return trains[i].ArrivalTime.Before(trains[j].ArrivalTime)
		})
	}
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {

	if len(departureStation) < minLen {
		return nil, errEmptyDepStation
	}

	if len(arrivalStation) < minLen {
		return nil, errEmptyArrStation
	}

	departureStationID, err := strconv.Atoi(departureStation)
	if err != nil || departureStationID < minIDNumber {
		return nil, errBadDepStation
	}

	arrivalStationID, err := strconv.Atoi(arrivalStation)
	if err != nil || arrivalStationID < minIDNumber {
		return nil, errBadArrStation
	}

	if _, ok := criteriaOfSort[criteria]; !ok {
		return nil, errCriteria
	}

	data := readDataFromJSON(dataPath)
	trains := make(Trains, 0)

	for _, v := range data {
		if v.DepartureStationID == departureStationID && v.ArrivalStationID == arrivalStationID {
			trains = append(trains, v)
		}
	}

	if len(trains) < minLen {
		return nil, nil
	}
	sortTrainsByAscending(trains, criteria)

	if len(trains) < numOfReturnTrains {
		return trains, nil
	}
	return trains[:numOfReturnTrains], nil
}

package main

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFindTrainsSuccess(t *testing.T) {
	assert := assert.New(t)
	testsOK := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		exp        Trains
		expErr     error
	}{
		"successful_price": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "price",
			exp: Trains{
				{TrainID: 1177, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 164.65, ArrivalTime: time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 36, 0, 0, time.UTC)},
				{TrainID: 1178, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 164.65, ArrivalTime: time.Date(0, time.January, 1, 10, 25, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 36, 0, 0, time.UTC)},
				{TrainID: 1141, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 176.77, ArrivalTime: time.Date(0, time.January, 1, 12, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 16, 48, 0, 0, time.UTC)},
			},
			expErr: nil,
		},
		"successful_arrival": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "arrival-time",
			exp: Trains{
				{TrainID: 978, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 258.53, ArrivalTime: time.Date(0, time.January, 1, 4, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 10, 0, 0, time.UTC)},
				{TrainID: 1316, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 209.73, ArrivalTime: time.Date(0, time.January, 1, 5, 55, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 52, 0, 0, time.UTC)},
				{TrainID: 2201, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 280, ArrivalTime: time.Date(0, time.January, 1, 6, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 14, 55, 0, 0, time.UTC)},
			},
			expErr: nil,
		},
		"successful_departure": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "departure-time",
			exp: Trains{
				{TrainID: 1386, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 220.49, ArrivalTime: time.Date(0, time.January, 1, 8, 30, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 3, 0, 0, time.UTC)},
				{TrainID: 978, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 258.53, ArrivalTime: time.Date(0, time.January, 1, 4, 15, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 10, 0, 0, time.UTC)},
				{TrainID: 1316, DepartureStationID: 1902, ArrivalStationID: 1929, Price: 209.73, ArrivalTime: time.Date(0, time.January, 1, 5, 55, 0, 0, time.UTC), DepartureTime: time.Date(0, time.January, 1, 13, 52, 0, 0, time.UTC)},
			},
			expErr: nil,
		},
	}

	for name, tc := range testsOK {
		t.Run(name, func(t *testing.T) {
			got, gotErr := FindTrains(tc.depStation, tc.arrStation, tc.criteria)
			if assert.NoError(gotErr) {
				assert.Equal(tc.exp, got)
			}
		})
	}
}
func TestFindTrainsNilNil(t *testing.T) {
	assert := assert.New(t)
	testsNIL := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		exp        Trains
		expErr     error
	}{
		"wrong_depStation": {
			depStation: "12",
			arrStation: "1929",
			criteria:   "price",
			exp:        nil,
			expErr:     nil,
		},
		"wrong_arrStation": {
			depStation: "1902",
			arrStation: "11",
			criteria:   "price",
			exp:        nil,
			expErr:     nil,
		},
	}

	for name, tc := range testsNIL {
		t.Run(name, func(t *testing.T) {
			got, gotErr := FindTrains(tc.depStation, tc.arrStation, tc.criteria)
			assert.Nil(got)
			assert.Nil(gotErr)
		})
	}
}
func TestFindTrainsErrors(t *testing.T) {
	assert := assert.New(t)
	testsNotOK := map[string]struct {
		arrStation string
		depStation string
		criteria   string
		exp        Trains
		expErr     error
	}{
		"wrong_criteria": {
			depStation: "1902",
			arrStation: "1929",
			criteria:   "awef",
			exp:        nil,
			expErr:     errors.New("unsupported criteria"),
		},
		"absent_depStationId": {
			depStation: "",
			arrStation: "1929",
			criteria:   "departure",
			exp:        nil,
			expErr:     errors.New("empty departure station"),
		},
		"absent_arrStation": {
			depStation: "1902",
			arrStation: "",
			criteria:   "departure",
			exp:        nil,
			expErr:     errors.New("empty arrival station"),
		},

		"bad_arrStation_input": {
			depStation: "1902",
			arrStation: "serg",
			criteria:   "price",
			exp:        nil,
			expErr:     errors.New("bad arrival station input"),
		},
		"bad_depStation_input": {
			depStation: "serg",
			arrStation: "1922",
			criteria:   "price",
			exp:        nil,
			expErr:     errors.New("bad departure station input"),
		},
	}

	for name, tc := range testsNotOK {
		t.Run(name, func(t *testing.T) {
			got, gotErr := FindTrains(tc.depStation, tc.arrStation, tc.criteria)
			if assert.Error(gotErr) && assert.Nil(got) {
				assert.EqualError(gotErr, tc.expErr.Error())
			}
		})
	}
}
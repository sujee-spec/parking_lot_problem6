package parkinglot

import (
	"errors"
	"math"
)

type ParkingType string

const (
	SimpleParking ParkingType = "simpleParking"
	EvenParking   ParkingType = "evenParking"
)

type Attendant struct {
	parkingPlan        ParkingType
	Parkinglot         []*ParkingLot
	parkingsFull       []bool
	slotsOccupiedCount []uint
}

func NewAttendantV2(parkingPlan ParkingType, parkingLots ...*ParkingLot) (*Attendant, error) {
	parkinglotSlice := []*ParkingLot{}
	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("attendant cannot have nil parkinglot")
		}
	}

	parkinglotSlice = append(parkinglotSlice, parkingLots...)

	statuses := make([]bool, len(parkingLots))
	slotsOccupiedCount := make([]uint, len(parkingLots))
	attendant := Attendant{
		parkingPlan:        parkingPlan,
		Parkinglot:         parkingLots,
		parkingsFull:       statuses,
		slotsOccupiedCount: slotsOccupiedCount,
	}

	for _, parkinglot := range parkinglotSlice {
		parkinglot.OnFull(&attendant)
	}

	return &attendant, nil
}

func NewAttendant(parkingLots ...*ParkingLot) (*Attendant, error) {
	parkinglotSlice := []*ParkingLot{}
	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("attendant cannot have nil parkinglot")
		}
	}

	parkinglotSlice = append(parkinglotSlice, parkingLots...)

	statuses := make([]bool, len(parkingLots))
	slotsOccupiedCount := make([]uint, len(parkingLots))
	attendant := Attendant{
		Parkinglot:         parkingLots,
		parkingsFull:       statuses,
		slotsOccupiedCount: slotsOccupiedCount,
	}

	for _, parkinglot := range parkinglotSlice {
		parkinglot.OnFull(&attendant)
	}

	return &attendant, nil
}

func (a *Attendant) Park(car *Car) error {
	if car == nil {
		return errors.New("car cannot be nil")
	}

	if a.checkIsCarParked(car) {
		return errors.New("attendant: car already parked")
	}

	index := a.findAvailableParkinglotIndex(a.parkingPlan)
	if index < 0 {
		return errors.New("parking lot is full, attendant cannot park the car")
	}

	err := a.Parkinglot[index].park(car)
	if err != nil {
		return err
	}
	a.slotsOccupiedCount[index]++

	return nil
}

func (a *Attendant) findAvailableParkinglotIndex(parkingPlan ParkingType) int {

	if parkingPlan == EvenParking {
		count := math.MaxInt64
		parkinglotIndex := -1
		for i, parkingCount := range a.slotsOccupiedCount {
			if parkingCount < uint(count) {
				count = int(parkingCount)
				parkinglotIndex = i
			}
		}
		if parkinglotIndex == -1 {
			return -1
		}
		return parkinglotIndex
	}

	//when simple parking plan or no parkingPlan
	for i := range a.Parkinglot {
		if a.parkingsFull[i] {
			continue
		}
		return i
	}
	return -1
}

// TODO: refactor
func (a *Attendant) UnPark(car *Car) error {
	if car == nil {
		return errors.New("attendant/unpark: car cannot be nil")
	}

	if !a.checkIsCarParked(car) {
		return errors.New("attendant/unpark: car is not parked")
	}

	var err error
	for i, parkinglot := range a.Parkinglot {
		if !parkinglot.isParked(car) {
			continue
		}
		err := parkinglot.unPark(car)
		if err != nil {
			return err
		}
		a.slotsOccupiedCount[i]--
		a.parkingsFull[i] = false
		return nil
	}

	return err
}

func (a *Attendant) receiveFull(i int) {
	a.parkingsFull[i] = true
}

func (a *Attendant) checkIsCarParked(car *Car) bool {
	for _, parkinglot := range a.Parkinglot {
		if parkinglot.isParked(car) {
			return true
		}
	}
	return false
}

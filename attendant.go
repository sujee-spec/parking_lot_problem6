package parkinglot

import (
	"errors"
	"math"
)

type ParkingType string

const (
	ParkInFirstAvailableParkinglot    ParkingType = "parkOnFirstAvailable"
	ParkInParkinglotWithLeastVehicles ParkingType = "parkinglotWithLeastVehicleParking"
	MostOccupiedParking               ParkingType = "parkonMostOccupiedParking"
)

type getParkinglotFn func(*Attendant) (*ParkingLot, error)

type Attendant struct {
	parkingPlanFn     getParkinglotFn
	parkinglots       []*ParkingLot
	parkingFullStatus []bool
}

func NewAttendantV2(parkingPlan ParkingType, parkingLots ...*ParkingLot) (*Attendant, error) {

	attendant, err := NewAttendant(parkingLots...)
	if err != nil {
		return nil, err
	}

	parkingPlanFn := decideParkingPlanFn(parkingPlan)

	attendant.parkingPlanFn = parkingPlanFn

	return attendant, nil
}

func decideParkingPlanFn(parkingPlan ParkingType) getParkinglotFn {

	switch parkingPlan {
	case ParkInParkinglotWithLeastVehicles:
		return findParkinglotWithLeastVehicle
	case ParkInFirstAvailableParkinglot:
		return findFirstAvailableParkingLot
	default:
		return findParkinglotWithMostVehicles
	}
}

func NewAttendant(parkingLots ...*ParkingLot) (*Attendant, error) {

	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("newattendant cannot have nil parkinglot")
		}
	}

	statuses := make([]bool, len(parkingLots))
	attendant := Attendant{
		parkingPlanFn:     findFirstAvailableParkingLot,
		parkinglots:       parkingLots,
		parkingFullStatus: statuses,
	}

	for _, parkinglot := range parkingLots {
		parkinglot.AddParkingFullReceiver(&attendant)
	}

	for _, parkinglot := range parkingLots {
		parkinglot.AddParkingAvailableReceiver(&attendant)
	}

	return &attendant, nil
}

func (a *Attendant) Park(car *Car) error {
	if car == nil {
		return errors.New("attendant cannot park nil car")
	}

	if a.isParked(car) {
		return errors.New("attendant cannot park already parked car")
	}

	parkinglot, err := a.parkingPlanFn(a)
	if err != nil {
		return err
	}

	return parkinglot.park(car)

}

func findParkinglotWithLeastVehicle(a *Attendant) (*ParkingLot, error) {
	minOccupied := math.MaxInt64
	selectedLotIndex := -1

	for i, parkingStatus := range a.parkingFullStatus {
		if parkingStatus {
			continue
		}
		occupiedCount := a.parkinglots[i].countOccupiedSlots()
		if occupiedCount < minOccupied {
			minOccupied = occupiedCount
			selectedLotIndex = i
		}
	}
	if selectedLotIndex == -1 {
		return nil, errors.New("no parkinglot with least vehicle available")
	}

	return a.parkinglots[selectedLotIndex], nil
}

func findFirstAvailableParkingLot(a *Attendant) (*ParkingLot, error) {

	for i := range a.parkingFullStatus {
		if !a.parkingFullStatus[i] {
			return a.parkinglots[i], nil
		}
	}
	return nil, errors.New("all parkinglots are full")

}

func findParkinglotWithMostVehicles(a *Attendant) (*ParkingLot, error) {
	maxOccupied := math.MinInt64
	selectedLotIndex := -1

	for i, parkingStatus := range a.parkingFullStatus {
		if parkingStatus {
			continue
		}
		occupiedCount := a.parkinglots[i].countOccupiedSlots()
		if occupiedCount > maxOccupied {
			maxOccupied = occupiedCount
			selectedLotIndex = i
		}
	}

	if selectedLotIndex == -1 {
		return nil, errors.New("no parkinglot with most vehicle available")
	}

	return a.parkinglots[selectedLotIndex], nil
}

func (a *Attendant) UnPark(car *Car) error {
	if car == nil {
		return errors.New("attendant cannot unpark nil car")
	}

	for _, parkinglot := range a.parkinglots {
		if !parkinglot.isParked(car) {
			continue
		}

		return parkinglot.unPark(car)
	}

	return errors.New("attendant cannot unpark, car not found in parlinglots")

}

func (a *Attendant) receiveParkingFullStatus(i int) {
	a.parkingFullStatus[i] = true
}

func (a *Attendant) receiveParkingAvailableStatus(i int) {
	a.parkingFullStatus[i] = false
}

func (a *Attendant) isParked(car *Car) bool {
	for _, parkinglot := range a.parkinglots {
		if parkinglot.isParked(car) {
			return true
		}
	}
	return false
}

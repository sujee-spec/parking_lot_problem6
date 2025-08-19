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
	parkingPlan  ParkingType
	Parkinglot   []*ParkingLot
	parkingsFull []bool
}

func NewAttendantV2(parkingPlan ParkingType, parkingLots ...*ParkingLot) (*Attendant, error) {

	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("attendant cannot have nil parkinglot")
		}
	}

	statuses := make([]bool, len(parkingLots))
	attendant := Attendant{
		parkingPlan:  parkingPlan,
		Parkinglot:   parkingLots,
		parkingsFull: statuses,
	}

	for _, parkinglot := range parkingLots {
		parkinglot.OnFull(&attendant)
	}

	return &attendant, nil
}

func NewAttendant(parkingLots ...*ParkingLot) (*Attendant, error) {

	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("attendant cannot have nil parkinglot")
		}
	}

	statuses := make([]bool, len(parkingLots))
	attendant := Attendant{
		Parkinglot:   parkingLots,
		parkingsFull: statuses,
	}

	for _, parkinglot := range parkingLots {
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

	parkinglot := a.findAvailableParkinglot(a.parkingPlan)
	if parkinglot == nil {
		return errors.New("parking lot is full, attendant cannot park the car")
	}

	return parkinglot.park(car)

}

func (a *Attendant) findAvailableParkinglot(parkingPlan ParkingType) *ParkingLot {
	switch parkingPlan {
	case EvenParking:
		return a.findEvenlyDistributedParkingLot()
	default:
		return a.findFirstAvailableParkingLot()
	}
}

func (a *Attendant) findEvenlyDistributedParkingLot() *ParkingLot {
	minOccupied := math.MaxInt64
	var selectedLot *ParkingLot

	for _, lot := range a.Parkinglot {
		occupiedCount := countOccupiedSlots(lot)
		if occupiedCount < minOccupied {
			minOccupied = occupiedCount
			selectedLot = lot
		}
	}

	return selectedLot
}

func (a *Attendant) findFirstAvailableParkingLot() *ParkingLot {
	for i, lot := range a.Parkinglot {
		if !a.parkingsFull[i] {
			return lot
		}
	}
	return nil
}

func countOccupiedSlots(lot *ParkingLot) int {
	count := 0
	for _, slot := range lot.slots {
		if slot.occupied {
			count++
		}
	}
	return count
}

// TODO: refactor
func (a *Attendant) UnPark(car *Car) error {
	if car == nil {
		return errors.New("attendant/unpark: car cannot be nil")
	}

	if !a.checkIsCarParked(car) {
		return errors.New("attendant/unpark: car is not parked")
	}

	for i, parkinglot := range a.Parkinglot {
		if !parkinglot.isParked(car) {
			continue
		}

		err := parkinglot.unPark(car)
		if err != nil {
			return err
		}

		a.parkingsFull[i] = false
	}

	return nil
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

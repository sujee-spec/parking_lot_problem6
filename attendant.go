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
	parkinglotSlice := []*ParkingLot{}
	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("attendant cannot have nil parkinglot")
		}
	}

	parkinglotSlice = append(parkinglotSlice, parkingLots...)

	statuses := make([]bool, len(parkingLots))
	attendant := Attendant{
		parkingPlan:  parkingPlan,
		Parkinglot:   parkingLots,
		parkingsFull: statuses,
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
	attendant := Attendant{
		Parkinglot:   parkingLots,
		parkingsFull: statuses,
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

	parkinglot, index := a.findAvailableParkinglot(a.parkingPlan)
	if parkinglot == nil || index < 0 {
		return errors.New("parking lot is full, attendant cannot park the car")
	}

	err := parkinglot.park(car)
	if err != nil {
		return err
	}

	return nil
}

func (a *Attendant) findAvailableParkinglot(parkingPlan ParkingType) (*ParkingLot, int) {

	if parkingPlan == EvenParking {
		count := math.MaxInt64
		parkinglotIndex := -1
		for i, plot := range a.Parkinglot {
			internalCount := 0
			for _, p := range plot.slots {
				if p.occupied {
					internalCount++
				}
			}
			if internalCount < count {
				count = internalCount
				parkinglotIndex = i
			}
		}
		if parkinglotIndex == -1 {
			return nil, -1
		}
		return a.Parkinglot[parkinglotIndex], parkinglotIndex
	}

	//when simple parking plan or no parkingPlan
	for i, p := range a.Parkinglot {
		if a.parkingsFull[i] {
			continue
		}
		return p, i
	}
	return nil, -1
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

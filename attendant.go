package parkinglot

import (
	"errors"
	"math"
)

type ParkingType string

const (
	FirstAvailableParking   ParkingType = "parkOnFirstAvailable"
	EvenDistributionParking ParkingType = "evenDistributionParking"
	MostOccupiedParking     ParkingType = "mostOccupiedParking"
)

type parkinglotMethod func(*Attendant) *ParkingLot

type Attendant struct {
	parkinglotMethodStyle parkinglotMethod
	Parkinglots           []*ParkingLot
	parkingsFull          []bool
}

func NewAttendantV2(parkingPlan ParkingType, parkingLots ...*ParkingLot) (*Attendant, error) {

	attendant, err := NewAttendant(parkingLots...)
	if err != nil {
		return nil, err
	}

	var decideParkinglotMethod parkinglotMethod
	if parkingPlan == EvenDistributionParking {
		decideParkinglotMethod = findEvenlyDistributedParkingLot
	} else if parkingPlan == FirstAvailableParking {
		decideParkinglotMethod = findFirstAvailableParkingLot
	} else {
		decideParkinglotMethod = findMostFilledParkingLot
	}

	attendant.parkinglotMethodStyle = decideParkinglotMethod

	return attendant, nil
}

func NewAttendant(parkingLots ...*ParkingLot) (*Attendant, error) {

	for _, parkingLot := range parkingLots {
		if parkingLot == nil {
			return nil, errors.New("newattendant cannot have nil parkinglot")
		}
	}

	statuses := make([]bool, len(parkingLots))
	attendant := Attendant{
		parkinglotMethodStyle: findFirstAvailableParkingLot,
		Parkinglots:           parkingLots,
		parkingsFull:          statuses,
	}

	for _, parkinglot := range parkingLots {
		parkinglot.AddParkingFullReceivers(&attendant)
	}

	for _, parkinglot := range parkingLots {
		parkinglot.AddParkingAvailableReceivers(&attendant)
	}

	return &attendant, nil
}

func (a *Attendant) Park(car *Car) error {
	if car == nil {
		return errors.New("attendant cannot park nil car")
	}

	if a.checkIsCarParked(car) {
		return errors.New("attendant cannot park already parked car")
	}

	parkinglot := a.parkinglotMethodStyle(a)
	if parkinglot == nil {
		return errors.New("parking lot is full, attendant cannot park the car")
	}

	return parkinglot.park(car)

}

func findEvenlyDistributedParkingLot(a *Attendant) *ParkingLot {
	minOccupied := math.MaxInt64
	var selectedLot *ParkingLot

	for _, lot := range a.Parkinglots {
		occupiedCount := countOccupiedSlots(lot)
		if occupiedCount < minOccupied {
			minOccupied = occupiedCount
			selectedLot = lot
		}
	}

	return selectedLot
}

func findFirstAvailableParkingLot(a *Attendant) *ParkingLot {
	for i, lot := range a.Parkinglots {
		if !a.parkingsFull[i] {
			return lot
		}
	}
	return nil
}

func findMostFilledParkingLot(a *Attendant) *ParkingLot {
	maxOccupied := math.MinInt64
	var selectedLot *ParkingLot

	for _, lot := range a.Parkinglots {
		occupiedCount := countOccupiedSlots(lot)
		if occupiedCount > maxOccupied {
			maxOccupied = occupiedCount
			selectedLot = lot
		}
	}

	return selectedLot
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
		return errors.New("attendant cannot unpark nil car")
	}

	if !a.checkIsCarParked(car) {
		return errors.New("attendant cannot unpark, car not found in parlinglots")
	}

	for _, parkinglot := range a.Parkinglots {
		if !parkinglot.isParked(car) {
			continue
		}

		return parkinglot.unPark(car)
	}

	return nil
}

func (a *Attendant) receiveFull(i int) {
	a.parkingsFull[i] = true
}

func (a *Attendant) receiveAvailable(i int) {
	a.parkingsFull[i] = false
}

func (a *Attendant) checkIsCarParked(car *Car) bool {
	for _, parkinglot := range a.Parkinglots {
		if parkinglot.isParked(car) {
			return true
		}
	}
	return false
}

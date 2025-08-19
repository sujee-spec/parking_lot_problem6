package parkinglot

import "errors"

type Attendant struct {
	Parkinglot      []*ParkingLot
	parkingStatuses []bool
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
		Parkinglot:      parkingLots,
		parkingStatuses: statuses,
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

	for i, p := range a.Parkinglot {
		if a.parkingStatuses[i] {
			continue
		}
		return  p.park(car)
	}

	return errors.New("parking lot is full, attendant cannot park the car")
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
		a.parkingStatuses[i] = false
		return nil
	}

	return err
}

func (a *Attendant) receiveFull(i int) {
	a.parkingStatuses[i] = true
}

func (a *Attendant) checkIsCarParked(car *Car) bool {
	for _, parkinglot := range a.Parkinglot {
		if parkinglot.isParked(car) {
			return true
		}
	}
	return false
}

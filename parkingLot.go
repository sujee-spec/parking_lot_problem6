package parkinglot

import (
	"errors"
)

type ParkingFullReceiver interface {
	receiveFull(int)
}

type ParkingAvailableReceiver interface {
	receiveAvailable(int)
}

type slot struct {
	id       int
	car      *Car
	occupied bool
}

func (s *slot) isEmpty() bool {
	return !s.occupied
}
func (s *slot) isNotEmpty() bool {
	return s.occupied
}

func (s *slot) occupy(c *Car) {
	s.car = c
	s.occupied = true
}

func (s *slot) free() {
	s.car = nil
	s.occupied = false
}

type ParkingLot struct {
	id                   int
	capacity             int
	slots                []slot
	fullSubscribers      []ParkingFullReceiver
	availableSubscribers []ParkingAvailableReceiver
}

func (p *ParkingLot) AddParkingAvailableReceivers(parkingAvailableReceiver ParkingAvailableReceiver) {
	p.availableSubscribers = append(p.availableSubscribers, parkingAvailableReceiver)
}

func (p *ParkingLot) AddParkingFullReceivers(r ParkingFullReceiver) {
	p.fullSubscribers = append(p.fullSubscribers, r)
}

type Car struct {
	numberPlate string
}

func (car1 *Car) isEqual(car2 *Car) bool {
	return car1.numberPlate == car2.numberPlate
}

func NewParkingLot(capacity int) (*ParkingLot, error) {
	if capacity < 1 {
		return nil, errors.New("cannot create parking lot with capacity less than 1")
	}
	lots := make([]slot, 0, capacity)
	for i := 0; i < capacity; i++ {
		newLot := slot{
			id: i + 1,
		}
		lots = append(lots, newLot)
	}

	return &ParkingLot{slots: lots, capacity: capacity}, nil
}

func (p *ParkingLot) park(c *Car) error {
	if c == nil {
		return errors.New("parkinglot cannot park nil car")
	}
	if p.isParked(c) {
		return errors.New("parkinglot cannot park already parked car")
	}
	for i := 0; i < p.capacity; i++ {
		if p.slots[i].isNotEmpty() {
			continue
		}
		p.slots[i].occupy(c)

		if p.isFullyFilled() {

			p.notifyFullReceivers()

		}

		return nil

	}
	return errors.New("parkinglot full, cannot park more cars")

}

func (p *ParkingLot) isFullyFilled() bool {
	for _, slot := range p.slots {
		if !slot.occupied {
			return false
		}
	}
	return true
}

func (p *ParkingLot) unPark(car *Car) error {
	if car == nil {
		return errors.New("parkinglot cannot unpark nil car")
	}

	if !p.isParked(car) {
		return errors.New("parkinglot cannot unpark, car not found in parkinglot")
	}

	for i := 0; i < p.capacity; i++ {
		if p.slots[i].isEmpty() {
			continue
		}

		if !p.slots[i].car.isEqual(car) {
			continue
		}

		isParkingFull := p.isFullyFilled()
		p.slots[i].free()

		if isParkingFull {
			p.notifyAvailableRecievers()
		}

	}
	return nil
}

func (p *ParkingLot) notifyAvailableRecievers() {
	if p.availableSubscribers == nil {
		return
	}
	for _, r := range p.availableSubscribers {
		r.receiveAvailable(p.id)
	}
}

func (p *ParkingLot) isParked(car *Car) bool {
	for _, slot := range p.slots {
		if slot.isEmpty() {
			continue
		}
		if car.isEqual(slot.car) {
			return true
		}
	}
	return false
}

func (p *ParkingLot) notifyFullReceivers() {
	for _, r := range p.fullSubscribers {
		r.receiveFull(p.id)
	}

}

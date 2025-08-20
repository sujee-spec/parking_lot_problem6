package parkinglot

import (
	"testing"
)

func TestCreateLot(t *testing.T) {
	firstLot := &slot{
		id:       1,
		car:      &Car{"KJ-09-AK-123"},
		occupied: false,
	}

	if firstLot.id == 0 {
		t.Error("Struct Lot cannot be created with Id 0")
	}
}

func TestCreateParkingLot(t *testing.T) {
	_, err := NewParkingLot(10)
	if err != nil {
		t.Error("Parking lot is not created")
	}
}

func TestCreateEmptyLotsInParkingLot(t *testing.T) {
	_, err := NewParkingLot(10)
	if err != nil {
		t.Error("Empty lots are not created in ParkingLot")
	}
}

func TestCannotCreateParkingLotWithCapacityLessThanOne(t *testing.T) {
	_, err := NewParkingLot(-1)
	if err == nil {
		t.Error("Parking lot cannot be created with capacity less than 1")
	}
}

func TestParkCar(t *testing.T) {
	parkingLot, _ := NewParkingLot(1)
	car := Car{
		numberPlate: "KJ-09-AK-123",
	}
	err := parkingLot.park(&car)
	if err != nil {
		t.Errorf("Vehicle not parked")
	}
}

func TestCreateParkingLotStruct(t *testing.T) {
	_, err := NewParkingLot((10))
	if err != nil {
		t.Errorf("Parking Lot is not created")
	}

}

func TestCreateCarStruct(t *testing.T) {
	car := &Car{

		numberPlate: "KJ-09-AK-123",
	}

	if car.numberPlate == "" {
		t.Error("Struct Car cannot be created without NumberPlate")
	}
}

func TestParkingLotCreationWithCapacity(t *testing.T) {
	p, _ := NewParkingLot(1)
	if p.capacity != 1 {
		t.Errorf("ParkingLot has not been created with capacity 1")
	}
}

func TestCheckIfLotIsEmptyBeforeParking(t *testing.T) {
	p, _ := NewParkingLot(1)
	car1 := &Car{
		numberPlate: "KK-09-AK-1234",
	}
	car2 := &Car{
		numberPlate: "KK-09-AK-2341",
	}
	p.park(car1)
	err := p.park(car2)

	if err == nil {
		t.Errorf("Cannot park car since no slots are empty")
	}

}

func TestCheckIfParkingLotIsFull(t *testing.T) {
	car1 := &Car{
		numberPlate: "KK-09-AK-1234",
	}
	car2 := &Car{
		numberPlate: "KK-09-AK-2341",
	}
	p, _ := NewParkingLot(1)
	p.park(car1)
	err := p.park(car2)

	if err == nil {
		t.Errorf("ParkingLot is full")
	}

}

type mockParkingFullReceiverCounter struct {
	parkingFullCalledTimes int
}

func (s *mockParkingFullReceiverCounter) receiveFull(int) {

	s.parkingFullCalledTimes++
}

func TestCheckIfFullReceiverGetsNotifiedWhenCarParkedAfterUnpark(t *testing.T) {

	car1 := &Car{
		numberPlate: "KK-09-AK-1234",
	}
	car2 := &Car{
		numberPlate: "KK-09-AK-2341",
	}

	car3 := &Car{
		numberPlate: "KK-09-AK-234sds",
	}
	s := &mockParkingFullReceiverCounter{}

	p, _ := NewParkingLot(3)
	p.OnFull(s)

	err := p.park(car1)
	if err != nil {
		t.Fatal("car1 should be parked")
	}

	err = p.park(car2)
	if err != nil {
		t.Fatal("car2 should be parked")
	}

	err = p.park(car3)
	if err != nil {
		t.Fatal("car3 should be parked")
	}
	if s.parkingFullCalledTimes != 1 {
		t.Fatalf("parking full reciver should have been called once but was %d", s.parkingFullCalledTimes)
	}

	err = p.unPark(car1)
	if err != nil {
		t.Fatal("car3 should be parked")
	}

	err = p.park(car1)
	if err != nil {
		t.Fatal("car1 should be parked")
	}

	if s.parkingFullCalledTimes != 2 {
		t.Fatalf("parking full reciver should have been called once but was %d", s.parkingFullCalledTimes)
	}

}

func TestUnparkCar(t *testing.T) {
	p, _ := NewParkingLot(1)
	car := &Car{
		numberPlate: "KK-09-AK-2341",
	}
	p.park(car)

	err := p.unPark(car)
	if err != nil {
		t.Errorf("Car not unparked")

	}
}

func TestUnparkCarNotFound(t *testing.T) {
	p, _ := NewParkingLot(1)
	car := &Car{
		numberPlate: "KK-09-AK-2341",
	}
	p.park(car)

	err := p.unPark(car)
	if err != nil {
		t.Errorf("Car not found")
	}
}

func TestCarIsParked(t *testing.T) {
	p, _ := NewParkingLot(1)
	car := &Car{
		numberPlate: "KK-09-AK-2341",
	}
	p.park(car)

	result := p.isParked(car)

	if !result {
		t.Errorf("Car is not parked in the ParkingLot")
	}

}
func TestCheckIfCarAlreadyParked(t *testing.T) {
	p, _ := NewParkingLot(2)
	const numberPlate = "KK-09-AK-1234"
	car1 := &Car{
		numberPlate: numberPlate,
	}
	car2 := &Car{
		numberPlate: numberPlate,
	}
	p.park(car1)
	err := p.park(car2)
	if err == nil {
		t.Errorf("Car is already parked")
	}

}

type mockParkingFullReceiver struct {
	parkingFull bool
}

func (s *mockParkingFullReceiver) receiveFull(int) {

	s.parkingFull = true
}

var car = Car{numberPlate: "AB-12-CD-3456"}

func TestMultipleReceiversShouldBeNotifiedWhenParkingFull(t *testing.T) {
	p, _ := NewParkingLot(1)
	s := &mockParkingFullReceiver{}
	another := &mockParkingFullReceiver{}

	p.OnFull(s)
	p.OnFull(another)
	p.park(&car)

	if !s.parkingFull {
		t.Errorf("The status should be changed to parking_full")
	}
	if !another.parkingFull {
		t.Errorf("Another person should have status parking_full")
	}

}

type mockParkingAvailableReceiver struct {
	receiveCalled bool
}

func (m *mockParkingAvailableReceiver) receiveAvailable() {
	m.receiveCalled = !m.receiveCalled
}

func TestSingleRecieverNotifiedParkingAvailable(t *testing.T) {
	p, _ := NewParkingLot(1)
	parkingAvailableReceiver := mockParkingAvailableReceiver{}
	p.OnAvailable(&parkingAvailableReceiver)
	p.park(&car)
	p.unPark(&car)

	if !parkingAvailableReceiver.receiveCalled {
		t.Errorf("Receive Function not called")
	}
}
func TestUnparkNotParkedCar(t *testing.T) {
	parkinglot, _ := NewParkingLot(2)

	err := parkinglot.unPark(&car)

	if err.Error() != "Car is not found in the parking lot" {
		t.Errorf("parkinglot cannot unpark the not existing car")
	}
}
func TestCannotNofifyMultiplePeopleWhenParkingAvailable(t *testing.T) {
	p, _ := NewParkingLot(1)
	parkingFullReceiver := mockParkingFullReceiver{}
	parkingAvailableReceiver := mockParkingAvailableReceiver{}

	p.OnFull(&parkingFullReceiver)

	p.OnAvailable(&parkingAvailableReceiver)

	p.park(&car)
	if !parkingFullReceiver.parkingFull {
		t.Errorf("Only one person should be notified for parking availability")
	}
	if parkingAvailableReceiver.receiveCalled == true {
		t.Errorf("Parking Available Receiver should not be notified when parking full")
	}

	p.unPark(&car)
	if !parkingFullReceiver.parkingFull {
		t.Errorf("Only one person should be notified for parking availability")
	}

	if !parkingAvailableReceiver.receiveCalled {
		t.Errorf("This receiver has to be Notified when parking is available")
	}

}

func TestMultipleReceiversShouldBeNotifiedWhenParkingAvailable(t *testing.T) {
	p, _ := NewParkingLot(1)
	s := mockParkingAvailableReceiver{}
	another := mockParkingAvailableReceiver{}

	p.OnAvailable(&s)
	p.OnAvailable(&another)

	err := p.park(&car)
	if err != nil {
		t.Fatal("car should be parked")
	}

	err = p.unPark(&car)
	if err != nil {
		t.Fatal("car should be unparked")
	}

	if !s.receiveCalled {
		t.Errorf("the status should be changed to parking_available")
	}
	if !another.receiveCalled {
		t.Errorf("another person should have status parking_available")
	}

}

type ReceiveBothNotification struct {
	notifiedFull      bool
	notifiedAvailable bool
}

func (o *ReceiveBothNotification) receiveFull(int) {
	o.notifiedFull = true
}

func (o *ReceiveBothNotification) receiveAvailable() {
	o.notifiedAvailable = true
}

func TestNotifyOwnerWhenFullAndWhenAvailable(t *testing.T) {
	p, _ := NewParkingLot(1)
	owner := ReceiveBothNotification{}
	p.OnAvailable(&owner)
	p.OnFull(&owner)
	p.park(&car)

	if !owner.notifiedFull {
		t.Errorf("owner should be notified when the parking lot gets full")
	}

	p.unPark(&car)

	if !owner.notifiedAvailable {
		t.Errorf("owner should be notified when the parking lot is available")
	}

}

func TestTwoCarsEqual(t *testing.T) {
	const carNumber = "MH12AA2345"
	car1 := Car{carNumber}
	car2 := Car{carNumber}

	result := car1.isEqual(&car2)

	if !result {
		t.Errorf("car1 should be equal to car2")
	}

}

func TestTwoCarsNotEqual(t *testing.T) {
	car1 := Car{"MH12AA2345"}
	car2 := Car{"UP12HH4009"}

	result := car1.isEqual(&car2)

	if result {
		t.Errorf("car1 should not be equal to car2")
	}

}

func TestParkAfterUnparkForSingleSlot(t *testing.T) {
	parkinglot, _ := NewParkingLot(1)

	parkinglot.park(&car)
	parkinglot.unPark(&car)

	err := parkinglot.park(&car)

	if err != nil {
		t.Errorf("car should park after nupark")
	}
}

func TestCannotParkNil(t *testing.T) {

	parkinglot, _ := NewParkingLot(1)

	err := parkinglot.park(nil)

	if err == nil {
		t.Errorf("car to be parked cannot be nil")
	}

}

func TestCannotUnparkNil(t *testing.T) {

	parkinglot, _ := NewParkingLot(1)

	err := parkinglot.unPark(nil)

	if err == nil {
		t.Error("nil car cannot be unparked")
	}
}
func TestCarShouldGetParkedAferUnpark(t *testing.T) {
	parkinglot, _ := NewParkingLot(2)
	car2 := Car{"AA10AK2345"}

	err1 := parkinglot.park(&car)
	if err1 != nil {
		t.Fatal("car should be parked")
	}
	err1 = parkinglot.park(&car2)
	if err1 != nil {
		t.Fatal("car2 should be parked")
	}

	err1 = parkinglot.unPark(&car2)
	if err1 != nil {
		t.Fatal("car2 should be unparked")
	}
	err := parkinglot.park(&car2)

	if err != nil {
		t.Fatalf("car2 should be parked after unpark : %v", err)
	}
}

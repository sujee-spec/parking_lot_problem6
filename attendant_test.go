package parkinglot

import "testing"

func TestCreateNewAttendant(t *testing.T) {
	parkinglog, _ := NewParkingLot(1)
	_, err := NewAttendant(parkinglog)

	if err != nil {
		t.Error("new attendant should be created")
	}
}

func TestAttedantCannotBeCreatedWithNilParkingLot(t *testing.T) {
	_, err := NewAttendant(nil)

	const expectedError = "newattendant cannot have nil parkinglot"

	if err.Error() != expectedError {
		t.Error("attendant should not create with nil parking lot")
	}
}

func TestParkCarByAttendant(t *testing.T) {
	parkingLot, _ := NewParkingLot(1)
	attendant, _ := NewAttendant(parkingLot)

	err := attendant.Park(&car)

	if err != nil {
		t.Errorf("attendant should be able to park the car")
	}
}

func TestAttendantCannotParkWhenParkingFull(t *testing.T) {
	parkingLot, _ := NewParkingLot(1)
	attendant, _ := NewAttendant(parkingLot)

	expectedError := "parking lot is full, attendant cannot park the car"

	attendant.Park(&Car{"KK10AA1234"})

	err := attendant.Park(&car)

	if err.Error() != expectedError {
		t.Errorf("attendant cannot park when parking full")
	}
}

func TestUnParkCarByAttendant(t *testing.T) {
	parkingLot, _ := NewParkingLot(1)
	attendant, _ := NewAttendant(parkingLot)

	attendant.Park(&car)
	err := attendant.UnPark(&car)

	if err != nil {
		t.Errorf("car should be unparked by attendant")
	}
}

func TestAttendantParkAfterParkingAvailable(t *testing.T) {
	parkingLot, _ := NewParkingLot(1)
	attendant, _ := NewAttendant(parkingLot)

	//park the car
	attendant.Park(&car)

	//unpark the car
	attendant.UnPark(&car)

	//should be able to park again
	err := attendant.Park(&car)

	if err != nil {
		t.Error("should be able to park after parking become available")
	}

}

//Multiple parking lots

func TestAttendantCannotParkNilCar(t *testing.T) {
	parkinglot, _ := NewParkingLot(2)
	attendant, _ := NewAttendant(parkinglot)

	err := attendant.Park(nil)
	expectedError := "attendant cannot park nil car"

	if err.Error() != expectedError {
		t.Errorf("attendant cannot park nil car")
	}
}

func TestAttendantShouldCheckCarIsParkedAfterUnpark(t *testing.T) {
	parkinglot, _ := NewParkingLot(2)
	attendant, _ := NewAttendant(parkinglot)
	car2 := Car{"UM-12-TK-1234"}

	err := attendant.Park(&car)
	if err != nil {
		t.Fatal("car should be parked")
	}

	err = attendant.Park(&car2)
	if err != nil {
		t.Fatal("car2 should be parked")
	}
	err = attendant.UnPark(&car2)
	if err != nil {
		t.Fatal("car2 should be unparked")
	}

	isCarParked := attendant.checkIsCarParked(&car)
	isCar2Parked := attendant.checkIsCarParked(&car2)

	if isCarParked == false {
		t.Error("car should be parked in parking lot")
	}
	if isCar2Parked == true {
		t.Errorf("car is already parked")
	}
}

func TestAttendantCanManageMultipleParkingLots(t *testing.T) {
	parkinglot1, err1 := NewParkingLot(1)
	if err1 != nil {
		t.Fatalf("failed to create parking lot 1: %v", err1)
	}

	parkinglot2, err2 := NewParkingLot(1)
	if err2 != nil {
		t.Fatalf("failed to create parking lot 2: %v", err2)
	}

	_, err := NewAttendant(parkinglot1, parkinglot2)

	if err != nil {
		t.Errorf("expected attendant to be created with multiple parking lots, got error: %v", err)
	}
}

func TestAttendantParkCarInNextParkingLotWithAvailableSlot(t *testing.T) {

	parkinglot1, _ := NewParkingLot(1)
	parkinglot2, _ := NewParkingLot(3)
	attendant, _ := NewAttendant(parkinglot1, parkinglot2)

	parkinglot1.park(&car)
	car2 := Car{"1234567"}
	err := attendant.Park(&car2)

	if err != nil {
		t.Fatal("car2 should be parked at parking lot 2")
	}

	if !parkinglot2.slots[0].car.isEqual(&car2) {
		t.Errorf("car2 should be parked at parking lot 2 ")

	}

}

func TestAttendantIsAbleToUnparkAfterParkForMultipleParkingLots(t *testing.T) {
	parkinglot1, _ := NewParkingLot(1)
	parkinglot2, _ := NewParkingLot(2)
	parkinglot3, _ := NewParkingLot(1)
	attendant, _ := NewAttendant(parkinglot1, parkinglot2, parkinglot3)

	err := parkinglot1.park(&car)
	if err != nil {
		t.Fatal("car should be parking in parkinglot 1")
	}

	anotherCar := &Car{"1234567"}
	err = parkinglot2.park(anotherCar)
	if err != nil {
		t.Fatal("another should be parking in parkinglot 2")
	}

	err = parkinglot3.park(&Car{"09876533"})
	if err != nil {
		t.Fatal("Car{09876533} should be parking in parkinglot 2")
	}

	err = attendant.UnPark(anotherCar)

	if err != nil {
		t.Error("another car should be unparked from parkinglot 2")
	}

}

func TestAttendantCannotParkSameCarAgain(t *testing.T) {
	parkinglot1, err := NewParkingLot(2)

	if err != nil {
		t.Fatal("parkinglot should be created with 2 slot")
	}
	attendant, err := NewAttendant(parkinglot1)
	if err != nil {
		t.Fatal("attendant should be created wiht parkinglot1")
	}

	err = attendant.Park(&car)
	if err != nil {
		t.Fatal("car should get parked in the parking lot")
	}

	err = attendant.Park(&car)

	expectedError := "attendant cannot park already parked car"
	if err.Error() != expectedError {
		t.Error("car cannot be parked again")
	}

}

func TestAttendantCannotUnParkNilCar(t *testing.T) {
	parkinglot1, err := NewParkingLot(2)

	if err != nil {
		t.Fatal("parkinglot should be created with 2 slot")
	}
	attendant, err := NewAttendant(parkinglot1)
	if err != nil {
		t.Fatal("attendant should be created wiht parkinglot1")
	}

	err = attendant.UnPark(nil)
	expectedError := "attendant cannot unpark nil car"
	if err.Error() != expectedError {
		t.Fatal("car should get parked")
	}
}

func TestAddentParkUsingTheEvenPlan(t *testing.T) {
	parkinglot1, _ := NewParkingLot(2)
	parkinglot2, _ := NewParkingLot(2)

	attendant, err := NewAttendantV2(ParkInParkinglotWithLeastVehicles, parkinglot1, parkinglot2)
	if err != nil {
		t.Fatal("attendant should be created with evenparking plan")
	}

	err = attendant.Park(&car)
	if err != nil {
		t.Fatal("car should be parked in parkinglot1 in first slot")
	}
	if attendant.Parkinglots[0].slots[0].car != &car {
		t.Fatal("count of parkinglot1 should increase to 1 after car gets parked")
	}

	car2 := &Car{"car2"}
	err = attendant.Park(car2)
	if err != nil {
		t.Fatal("car should be parked in parkinglot2 in first slot")
	}
	if attendant.Parkinglots[1].slots[0].car != car2 {
		t.Fatal("count of parkinglot2 should increase to 1 after car2 gets parked")
	}
}

func TestAttendantsAccessTheSameParkinglotReference(t *testing.T) {
	parkinglot1, _ := NewParkingLot(2)
	parkinglot2, _ := NewParkingLot(2)

	car2 := &Car{numberPlate: "car2"}
	car3 := &Car{numberPlate: "car3"}
	car4 := &Car{numberPlate: "car4"}

	simpleAttendant, _ := NewAttendantV2(ParkInFirstAvailableParkinglot, parkinglot1, parkinglot2)
	complexAttendant, _ := NewAttendantV2(ParkInParkinglotWithLeastVehicles, parkinglot1, parkinglot2)

	err := simpleAttendant.Park(&car)
	if err != nil {
		t.Fatalf("simple attendant should be able to park car in parkinglot 1,  %v", err)
	}
	if !parkinglot1.slots[0].car.isEqual(&car) {
		t.Fatal("car should be parked in parkinglot 1, first slot")
	}

	err = simpleAttendant.Park(car2)
	if err != nil {
		t.Fatalf("simple attendant should be able to park car2 in parkinglot 1,  %v", err)
	}
	if car.isEqual(parkinglot1.slots[0].car) == false {
		t.Fatal("attendant should park in lot1, second slot")
	}

	err = complexAttendant.Park(car3)
	if err != nil {
		t.Fatal("car3 should be parked")
	}
	if !parkinglot2.slots[0].car.isEqual(car3) {
		t.Fatal("car3 should be parked in parkinglot 2, first slot")
	}

	err = simpleAttendant.UnPark(&car)
	if err != nil {
		t.Fatalf("simpleattendant should unpark the car %v", err)
	}
	if parkinglot1.slots[0].occupied == true {
		t.Fatal("parkinglot 1, first slot should be empty after car unpark")
	}

	err = simpleAttendant.UnPark(car2)
	if err != nil {
		t.Fatal("ca2 should get unpark")
	}
	if parkinglot1.slots[1].occupied == true {
		t.Fatal("parkinglot 1, second slot should be empty afer car2 unpark")
	}

	err = complexAttendant.Park(car4)
	if err != nil {
		t.Fatal("car4 should get park")
	}
	if car4.isEqual(parkinglot1.slots[0].car) == false {
		t.Fatalf("car should have been parked in first slot of parkinglot 1")
	}
}

func TestAttendParkInLotWithMostCapacity(t *testing.T) {
	parkinglot1, _ := NewParkingLot(2)
	parkinglot2, _ := NewParkingLot(3)

	car2 := &Car{numberPlate: "car2"}
	car3 := &Car{numberPlate: "car3"}
	evenTypeAttendant, _ := NewAttendantV2(ParkInParkinglotWithLeastVehicles, parkinglot1, parkinglot2)
	mostTypeAttendant, _ := NewAttendantV2(MostOccupiedParking, parkinglot1, parkinglot2)

	err := evenTypeAttendant.Park(&car)
	if err != nil {
		t.Fatal("car should be parked in parkinglot1")
	}

	err = evenTypeAttendant.Park(car2)
	if err != nil {
		t.Fatal("car2 should be parked in the parkinglot")
	}
	if !parkinglot2.slots[0].car.isEqual(car2) {
		t.Fatal("car2 should be parked in the parkinglot2, first slot")
	}

	err = evenTypeAttendant.UnPark(&car)
	if err != nil {
		t.Fatal("evenattendant should be unpark car")
	}
	if !parkinglot1.slots[0].isEmpty() {
		t.Fatal("after car unpark, the parkinglot1, first slot should be empty")
	}

	err = mostTypeAttendant.Park(car3)
	if err != nil {
		t.Fatal("mosttypeattendant should be able to park car3")
	}
	if !parkinglot2.slots[1].car.isEqual(car3) {
		t.Error("mosttypeattendant should park car3 in parkiglot2, second slot")
	}

}

/*What do you want to test here? */
func TestMultipleAttendantsSubscribeToBothParkingNotifications(t *testing.T) {
	parkinglot1, _ := NewParkingLot(2)
	parkinglot2, _ := NewParkingLot(2)

	car2 := &Car{numberPlate: "car2"}

	simpleAttendant, _ := NewAttendantV2(ParkInFirstAvailableParkinglot, parkinglot1, parkinglot2)
	complexAttendant, _ := NewAttendantV2(ParkInParkinglotWithLeastVehicles, parkinglot1, parkinglot2)

	err := simpleAttendant.Park(&car)
	if err != nil {
		t.Fatalf("simple attendant should be able to park car in parkinglot 1")
	}
	//TODO is this validation already done in other dedicated test case which validates the parking plan ?
	if !parkinglot1.slots[0].car.isEqual(&car) {
		t.Fatal("car should be parked in parkinglot 1, first slot")
	}

	// attendant gets notification for parking full
	err = simpleAttendant.Park(car2)
	if err != nil {
		t.Fatalf("simple attendant should be able to park car2 in parkinglot 1")
	}
	if car.isEqual(parkinglot1.slots[0].car) == false {
		t.Fatal("attendant should park in lot1, first slot")
	}
	if simpleAttendant.parkingFullStatus[0] != true {
		t.Fatal("simple attendant should receive parking full notification")
	}
	if complexAttendant.parkingFullStatus[0] != true {
		t.Fatal("complex attendant should receive parking full notification")
	}

	// attendant gets notification for parking available
	err = simpleAttendant.UnPark(&car)
	if err != nil {
		t.Fatalf("simpleattendant should unpark the car %v", err)
	}
	if simpleAttendant.parkingFullStatus[0] != false {
		t.Fatal("simple attendant should receive parking full notification")
	}
	if complexAttendant.parkingFullStatus[0] != false {
		t.Fatal("complex attendant should receive parking full notification")
	}

}

package parkinglot

import "testing"

func TestAttedantCannotBeCreatedWithNilParkingLot(t *testing.T) {
	_, err := NewAttendant(nil)

	const expectedError = "newattendant cannot have nil parkinglot"

	if err.Error() != expectedError {
		t.Error("attendant should not create with nil parking lot")
	}
}

func TestAttendantCannotParkWhenParkingFull(t *testing.T) {
	parkingLot, _ := NewParkingLot(1)
	attendant, _ := NewAttendant(parkingLot)

	expectedError := "all parkinglots are full"

	attendant.Park(&Car{"KK10AA1234"})

	err := attendant.Park(&car)

	if err.Error() != expectedError {
		t.Errorf("attendant cannot park when parking full")
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

func TestAttendantCannotUnparkACarThatWasNeverParkedInParkingLot(t *testing.T) {
	parkinglot, _ := NewParkingLot(2)

	attendant, err := NewAttendantV2(ParkInFirstAvailableParkinglot, parkinglot)
	if err != nil {
		t.Fatal("attendant should be created")
	}

	err = attendant.UnPark(&car)
	const expectedError = "attendant cannot unpark, car not found in parlinglots"
	if err.Error() != expectedError {
		t.Error("attendant cannot unpark the car that was never parked")
	}
}

func TestMostFilledParkingAttendantShouldParkInMostFilledParkingLot(t *testing.T) {
	parkinglot1, _ := NewParkingLot(1)
	parkinglot2, _ := NewParkingLot(2)

	mostFilledTypeAttendant, err := NewAttendantV2(MostOccupiedParking, parkinglot1, parkinglot2)
	if err != nil {
		t.Fatal("attendant should be created")
	}

	err = mostFilledTypeAttendant.Park(&car)
	if err != nil {
		t.Fatal("attendant should park the car")
	}

	car2 := &Car{"car2"}
	err = mostFilledTypeAttendant.Park(car2)
	if err != nil {
		t.Fatal("attendant should be able to park car2")
	}
	if !mostFilledTypeAttendant.parkinglots[1].slots[0].car.isEqual(car2) {
		t.Error("car2 shuld be parked in parkinglot2, first slot")
	}

}

func TestAttendantCannotFindParkinglotWithLeastFilledVehicle(t *testing.T) {

	parkinglot, err := NewParkingLot(1)
	if err != nil {
		t.Fatal("parkinglot should be created")
	}

	attendant, err := NewAttendantV2(ParkInParkinglotWithLeastVehicles, parkinglot)
	if err != nil {
		t.Fatal("attendant should be created")
	}
	err = attendant.Park(&car)
	if err != nil {
		t.Fatal("attendant should be able to park the car")
	}

	err = attendant.Park(&Car{"car2"})
	const expectedError = "no parkinglot with least vehicle available"
	if err.Error() != expectedError {
		t.Error("attendant should not be able to park the car")
	}
}

func TestAttendantCannotFindParkinglotWithMostFilledVehicle(t *testing.T) {

	parkinglot, err := NewParkingLot(1)
	if err != nil {
		t.Fatal("parkinglot should be created")
	}

	attendant, err := NewAttendantV2(MostOccupiedParking, parkinglot)
	if err != nil {
		t.Fatal("attendant should be created")
	}
	err = attendant.Park(&car)
	if err != nil {
		t.Fatal("attendant should be able to park the car")
	}

	err = attendant.Park(&Car{"car2"})
	const expectedError = "no parkinglot with most vehicle available"
	if err.Error() != expectedError {
		t.Error("attendant should not be able to park the car")
	}
}

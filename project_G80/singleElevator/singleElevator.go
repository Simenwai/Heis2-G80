package singleElevator

//main for single elevator

import (
	"time"

	"project/types"
	"project/singleElevator/fsm"
	"project/singleElevator/fsm/elevio"
)

func SingleElevator(
	simulatorPort string,
	elevLight <-chan types.LightEvent,
	elevOrderIn <-chan types.ButtonEvent,
	elevCostReqIn <-chan types.Order,
	elevOrderCompleted chan<- types.ButtonEvent,
	elevCostReqOut chan<- types.Cost,  
	elevObstruction chan<- bool, 
	elevTimeOfFloorEvent chan<- time.Time) {

	elevio.Init("localhost:"+simulatorPort, types.NUM_FLOORS)

	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	timeout := make(chan bool)

	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)

	fsm.fsm_onInitBetweenFloors()

	for {
		select {
		case buttonEvent := <-elevOrderIn:
			fsm.fsm_onRequestButtonPress(buttonEvent.Floor, buttonEvent.Button, timeout)
		case floor := <-drv_floors:
			elevTimeOfFloorEvent <- time.Now()
			fsm.fsm_onFloorArrival(floor, timeout, elevOrderCompleted)
		case <-timeout:
			fsm.fsm_onDoorTimeout(elevOrderCompleted)
		//case buttonEvent := <-elevCostReqIn:
		//	go elevator.CostCalculation(buttonEvent, elevCostReqOut)
		case lightEvent := <-elevLight:
			elevio.SetButtonLamp(lightEvent.Light.Button, lightEvent.Light.Floor, lightEvent.Switch)
		case obstructionEvent := <- drv_obstr:
			fsm.Fsm_onObstructionSwitch(obstructionEvent)
			elevObstruction <- obstructionEvent
		}
	}
}

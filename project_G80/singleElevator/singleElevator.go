package singleElevator

//main for single elevator

import (
	"fmt"
	//"time"

	"project/singleElevator/fsm"
	"project/singleElevator/fsm/elevio"
	"project/types"
)

func SingleElevator(simulatorPort string) {
	
	drv_floors := make(chan int)
	drv_obstr := make(chan bool)
	drv_buttons := make(chan types.ButtonEvent)
	drv_timer := make(chan bool)
	
	go elevio.PollFloorSensor(drv_floors)
	go elevio.PollObstructionSwitch(drv_obstr)
	go elevio.PollButtons(drv_buttons)

	
	fmt.Println("sas")
	for {
		select {
		case buttonEvent := <-drv_buttons:
			fmt.Println("knap trykt")
			fsm.Fsm_onRequestButtonPress(buttonEvent.Floor, buttonEvent.Button, drv_timer)
		case floor := <-drv_floors:
		//	elevTimeOfFloorEvent := time.Now()
			fsm.Fsm_onFloorArrival(floor, drv_timer)
		case <-drv_timer:
			fsm.Fsm_onDoorTimeout()
		//case buttonEvent := <-elevCostReqIn:
		//	go elevator.CostCalculation(buttonEvent, elevCostReqOut)
		//case lightEvent := <-elevLight:
		//	elevio.SetButtonLamp(lightEvent.Light.Button, lightEvent.Light.Floor, lightEvent.Switch)
		case obstructionEvent := <- drv_obstr:
			fsm.Fsm_onObstructionSwitch(obstructionEvent)
		}
		fmt.Println("www")
	}
}

package fsm

import (
	"fmt"
	"time"

	requests "project/requests"
	elevio "project/singleElevator/fsm/elevio"
	"project/types"
)

var elev types.Elevator

//sjekk første funksjon i fsm: func fsm_init

func SetAllLights(){
	for floor := 0; floor < types.NUM_FLOORS; floor++ {
		for btn := 0; btn < types.NUM_BUTTONS; btn++{
			elevio.SetButtonLamp(types.ButtonType(btn), floor, elev.Requests[floor][btn])
		}
	}
}

func TimerDoor(timedOut chan<- bool){
	for {
		time.Sleep((time.Duration(3*time.Second)))
		if !elev.Obstruction{
			break
		}
	}

	timedOut <- true
}

func OpenDoor(timedOut chan<- bool){
	elevio.SetDoorOpenLamp(true)
	elev.Behaviour = types.EB_DoorOpen
	go TimerDoor(timedOut)
	
}


func Fsm_onInitBetweenFloors(){
	elevio.SetDoorOpenLamp(false)
	elevio.SetMotorDirection(types.MD_Down)
	elev.Dirn = types.MD_Down
	elev.Behaviour = types.EB_Moving
	SetAllLights()
	a := -5

	for a != 0{
		a = elevio.GetFloor()
	}
	elevio.SetMotorDirection(types.MD_Stop)

	elev.Dirn = types.MD_Stop
	elev.Behaviour = types.EB_Idle
	fmt.Println("ss")
}

func Fsm_onRequestButtonPress(btn_floor int, btn_type types.ButtonType, timedOut chan<- bool){
	fmt.Println(elev.Behaviour)
	fmt.Println(elev.Floor)
	fmt.Println(elev.Dirn)

	
	
	switch elev.Behaviour{
	case types.EB_DoorOpen:
		if elev.Floor == btn_floor{
			OpenDoor(timedOut)

		} else {
			elev.Requests[btn_floor][btn_type] = true
		}
		break
	
	case types.EB_Moving:
		elev.Requests[btn_floor][btn_type] = true
		break
	
	case types.EB_Idle:
		if elev.Floor == btn_floor{
			OpenDoor(timedOut)
		} else {
			elev.Requests[btn_floor][btn_type] = true
			elev.Dirn = requests.Requests_nextAction(elev)
			elevio.SetMotorDirection(elev.Dirn)
			elev.Behaviour = types.EB_Moving
		}
		break
	}
	SetAllLights()
	//print new state
}

func Fsm_onFloorArrival(newFloor int, timedOut chan<- bool) {
	elev.Floor = newFloor

	elevio.SetFloorIndicator(elev.Floor)

	switch elev.Behaviour {
	case types.EB_Moving:
		if requests.Requests_shouldStop(elev) {
			elev = requests.Requests_clearAtCurrentFloor(elev)
			elevio.SetMotorDirection(types.MD_Stop)
			OpenDoor(timedOut)
		}
		break
	default:
		fmt.Println("Floor arrival break")
		break
	}
	SetAllLights()

	//print new state
}

func Fsm_onDoorTimeout() {
	switch elev.Behaviour {
	case types.EB_DoorOpen:
		elev.Dirn = requests.Requests_nextAction(elev)
		elev = requests.Requests_clearAtCurrentFloor(elev)
		elevio.SetDoorOpenLamp(false)
		elevio.SetMotorDirection(elev.Dirn)

		if elev.Dirn == types.MD_Stop {
			elev.Behaviour = types.EB_Idle

		} else {
			elev.Behaviour = types.EB_Moving
		}

		break
	default:
		break
	}
}

func Fsm_onObstructionSwitch(obstructionEvent bool){
		elev.Obstruction = obstructionEvent

}


package fsm

import (
	"time"

	"project/types"
	"project/singleElevator/fsm/elevio"
	"project/requests"
)

var elev types.Elevator

//sjekk første funksjon i fsm: func fsm_init

func setAllLights(){
	for floor := 0; floor < types.NUM_FLOORS; floor++ {
		for btn := 0; btn < types.NUM_BUTTONS; btn++{
			elevio.SetButtonLamp(types.ButtonType(btn), floor, elev.Requests[floor][btn])
		}
	}
}

func timerDoor(isTimedOut chan <- bool){
	for{
		time.Sleep(time.Duration(3 * time.Second))
		if !elev.Obstruction{
			break
		} 
	}
	isTimedOut <- true	
}

func openDoor(isTimedOut chan<- bool){
	elevio.SetDoorOpenLamp(true)
	go timerDoor(isTimedOut)
	elev.Behaviour = types.EB_DoorOpen
}


func fsm_onInitBetweenFloors(){
	elevio.SetMotorDirection(types.MD_Down)
	elev.Dirn = types.MD_Down
	elev.Behaviour = types.EB_Moving
	setAllLights()	
}

func fsm_onRequestButtonPress(btn_floor int, btn_type types.ButtonType, isTimedOut chan<- bool ){
	//print(elevator)?
	
	switch elev.Behaviour{
	case types.EB_DoorOpen:
		if elev.Floor == btn_floor{
			go timerDoor(isTimedOut)
		} else {
			elev.Requests[btn_floor][btn_type] = true
		}
		break
	
	case types.EB_Moving:
		elev.Requests[btn_floor][btn_type] = true
		break
	
	case types.EB_Idle:
		if elev.Floor == btn_floor{
			openDoor(isTimedOut)
		} else {
			elev.Requests[btn_floor][btn_type] = true
			elev.Dirn = requests_chooseDirection(elev)
			elevio.SetMotorDirection(elev.Dirn)
			elev.Behaviour = types.EB_Moving
		}
		break
	}
	setAllLights()
	//print new state
}

func fsm_onFloorArrival(newFloor int, isTimedOut chan<- bool, elevOrderCompleted chan<- types.ButtonEvent) {
	elev.Floor = newFloor

	elevio.SetFloorIndicator(elev.Floor)

	switch elev.Behaviour {
	case types.EB_Moving:
		if requests_shouldStop(elev) {
			elev = requests_clearAtCurrentFloor(elev, elevOrderCompleted)
			elevio.SetMotorDirection(types.MD_Stop)
			openDoor(isTimedOut)
		}
		break
	default:
		break
	}

	//print new state
}

func fsm_onDoorTimeout(elevOrderCompleted chan<- types.ButtonEvent) {
	switch elev.Behaviour {
	case types.EB_DoorOpen:
		elev.Dirn = requests_chooseDirection(elev)
		elev = requests_clearAtCurrentFloor(elev, elevOrderCompleted)
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

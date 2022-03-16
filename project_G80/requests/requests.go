package requests

import (
	"project/types"
)

func RequestsAbove(e types.Elevator) bool {
	for f := e.Floor + 1; f < types.NUM_FLOORS; f++ {
		for btn := 0; btn < types.NUM_BUTTONS; btn++ {
			if e.Requests[f][btn] == true {
				return true
			}
		}
	}
	return false
}

func RequestsBelow(e types.Elevator) bool {
	for f := 0; f < e.Floor; f++ {
		for btn := 0; btn < types.NUM_BUTTONS; btn++ {
			if e.Requests[f][btn] == true {
				return true
			}
		}
	}
	return false
}

func Requests_nextAction(e types.Elevator) types.MotorDirection {
	switch e.Dirn {
	case types.MD_Up:
		if RequestsAbove(e){
			return types.MD_Up
		} else if RequestsBelow(e){
			return types.MD_Down
		} else {
			return types.MD_Stop
		}
	case types.MD_Down:
		if RequestsBelow(e){
			return types.MD_Down
		} else if RequestsAbove(e) {
			return types.MD_Up
		} else {
			return types.MD_Stop
		}
	case types.MD_Stop:
		if RequestsBelow(e) {
			return types.MD_Down
		} else if RequestsAbove(e) {
			return types.MD_Up
		} else {
			return types.MD_Stop
		}
	default:
		return types.MD_Stop
	}
}

func Requests_chooseDirection(e types.Elevator) types.MotorDirection {
	switch e.Dirn {
	case types.MD_Up:
		if RequestsAbove(e) {
			return types.MD_Up
		} else if RequestsBelow(e) {
			return types.MD_Down
		} else {
			return types.MD_Stop
		}
	case types.MD_Down:
		if RequestsBelow(e) {
			return types.MD_Down
		} else if RequestsAbove(e) {
			return types.MD_Up
		} else {
			return types.MD_Stop
		}
	case types.MD_Stop:
		if RequestsBelow(e) {
			return types.MD_Down
		} else if RequestsAbove(e) {
			return types.MD_Up
		} else {
			return types.MD_Stop
		}
	default:
		return types.MD_Stop
	}
}


func Requests_shouldStop(e types.Elevator) bool {
	switch e.Dirn {
	case types.MD_Down:
		return e.Requests[e.Floor][types.BT_HallDown] ||
			e.Requests[e.Floor][types.BT_Cab] ||
			!RequestsBelow(e)
	case types.MD_Up:
		return e.Requests[e.Floor][types.BT_HallUp] ||
			e.Requests[e.Floor][types.BT_Cab] ||
			!RequestsAbove(e)
	case types.MD_Stop:
		return true
	}
	return true
}
/*
func request_shouldClearImmediately(e types.Elevator, btn_floor int, btn_type types.ButtonType) int{
	switch e.Config.CRVariant {
	case CV_types.CV_All:
		e.Floor == btn_floor
	case CV_types.CV_InDirn:
		return (e.Floor == btn_floor) && 
		(
			(e.dirn == types.MD_Up && btn_ype == types.BT_HallUp) ||
			(e.dirn == types.MD_Down && btn_ype == types.BT_HallDown) ||
			(e.dirn == types.MD_Stop) ||
			(btn_type == types.BT_Cab)
		)
	default:
		return 0
	}
}
*/

func Requests_clearAtCurrentFloor(e types.Elevator/*, elevOrderCompleted chan<- types.ButtonEvent*/) types.Elevator {
	switch e.Config.CRVariant {
	case types.CV_All:
		for btn := 0; btn < types.NUM_BUTTONS; btn++ {
			e.Requests[e.Floor][btn] = false
			/*elevOrderCompleted <- types.ButtonEvent{
				Floor:  e.Floor,
				Button: types.ButtonType(btn)}*/
		}
		break
	case types.CV_InDirn:
		e.Requests[e.Floor][types.BT_Cab] = false
		//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_Cab}
		switch e.Dirn {
		case types.MD_Up:
			e.Requests[e.Floor][types.BT_HallUp] = false
			//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_HallUp}
			if !RequestsAbove(e) {
				e.Requests[e.Floor][types.BT_HallDown] = false
				//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_HallDown}
			}
			break
		case types.MD_Down:
			e.Requests[e.Floor][types.BT_HallDown] = false
			//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_HallDown}
			if !RequestsBelow(e) {
				e.Requests[e.Floor][types.BT_HallUp] = false
				//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_HallUp}
			}
			break
		case types.MD_Stop:
		default:
			e.Requests[e.Floor][types.BT_HallUp] = false
			e.Requests[e.Floor][types.BT_HallDown] = false
			//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_HallUp}
			//elevOrderCompleted <- types.ButtonEvent{Floor: e.Floor, Button: types.BT_HallDown}
			break
		}
		break
	default:
		break
	}
	return e
}


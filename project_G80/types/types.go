package types

import (
	"time"
)

const NUM_FLOORS = 4
const NUM_BUTTONS = 3
const NUM_ELEVATORS = 3

type MotorDirection int

const (
	MD_Up   MotorDirection = 1
	MD_Down                = -1
	MD_Stop                = 0
)

type ElevatorBehaviour int

const (
	EB_Idle ElevatorBehaviour = iota
	EB_DoorOpen
	EB_Moving
)

type ClearRequestVariant int

const (
	CV_All ClearRequestVariant = iota
	CV_InDirn
)

type Elevator struct {
	Floor     int
	Dirn      MotorDirection
	Requests  [NUM_FLOORS][NUM_BUTTONS]bool
	Behaviour ElevatorBehaviour
	Obstruction bool
	Config    struct {
		CRVariant     ClearRequestVariant
		doorOpenDurationSeconds int64
	}
}

type ButtonType int

const (
	BT_HallUp   ButtonType = 0
	BT_HallDown            = 1
	BT_Cab                 = 2
)

type ButtonEvent struct {
	Floor  int
	Button ButtonType
}

type Order struct {
	ID                   int
	Taker                int
	Button               ButtonEvent
	TimeFromDistribution time.Time
	Distributed 		 bool
	Completed            bool
}

type OrderMessage struct{
	ID int
	OrderStruct Order
}

type Cost struct {
	ID     int
	Sender int
	Value  int
}

type CostMessage struct{
	ID int
	CostStruct Cost
}

type LightEvent struct {
	Light 	ButtonEvent
	Switch		bool
}
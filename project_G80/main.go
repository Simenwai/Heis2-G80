package main

import (
	"project/singleElevator"
	//"project/types"
	//"time"
	"project/singleElevator/fsm/elevio"
	"fmt"
	"project/singleElevator/fsm"
)


func main(){
	

	port := "localhost:15657"
	elevio.Init(port,4)
	fsm.Fsm_onInitBetweenFloors()
	fmt.Println("hei")

	go singleElevator.SingleElevator(port)

	select{

	}
}

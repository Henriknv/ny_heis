package main

import . "./network"
import . "./elev"
import . "./constants"

import (
	//. "fmt"
	. "time"
)

func main() {

	send_ch := make(chan Elev_info, 100)
	receive_ch := make(chan Elev_info, 100)
	local_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 100)
	rem_local_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 100)
	calculate_order_ch := make(chan map[string]Elev_info, 100)
	lost_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 100)
	next_order_ch := make(chan int, 100)
	elev_dir_ch := make(chan int, 100)

	//Second return type is error handling.

	local_addr, _ := Udp_init(send_ch, receive_ch)
	Elevator_init()

	go Get_local_orders(local_order_ch, rem_local_order_ch, lost_order_ch)
	go Broadcast_orders(local_order_ch, send_ch, local_addr, elev_dir_ch)
	go Get_network_orders(receive_ch, calculate_order_ch, lost_order_ch)
	go Calculate_next_order(calculate_order_ch, next_order_ch, local_addr)
	go Execute_orders(next_order_ch, elev_dir_ch)

	for {

		Sleep(1 * Millisecond)

	}
}

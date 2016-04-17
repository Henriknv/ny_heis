package main

import . "./network"
import . "./elev"
import . "./constants"

import (
	//. "fmt"
	. "time"
)

func main() {

	send_ch := make(chan Elev_info, 1)
	receive_ch := make(chan Elev_info, 1)
	local_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 1)
	rem_local_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 1)
	calculate_order_ch := make(chan map[string]Elev_info, 1)
	lost_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 1)
	next_order_ch := make(chan int, 1)
	elev_dir_ch := make(chan int, 1)
	system_update_ch := make(chan map[string]Elev_info, 1)

	//Second return type is error handling.

	local_addr, _ := Udp_init(send_ch, receive_ch)
	Elevator_init()

	go Get_local_orders(local_order_ch, rem_local_order_ch, lost_order_ch)
	go Broadcast_orders(local_order_ch, send_ch, local_addr, elev_dir_ch)
	go Get_network_orders(receive_ch, calculate_order_ch, lost_order_ch, system_update_ch)
	go Calculate_next_order(calculate_order_ch, next_order_ch, local_addr)
	go Execute_orders(next_order_ch, elev_dir_ch)
	go Update_orders_and_lights(system_update_ch, rem_local_order_ch, local_addr)

	for {

		Sleep(1 * Millisecond)

	}
}

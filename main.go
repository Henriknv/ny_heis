package main

import . "./network"
import . "./elev"
import . "./constants"

import (
	. "fmt"
	. "time"
)

func main() {

	send_chan := make(chan Elev_info, 100)
	receive_chan := make(chan Elev_info, 100)
	local_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 100)
	rem_local_order_ch := make(chan [N_FLOORS][N_BUTTONS]int, 100)

	Udp_init(send_chan, receive_chan)
	Elevator_init()

	go Get_local_orders(local_order_ch, rem_local_order_ch)
	go Broadcast_orders(local_order_ch, send_chan)

	for {
		msg := <-receive_chan
		Println(msg)

		Sleep(1 * Millisecond)

	}
}

package elev

import . ".././driver"
import . ".././constants"
import . ".././network"

import (
	. "fmt"
	. "time"
)

func Elevator_init() {

	Elev_init()

	Elev_set_motor_direction(0)

	if Elev_get_floor_sensor_signal() == -1 {
		Elev_set_motor_direction(-1)
		for {
			if Elev_get_floor_sensor_signal() > -1 {
				Elev_set_motor_direction(0)
				Elev_set_floor_indicator(Elev_get_floor_sensor_signal())
				break
			}
		}
	}
}

func Get_local_orders(local_order_ch chan<- [N_FLOORS][N_BUTTONS]int, rem_local_order_ch <-chan [N_FLOORS][N_BUTTONS]int, lost_order_ch <-chan [N_FLOORS][N_BUTTONS]int) {

	var new_local_order_matrix [N_FLOORS][N_BUTTONS]int

	for {

		for i := 0; i < N_FLOORS; i++ {

			for j := 0; j < N_BUTTONS; j++ {

				if Elev_get_button_signal(j, i) {

					new_local_order_matrix[i][j] = 1

				}
			}
		}

		select {

		case lost_order_matrix := <-lost_order_ch:

			for i := 0; i < N_FLOORS; i++ {

				for j := 0; j < N_BUTTONS-1; j++ {

					if lost_order_matrix[i][j] == 1 {

						new_local_order_matrix[i][j] = 1
					}
				}
			}

			local_order_ch <- new_local_order_matrix

		case rem_local_order_matrix := <-rem_local_order_ch:

			new_local_order_matrix = rem_local_order_matrix

			for i := 0; i < N_FLOORS; i++ {

				for j := 0; j < N_BUTTONS; j++ {

					if Elev_get_button_signal(j, i) {

						new_local_order_matrix[i][j] = 1

					}
				}
			}

		case local_order_ch <- new_local_order_matrix:

		}

		Sleep(1 * Millisecond)

	}
}

func Broadcast_orders(local_order_ch <-chan [N_FLOORS][N_BUTTONS]int, send_ch chan<- Elev_info, local_addr string) {

	var floor int
	var dir int

	for {

		if Elev_get_floor_sensor_signal() == -1 {
			floor = Elev_get_floor_sensor_signal()
		}

		select {

		// case Get direction from Execute_orders:

		case local_order_matrix := <-local_order_ch:

			send_ch <- Elev_info{Elev_id: local_addr, Alive_counter: ALIVE_COUNTER, Floor: floor, Dir: dir, Local_order_matrix: local_order_matrix}

		}

		Sleep(10 * Millisecond)

	}
}

//func Get_network_orders(receive_ch <-chan Elev_info, calculate_order_ch chan<- map[string]Elev_info, lost_order_ch chan<- [N_FLOORS][N_BUTTONS]int)

func Get_network_orders(receive_ch <-chan Elev_info, lost_order_ch chan<- [N_FLOORS][N_BUTTONS]int) {

	online_elevators := make(map[string]Elev_info)

	var temp_elev Elev_info

	for {

		select {
		case new_info := <-receive_ch:

			online_elevators[new_info.Elev_id] = new_info

			//Disconnected elevator handling:

			for elevator := range online_elevators {

				temp_elev = online_elevators[elevator]
				temp_elev.Alive_counter = temp_elev.Alive_counter - 1

				online_elevators[elevator] = temp_elev

				if online_elevators[elevator].Alive_counter < 0 {

					lost_order_ch <- online_elevators[elevator].Local_order_matrix

					delete(online_elevators, elevator)

				}
			}
		}

		// select {

		// case calculate_order_ch <- online_elevators:

		// }

		Println(online_elevators)

		Sleep(1 * Millisecond)

	}
}

package elev

import . ".././driver"
import . ".././constants"
import . ".././network"
import .".././fileio"

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

	new_local_order_matrix := Read()
	

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

			Write(new_local_order_matrix)

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

			Write(new_local_order_matrix)

			local_order_ch <- new_local_order_matrix

		case local_order_ch <- new_local_order_matrix:
			
		}
	}
}

func Broadcast_orders(local_order_ch <-chan [N_FLOORS][N_BUTTONS]int, send_ch chan<- Elev_info, local_addr string, elev_dir_ch <-chan int) {

	var floor int
	var prev_floor int
	prev_dir := DIR_IDLE
	dir := DIR_IDLE


	for {

		if Elev_get_floor_sensor_signal() != LIMBO {
			floor = Elev_get_floor_sensor_signal()

		}



		prev_floor = floor - dir

		select {

		// case Get direction from Execute_orders:
		case temp_dir := <-elev_dir_ch:	

			if temp_dir != dir{

				prev_dir = dir
				dir = temp_dir

			}
		

		case local_order_matrix := <-local_order_ch:
			
			send_ch <- Elev_info{Elev_id: local_addr, Alive_counter: ALIVE_COUNTER, Floor: floor, Prev_floor: prev_floor, Dir: dir, Prev_dir: prev_dir, Local_order_matrix: local_order_matrix}

		}

		Sleep(3*Millisecond)

	}
}

func Get_network_orders(receive_ch <-chan Elev_info, calculate_order_ch chan<- map[string]Elev_info, lost_order_ch chan<- [N_FLOORS][N_BUTTONS]int, system_update_ch chan <- map[string]Elev_info) {

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

		select {

		case system_update_ch <- online_elevators:

		case calculate_order_ch <- online_elevators:

		}
	}
}

func abs_val(val int) int {

	if val < 0 {
		return -val
	}
	return val

}

func calculate_cost(current_floor int, prev_floor int, target_floor int, prev_dir int, button_type int) (cost int) {

	dir := current_floor - prev_floor
	floor_dif := target_floor - current_floor 

	if dir < 0{
		dir = DIR_DOWN
	}else if dir > 0{
		dir = DIR_UP
	}else{
		dir = DIR_IDLE
	}

	
	if dir == DIR_DOWN && button_type == EXT_UP_BUTTONS{

		cost = TURN_COST + abs_val(floor_dif)*FLOOR_COST + 1

	}else if dir == DIR_UP && button_type == EXT_DOWN_BUTTONS{

		cost = TURN_COST + abs_val(floor_dif)*FLOOR_COST + 1

	} else if prev_dir == DIR_DOWN && button_type == EXT_UP_BUTTONS{

		cost = TURN_COST + abs_val(floor_dif)*FLOOR_COST + 1

	}else if prev_dir == DIR_UP && button_type == EXT_DOWN_BUTTONS{

		cost = TURN_COST + abs_val(floor_dif)*FLOOR_COST + 1
		
	}else{

		cost = abs_val(floor_dif)*FLOOR_COST + 1

	}


	// if dir == DIR_DOWN && button_type == EXT_UP_BUTTONS{

	// 	cost = TURN_COST + FLOOR_COST*(N_FLOORS-1-abs_val(floor_dif)) + FLOOR_COST*(N_FLOORS- abs_val(floor_dif)) + 1


	// }else if dir == DIR_UP && button_type == EXT_DOWN_BUTTONS{

	// 	cost = TURN_COST + FLOOR_COST*(N_FLOORS-1-abs_val(floor_dif)) +  FLOOR_COST*(N_FLOORS- abs_val(floor_dif)) + 1

	// } else if prev_dir == DIR_DOWN && button_type == EXT_UP_BUTTONS{

	// 	cost = TURN_COST + FLOOR_COST*(N_FLOORS-1-abs_val(floor_dif)) +  FLOOR_COST*(N_FLOORS- abs_val(floor_dif)) + 1

	// }else if prev_dir == DIR_UP && button_type == EXT_DOWN_BUTTONS{

	// 	cost = TURN_COST + FLOOR_COST*(N_FLOORS-1-abs_val(floor_dif)) +  FLOOR_COST*(N_FLOORS- abs_val(floor_dif)) + 1
	// }else{
	// 	cost = abs_val(floor_dif)*FLOOR_COST + 1
	// }

	return cost
}

// func Calculate_next_order(calculate_order_ch <-chan map[string]Elev_info, next_order_ch chan<- int, elev_id string)
func Calculate_next_order(calculate_order_ch <-chan map[string]Elev_info, next_order_ch chan<- int, elev_id string) {

	lowest_cost_floor := NO_ORDER
	var lowest_cost int
	var local_cost_this_order int
	var lowest_network_cost int

	for {

		select {

		case online_elevators := <-calculate_order_ch:

			lowest_cost_floor = NO_ORDER
			lowest_cost = N_FLOORS * N_BUTTONS * len(online_elevators) * 100

			for i := 0; i < N_FLOORS; i++ {
				if online_elevators[elev_id].Local_order_matrix[i][INTERNAL_BUTTONS] == 1 && calculate_cost(online_elevators[elev_id].Floor,online_elevators[elev_id].Prev_floor, i, online_elevators[elev_id].Prev_dir, INTERNAL_BUTTONS) < lowest_cost {

					lowest_cost = calculate_cost(online_elevators[elev_id].Floor,online_elevators[elev_id].Prev_floor, i,online_elevators[elev_id].Prev_dir, INTERNAL_BUTTONS)
					lowest_cost_floor = i

				}
			}

			for i := 0; i < N_FLOORS; i++ {

				for j := 0; j < N_BUTTONS-1; j++ {

				
					local_cost_this_order = N_FLOORS * N_BUTTONS * len(online_elevators) * 100
					lowest_network_cost = N_FLOORS * N_BUTTONS * len(online_elevators) * 100

					for order_elevator := range online_elevators {

						if online_elevators[order_elevator].Local_order_matrix[i][j] == 1 {

							for elevator := range online_elevators {

								if elevator != elev_id && calculate_cost(online_elevators[elevator].Floor, online_elevators[elevator].Prev_floor, i,online_elevators[elevator].Prev_dir, j) < lowest_network_cost {

									lowest_network_cost = calculate_cost(online_elevators[elevator].Floor, online_elevators[elevator].Prev_floor, i,online_elevators[elevator].Prev_dir, j)
									
								}

								if elevator == elev_id{

									local_cost_this_order = calculate_cost(online_elevators[elev_id].Floor, online_elevators[elev_id].Prev_floor, i, online_elevators[elevator].Prev_dir, j)

									if order_elevator == elev_id{
										local_cost_this_order = local_cost_this_order-1
									}
								}
							}
						}
					
						if local_cost_this_order < lowest_cost && local_cost_this_order < lowest_network_cost {

							lowest_cost = local_cost_this_order
							lowest_cost_floor = i
							
						}
					}	
				}
			}
			
			Println("FLOOR:  ", lowest_cost_floor, "  COST: ", lowest_cost)
			case next_order_ch <- lowest_cost_floor:
				Sleep(1*Millisecond)

		}
	}
}

func Execute_orders(next_order_ch <-chan int, elev_dir_ch chan <-int){

	target_floor := NO_ORDER
	current_floor := Elev_get_floor_sensor_signal()
	var dir int
	

	for{

		select{

		case target_floor = <- next_order_ch:

			if target_floor != NO_ORDER{
				
				
				if Elev_get_floor_sensor_signal() != LIMBO{				
					current_floor = Elev_get_floor_sensor_signal()
				}

				if target_floor < current_floor{
					Elev_set_motor_direction(DIR_DOWN)
					dir = DIR_DOWN
				}else if target_floor > current_floor{
					Elev_set_motor_direction(DIR_UP)
					dir = DIR_UP
				}else{
					Elev_set_motor_direction(DIR_IDLE)
					dir = DIR_IDLE
				}

				for current_floor != target_floor{

					select{
					case elev_dir_ch <- dir:
					case next_target_floor := <- next_order_ch:
						if next_target_floor != LIMBO{
							target_floor = next_target_floor
						}
					}

					current_floor = Elev_get_floor_sensor_signal()
					
				}

				target_floor = NO_ORDER
				
				Elev_stop_motor()
				Elev_open_door()
				
				dir = DIR_IDLE
			}

		case elev_dir_ch <- dir:

		}
	}
}

func Update_orders_and_lights(system_update_ch <- chan map[string]Elev_info, rem_local_order_ch chan <-[N_FLOORS][N_BUTTONS]int, elev_id string){

	var new_local_order_matrix [N_FLOORS][N_BUTTONS]int

	for{

		select{

		case online_elevators := <- system_update_ch:

			new_local_order_matrix = online_elevators[elev_id].Local_order_matrix

			for i:= 0; i < N_FLOORS; i++{

				for j:= 0; j < N_BUTTONS-1; j++{

					for order_elevator := range online_elevators {

						if online_elevators[order_elevator].Local_order_matrix[i][j] == 1{

							Elev_set_button_lamp(j,i,1)

							for elevator := range online_elevators {

								if online_elevators[elevator].Floor == i{

									if online_elevators[elevator].Dir == DIR_UP && j == EXT_UP_BUTTONS {

										new_local_order_matrix[i][j] = 0
										Elev_set_button_lamp(j,i,0)

									}else if online_elevators[elevator].Dir == DIR_DOWN && j == EXT_DOWN_BUTTONS{

										new_local_order_matrix[i][j] = 0
										Elev_set_button_lamp(j,i,0)

									}else if online_elevators[elevator].Dir == DIR_IDLE && (j == EXT_UP_BUTTONS || j == EXT_DOWN_BUTTONS){

										new_local_order_matrix[i][j] = 0
										Elev_set_button_lamp(j,i,0)

									}
								}
							}
						}
					}
				}
			}

			if online_elevators[elev_id].Local_order_matrix[online_elevators[elev_id].Floor][INTERNAL_BUTTONS] == 1{

				new_local_order_matrix[online_elevators[elev_id].Floor][INTERNAL_BUTTONS] = 0
				Elev_set_button_lamp(INTERNAL_BUTTONS, online_elevators[elev_id].Floor, 0)
				
			}

			for i := 0; i < N_FLOORS; i++{

				if new_local_order_matrix[i][INTERNAL_BUTTONS] == 1{

					Elev_set_button_lamp(INTERNAL_BUTTONS, i, 1)

				}
			}
			if Elev_get_floor_sensor_signal() != LIMBO{
				Elev_set_floor_indicator(Elev_get_floor_sensor_signal())
			}	
			rem_local_order_ch <- new_local_order_matrix
			//Println("new_local_order_matrix:  ", new_local_order_matrix)
		}
	}
}
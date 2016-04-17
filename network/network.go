package network

import . ".././constants"

import (
	. "encoding/json"
	. "fmt"
	"net"
	. "strconv"
)

//Ports and addresses for send/receive:

var local_addr *net.UDPAddr
var broadcast_addr *net.UDPAddr

const broadcast_listen_port = 25001
const local_listen_port = 20020

//Struct with elevator information:

type Elev_info struct {
	Elev_id            string
	Alive_counter      int
	Floor              int
	Prev_floor         int
	Dir                int
	Prev_dir           int
	Local_order_matrix [N_FLOORS][N_BUTTONS]int
}

// Functions for aquiring broadcast and local addresses:

func get_broadcast_addr(broadcast_listen_port int) (err error) {

	broadcast_addr, err = net.ResolveUDPAddr("udp", "129.241.187.255:"+Itoa(broadcast_listen_port))
	if err != nil {
		return err
	}

	Println("Broadcast address: " + broadcast_addr.String())

	return

}

func get_local_addr(local_listen_port int) (err error) {

	temp_conn, err := net.DialUDP("udp", nil, broadcast_addr)
	if err != nil {
		return err
	}

	defer temp_conn.Close()

	temp_addr := temp_conn.LocalAddr()
	local_addr, err = net.ResolveUDPAddr("udp", temp_addr.String())

	local_addr.Port = local_listen_port

	Println("Local address:    ", local_addr.String())

	return

}

//Functions to send and receive UDP packages:

func udp_send(send_chan <-chan Elev_info) {

	conn, _ := net.DialUDP("udp", local_addr, broadcast_addr)

	for {

		select {

		case msg := <-send_chan:

			buf, _ := Marshal(msg)

			conn.Write(buf)

		}
	}
}

func udp_receive(receive_chan chan<- Elev_info) {

	conn, _ := net.ListenUDP("udp", broadcast_addr)

	buf := make([]byte, 256)
	var msg Elev_info

	for {

		n, _, _ := conn.ReadFromUDP(buf)

		Unmarshal(buf[:n], &msg)
		receive_chan <- msg

	}
}

//Function to initialize UDP connections:

func Udp_init(send_chan chan Elev_info, receive_chan chan Elev_info) (local_address string, err bool) {

	err = false

	err_broadcast_addr := get_broadcast_addr(broadcast_listen_port)
	err_local_addr := get_local_addr(local_listen_port)

	go udp_send(send_chan)
	go udp_receive(receive_chan)

	if err_broadcast_addr != nil || err_local_addr != nil {

		err = true

	}

	return local_addr.String(), err

}
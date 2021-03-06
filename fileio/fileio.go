package fileio

import (
	"encoding/json"
	. "fmt"
	"io/ioutil"
	"os"
)

import .".././constants"

const filename = "order_backup.txt"

var read_buf [N_FLOORS][N_BUTTONS]int

//Function to check error on read/write operation:

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//Function to read from local file. Creates the file if it doesn't exist:

func Read() [N_FLOORS][N_BUTTONS]int{

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		Println("File does not exist at location. Making new queue file.")
		os.Create(filename)
	}

	dat, err := ioutil.ReadFile(filename)
	check(err)

	json.Unmarshal(dat, &read_buf)

	return read_buf
}

//Function to write to local file. 

func Write(Local_order_matrix [N_FLOORS][N_BUTTONS]int) {

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		os.Create(filename)
	}

	buf, _ := json.Marshal(Local_order_matrix)

	err := ioutil.WriteFile(filename, buf, 0644)
	check(err)
}
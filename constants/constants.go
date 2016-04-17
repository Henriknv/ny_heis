package constants

//Set number of floors and buttons:

const N_FLOORS = 4
const N_BUTTONS = 3

//Elevator direction:

const DIR_UP = 1
const DIR_DOWN = -1
const DIR_IDLE = 0

//Alive counter to detect disconnected elevators:

const ALIVE_COUNTER = 250 * 2

//Cost constants:

const FLOOR_COST = 10
const TURN_COST = 35

//Other:

const LIMBO = -1
const NO_ORDER = -2

//Buttons on the PLS:

const INTERNAL_BUTTONS = 2
const EXT_UP_BUTTONS = 0
const EXT_DOWN_BUTTONS = 1

//Filename for backup:

const filename = "order_backup.txt"
///////////////////////////////////////////////////////////////
//
// Cellular automaton version of diagonal motion at arbitrary angles
// on the lattice. It requires agents ot be able to count for themselves
// by their internal state machines
//
///////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"time"
	C "Cellibrium"
)

// **********************************************************

const DoF = 20000
var S_TRANSITION_MATRIX = make(map[string]byte)

// ****************************************************************

func main () {

	C.MODEL_NAME = "LatticeIndependentAngles"

	var st [C.Ylim]string
	
	 st[0] = "....................................."
	 st[1] = "....................................."
	 st[2] = "....................................."
	 st[3] = "....................................."
	 st[4] = "....................................."
	 st[5] = "....................................."
	 st[6] = "....................................."
	 st[7] = "....................................."
	 st[8] = "....................................."
	 st[9] = "....................................."
	st[10] = "....................................."
	st[11] = "....................................."
	st[12] = "....................................."
	st[13] = "....................................."
	st[14] = "....................................."
	st[15] = "....................................."
	st[16] = "....................................."
	st[17] = "....................................."
	st[18] = "....................................."
	st[19] = "....................................."
	st[20] = "....................................."
	st[21] = "....................................."
	st[22] = "....................................."
	st[23] = "....................................."
	st[24] = "....................................."
	st[25] = "....................................."
	st[26] = "....................................."
	st[27] = "....................................."  // X
	st[28] = "....................................."  // X
	st[29] = "....................................."  // X
	st[30] = "....................................."
	st[31] = "....................................."
	st[32] = "....................................."
	st[33] = "....................................."
	st[34] = "....................................."
	st[35] = "....................................."
	st[36] = "....................................."
	st[37] = "....................................." // >
	st[38] = "....................................." // >
	st[39] = "....................................." // >
	st[40] = "....................................."
	st[41] = "....................................."
	st[42] = "....................................."
	st[43] = "....................................."
	st[44] = "....................................."
	st[45] = "....................................."
	st[46] = "....................................."
	st[47] = "....................................."
	st[48] = "....................................."
	st[49] = "....................................."
	st[50] = "....................................."
	st[51] = "....................................."
	st[52] = "....................................."
	st[53] = "....................................."
	st[54] = "....................................."
	st[55] = "....................................."
	st[56] = "....................................."
	st[57] = "....................................."
	st[58] = "....................................."
	st[59] = "....................................."
	st[60] = "....................................."
	st[61] = ">...H................................"
	st[62] = "...L1R..............................."
	st[63] = "....2................................"
	st[64] = "....3................................"
	st[65] = "....4................................"
	st[66] = "....T................................"
	st[67] = "....................................."
	st[68] = "....................................."
	st[69] = "....................................."
	st[70] = "....................................."
	st[71] = "....................................."
	st[72] = "....................................."
	st[73] = "....................................."
	st[74] = "....*................................"
	st[75] = "....*................................"

	C.Initialize(st,DoF)
	InitTransitionMatrix()
	EquilGuideRail()
	ShowStates(st,C.MAXTIME,37,76)
}

// ****************************************************************

func InitTransitionMatrix() {

	// This represents anticlockwise oriented cycles, on screen (+Y is down)
	// If there is a fifth element, then it's a conditional prior
	// else it's a wildcard null entry

	// First shift forward phase

	S_TRANSITION_MATRIX["LH.."]  = 'l'
	S_TRANSITION_MATRIX["HR.."]  = 'r'

	S_TRANSITION_MATRIX["R2.."] = 'R'
	S_TRANSITION_MATRIX["2L.."] = 'L'

	S_TRANSITION_MATRIX["R3.."] = 'R'
	S_TRANSITION_MATRIX["3L.."] = 'L'

	S_TRANSITION_MATRIX["R4.."] = 'R'
	S_TRANSITION_MATRIX["4L.."] = 'L'

	S_TRANSITION_MATRIX["H...|."] = 'h'

	// 

	S_TRANSITION_MATRIX["hl1r"] = '1'
	S_TRANSITION_MATRIX["1L2R"] = '2'
	S_TRANSITION_MATRIX["2L3R"] = '3'
	S_TRANSITION_MATRIX["3L4R"] = '4'
	S_TRANSITION_MATRIX["4LTR"] = 'T'

	S_TRANSITION_MATRIX["T...|T"] = '.'

	// Now we shifted everything, and t -> restores previous state
	// but now signal back the count
	// But we don't want to do this, we want to shift right and then return

	S_TRANSITION_MATRIX["TL..|L"] = 'E'
	S_TRANSITION_MATRIX["RT..|R"] = 'W'

	S_TRANSITION_MATRIX["E4L."] = 'E'
	S_TRANSITION_MATRIX["E3L."] = 'E'
	S_TRANSITION_MATRIX["E2l."] = 'E'
	S_TRANSITION_MATRIX["E1..|l"] = 'E'
	S_TRANSITION_MATRIX["Eh.."] = 'e'

	S_TRANSITION_MATRIX["R4W."] = 'W'
	S_TRANSITION_MATRIX["R3W."] = 'W'
	S_TRANSITION_MATRIX["r2W."] = 'W'
	S_TRANSITION_MATRIX[".1W.|r"] = 'W'
	S_TRANSITION_MATRIX["hW.."] = 'w'

	// ** 2 **
	// relabel head
	S_TRANSITION_MATRIX["e1w."]  = '^'

	// Shart shifting spine right

	S_TRANSITION_MATRIX["^W..|w"] = 'F'

	S_TRANSITION_MATRIX["F1W."]  = '1'
	S_TRANSITION_MATRIX["12W."]  = '2'
	S_TRANSITION_MATRIX["23W."]  = '3'
	S_TRANSITION_MATRIX["34W."]  = '4'
	S_TRANSITION_MATRIX["4T..|W"]  = 'T'

	S_TRANSITION_MATRIX["^E21"]  = '.'
	S_TRANSITION_MATRIX[".E32"]  = '.'
	S_TRANSITION_MATRIX[".E43"]  = '.'
	S_TRANSITION_MATRIX[".ET4"]  = '.'
	S_TRANSITION_MATRIX["E.T.|T"]  = '.'

	S_TRANSITION_MATRIX["E..."]  = '.'
	S_TRANSITION_MATRIX["....|E"]  = '.'
	S_TRANSITION_MATRIX["E^..|e"]  = '.'

	S_TRANSITION_MATRIX["F..."]  = '/'
	S_TRANSITION_MATRIX["1/..|."]  = 'L'
	S_TRANSITION_MATRIX["/1..|."]  = 'R'

	S_TRANSITION_MATRIX["/l1r|F"]  = 'H'
	S_TRANSITION_MATRIX["LF.."]  = 'l'
	S_TRANSITION_MATRIX["FR.."]  = 'r'

	S_TRANSITION_MATRIX["H...|/"]  = '.'
}

// ****************************************************************

func EquilGuideRail() {

	for i := 1; i < C.Adim; i++ {
				
		go UpdateAgent_Flow(i)
	}
}

// ****************************************************************

func UpdateAgent_Flow(agent int) {
	
	// Simplify communication for now

	C.CausalIndependence(true)

	for t := 0; t < C.MAXTIME; t++ {

		TransformOrientedState(agent)
		C.CausalIndependence(false)
	}
}

// ****************************************************************

func TransformOrientedState(agent int) {

	// Search anti clockwise for pattern orientation (NB:Y inverted on screen)
	// form a string of the neighbour states clockwise around centre

	var state = make([]byte,C.N+2)
	var vacuum int = 0

	for d := 0; d < C.N; d++ {

		neigh := C.AGENT[agent].Neigh[d]
		state[d] = C.AGENT[neigh].ID

		if state[d] == '.' && C.AGENT[agent].ID == '.' {
			vacuum++
		}
	}

	if vacuum >= C.N {
		return 
	}

	// The centre is my state

	state[C.N] = '|' // conditional
	state[C.N+1] = C.AGENT[agent].ID

	for spin := 0; spin < C.N; spin++ {

		matched,newstate := MatchConfiguration(state)

		if matched {
			C.AGENT[agent].ID = newstate
			return
		}

		state = Rotate(state)
	}
}

// ****************************************************************

func MatchConfiguration(fullstate []byte) (bool,byte) {

	newstate, conditional_promise := S_TRANSITION_MATRIX[string(fullstate)]

	if conditional_promise {
		return true, newstate
	}

	var freestate = make([]byte,4) // truncate

	for i := 0; i < C.N; i++ {
		freestate[i] = fullstate[i]
	}

	genstate, unconditional_promise := S_TRANSITION_MATRIX[string(freestate)]

	if unconditional_promise {
		return true, genstate
	}

	return false, fullstate[C.N]
}

// ****************************************************************

func Rotate(state []byte) []byte {

	var newstate = make([]byte,C.N+2)

	for i := 1; i <= C.N; i++ {

		newstate[i-1] += state[i % C.N]
	}

	newstate[C.N] = state[C.N]
	newstate[C.N+1] = state[C.N+1]
	return newstate
}

// ****************************************************************

func ShowStates(st_rows [C.Ylim]string,tmax,xlim,ylim int) {
	
	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		
		for y := 0; y < ylim; y++ {
			
			for x := 0; x < xlim; x++ {
				
				if !C.Blocked(st_rows,x,y) {

					observable := C.AGENT[C.COORDS[x][y]].ID
					
					fmt.Printf("%3c",observable)
					
				} else {
					fmt.Printf("%3c",' ')
				}
			}
			
			fmt.Println("")
		}

		const base_timescale = 15  // smaller is faster
		const noflicker = 10
		time.Sleep(noflicker * time.Duration(base_timescale) * time.Millisecond) // random noise
	}
}


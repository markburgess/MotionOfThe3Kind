///////////////////////////////////////////////////////////////
//
// Cellular automaton replicator of string agent information, like DNA
// This is a form of semantic translation by cloning
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
	st[35] = ".....xxx............................."
	st[36] = ".....s..............................."
	st[37] = ".....s..............................." // >
	st[38] = ".....s..............................." // >
	st[39] = ".....s..............................." // >
	st[40] = "....................................."
	st[41] = "....................................."
	st[42] = "....................................."
	st[43] = "....................................."
	st[44] = "....................................."
	st[45] = "....................................."
	st[46] = "....................................."
	st[47] = ".......................M............."
	st[48] = ".......................G............."
	st[49] = ".......................A............."
	st[50] = ".......................T............."
	st[51] = ".......................T............."
	st[52] = ".......................A............."
	st[53] = ".......................C............."
	st[54] = ".......................C............."
	st[55] = ".......................A............."
	st[56] = "....................................."
	st[57] = "....................................."
	st[58] = "....................................."
	st[59] = "....................................."
	st[60] = "....................................."
	st[61] = "....................................."
	st[62] = "....................................."
	st[63] = "*...........+........................"
	st[64] = "............1........................"
	st[65] = "............2........................"
	st[66] = "............3........................"
	st[67] = "............4........................"
	st[68] = "............5........................"
	st[69] = "............6........................"
	st[70] = "............7........................"
	st[71] = "............8........................"
	st[72] = "............9........................"
	st[73] = "....................................."
	st[74] = "....................................."
	st[75] = "............*.........X.............."

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

	// 1) Copy number sequence

	S_TRANSITION_MATRIX["1...."] = 'A'
	S_TRANSITION_MATRIX["A+.."]  = 'L'
	S_TRANSITION_MATRIX["+A.."]  = 'R'
	S_TRANSITION_MATRIX["AA.+"]  = 'L'


	S_TRANSITION_MATRIX["LR.."]  = 'X'
	S_TRANSITION_MATRIX["Xr.."]  = 'x'
	S_TRANSITION_MATRIX["RA.."]  = 'r'

	// Duplicate Right

	S_TRANSITION_MATRIX["R1..A"]  = '1'
	S_TRANSITION_MATRIX["r1.."]  = '1'
	S_TRANSITION_MATRIX["12.."]  = '2'
	S_TRANSITION_MATRIX["23.."]  = '3'
	S_TRANSITION_MATRIX["34.."]  = '4'
	S_TRANSITION_MATRIX["45.."]  = '5'
	S_TRANSITION_MATRIX["56.."]  = '6'
	S_TRANSITION_MATRIX["67.."]  = '7'
	S_TRANSITION_MATRIX["78.."]  = '8'
	S_TRANSITION_MATRIX["89.."]  = '9'

	S_TRANSITION_MATRIX["R121"]  = '.'
	S_TRANSITION_MATRIX[".232"]  = '.'
	S_TRANSITION_MATRIX[".343"]  = '.'
	S_TRANSITION_MATRIX[".454"]  = '.'
	S_TRANSITION_MATRIX[".565"]  = '.'
	S_TRANSITION_MATRIX[".676"]  = '.'
	S_TRANSITION_MATRIX[".787"]  = '.'
	S_TRANSITION_MATRIX[".898"]  = '.'
	S_TRANSITION_MATRIX[".9.9"]  = '.'

	S_TRANSITION_MATRIX["+.r."]  = '.'
	S_TRANSITION_MATRIX["A+.."]  = '.'

	// 2) ACGT : G-C, T-U, A-T

	S_TRANSITION_MATRIX["M..."] = 't'
	S_TRANSITION_MATRIX["t+.."]  = 'L'
	S_TRANSITION_MATRIX["+t.."]  = 'R'
	S_TRANSITION_MATRIX["tt.+"]  = 'L'
	S_TRANSITION_MATRIX["tG.."]  = 'C'
	S_TRANSITION_MATRIX["tT.."]  = 'U'
	S_TRANSITION_MATRIX["tA.."]  = 'T'

	S_TRANSITION_MATRIX["CA.."]  = 'T'
	S_TRANSITION_MATRIX["AA.."]  = 'T'
	S_TRANSITION_MATRIX["TA.."]  = 'T'
	S_TRANSITION_MATRIX["GA.."]  = 'T'

	S_TRANSITION_MATRIX["CU.."]  = 'T'
	S_TRANSITION_MATRIX["AU.."]  = 'T'
	S_TRANSITION_MATRIX["TU.."]  = 'T'
	S_TRANSITION_MATRIX["GU.."]  = 'T'

	S_TRANSITION_MATRIX["CT.."]  = 'U'
	S_TRANSITION_MATRIX["AT.."]  = 'U'
	S_TRANSITION_MATRIX["TT.."]  = 'U'
	S_TRANSITION_MATRIX["GT.."]  = 'U'

	S_TRANSITION_MATRIX["CG.."]  = 'C'
	S_TRANSITION_MATRIX["AG.."]  = 'C'
	S_TRANSITION_MATRIX["TG.."]  = 'C'
	S_TRANSITION_MATRIX["GG.."]  = 'C'

	S_TRANSITION_MATRIX["CC.."]  = 'G'
	S_TRANSITION_MATRIX["AC.."]  = 'G'
	S_TRANSITION_MATRIX["TC.."]  = 'G'
	S_TRANSITION_MATRIX["GC.."]  = 'G'

	S_TRANSITION_MATRIX["UA.."]  = 'T'
	S_TRANSITION_MATRIX["UC.."]  = 'G'
	S_TRANSITION_MATRIX["UT.."]  = 'U'
	S_TRANSITION_MATRIX["UG.."]  = 'C'


	// 3) Simple shape filling

	S_TRANSITION_MATRIX["xs..."]  = 's'
	S_TRANSITION_MATRIX["ss..."]  = 's'

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

	var state = make([]byte,C.N+1)
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

	state[C.N] = C.AGENT[agent].ID

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

	var newstate = make([]byte,C.N+1)

	for i := 1; i <= C.N; i++ {

		newstate[i-1] += state[i % C.N]
	}

	newstate[C.N] = state[C.N]
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


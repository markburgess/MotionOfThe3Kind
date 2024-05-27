///////////////////////////////////////////////////////////////
//
// Cellular automaton version of conserved token passing.
// This is the pass the parcel protocol with conservation
// the only wave to ensure true conservation under virtual diffusion.
//
// When we pass PSI with conserved number, we don't get waves in the
// usual way, but something between diffusion and wave interference
// because each direction is treated like a private transaction
// whereas the state is shared between directions otherwise
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
const wrange = 100
const PERIOD = C.WAVELENGTH * wrange

type TrMatrix struct {

	Next byte   // my next state
	Prior byte  // mystate
}

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
	st[47] = "....................................."
	st[48] = "....................................."
	st[49] = "....................................."
	st[50] = "....................................."
	st[51] = ".......................+12..........X"
	st[52] = ".......................1............."
	st[53] = ".......................2............."
	st[54] = ".......................3............."
	st[55] = ".......................4............."
	st[56] = "....................................."
	st[57] = "....................................."
	st[58] = "....................................."
	st[59] = "....................................."
	st[60] = "....................................."
	st[61] = "....................................."
	st[62] = "....................................."
	st[63] = "*.........21+........................"
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

	// Duplicate left

	S_TRANSITION_MATRIX["1L..A"]  = '1'
	S_TRANSITION_MATRIX["21..."]  = '2'
	S_TRANSITION_MATRIX["32..."]  = '3'
	S_TRANSITION_MATRIX["43..."]  = '4'
	S_TRANSITION_MATRIX["54..."]  = '5'
	S_TRANSITION_MATRIX["65..."]  = '6'
	S_TRANSITION_MATRIX["76..."]  = '7'
	S_TRANSITION_MATRIX["87..."]  = '8'
	S_TRANSITION_MATRIX["98..."]  = '9'

	// Unzip


	// Clear

/*	S_TRANSITION_MATRIX["RLL."]  = '.'
	S_TRANSITION_MATRIX["AAL."]  = '.'
	S_TRANSITION_MATRIX["RA.."]  = '+'

	S_TRANSITION_MATRIX["....A"]  = '.'
	S_TRANSITION_MATRIX["....L"]  = '.'
	S_TRANSITION_MATRIX["....R"]  = '.'

	S_TRANSITION_MATRIX["...1A"]  = '+'

	S_TRANSITION_MATRIX["...LA"]  = '.'

	S_TRANSITION_MATRIX["..RA."]  = '+'
	S_TRANSITION_MATRIX["LR1R+"]  = 'X' */

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


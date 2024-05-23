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

var Y_TRANSITION_MATRIX = make(map[string]byte)
var X_TRANSITION_MATRIX = make(map[string]byte)

// ****************************************************************

func main () {

	C.MODEL_NAME = "Translation"

	var st [C.Ylim]string
	
	 st[0] = "*************************************"
	 st[1] = "*...................................*"
	 st[2] = "*...................................*"
	 st[3] = "*...................................*"
	 st[4] = "*...................................*"
	 st[5] = "*...................................*"
	 st[6] = "*...................................*"
	 st[7] = "*...................................*"
	 st[8] = "*...................................*"
	 st[9] = "*...................................*"
	st[10] = "*...................................*"
	st[11] = "*...................................*"
	st[12] = "*...................................*"
	st[13] = "*...................................*"
	st[14] = "*...................................*"
	st[15] = "*...................................*"
	st[16] = "*...................................*"
	st[17] = "*...................................*"
	st[18] = "*...................................*"
	st[19] = "*...................................*"
	st[20] = "*...................................*"
	st[21] = "*...................................*"
	st[22] = "*...................................*"
	st[23] = "*...................................*"
	st[24] = "*...................................*"
	st[25] = "*...................................*"
	st[26] = "*...................................*"
	st[27] = "*...................................*"  // X
	st[28] = "*...................................*"  // X
	st[29] = "*...................................*"  // X
	st[30] = "*...................................*"
	st[31] = "*...................................*"
	st[32] = "*...................................*"
	st[33] = "*...................................*"
	st[34] = "*...................................*"
	st[35] = "*...................................*"
	st[36] = "*........................S..........*"
	st[37] = "*........................+..........*" // >
	st[38] = "*........................N..........*" // >
	st[39] = "*...................................*" // >
	st[40] = "*...................................*"
	st[41] = "*...................................*"
	st[42] = "*...................................*"
	st[43] = "*...................................*"
	st[44] = "*...................................*"
	st[45] = "*...................................*"
	st[46] = "*...................................*"
	st[47] = "*...................................*"
	st[48] = "*...................................*"
	st[49] = "*...................................*"
	st[50] = "*...................................*"
	st[51] = "*...................................*"
	st[52] = "*..W+E..............................*"
	st[53] = "*...................................*"
	st[54] = "*...................................*"
	st[55] = "*...................................*"
	st[56] = "*........E+W........................*"
	st[57] = "*...................................*"
	st[58] = "*...................................*"
	st[59] = "*...................................*"
	st[60] = "*...................................*"
	st[61] = "*...................................*"
	st[62] = "*...................................*"
	st[63] = "*...................................*"
	st[64] = "*...................................*"
	st[65] = "*...................................*"
	st[66] = "*...................................*"
	st[67] = "*...................................*"
	st[68] = "*...................................*"
	st[69] = "*......N............................*"
	st[70] = "*......+............................*"
	st[71] = "*......S............................*"
	st[72] = "*...................................*"
	st[73] = "*...................................*"
	st[74] = "*...................................*"
	st[75] = "*************************************"

	C.Initialize(st,DoF)
	InitTransitionMatrix()
	EquilGuideRail()
	ShowStates(st,C.MAXTIME,37,76)
}

// ****************************************************************

func InitTransitionMatrix() {

	// This represents oriented cycles

	X_TRANSITION_MATRIX["..E"] = 'E'
	X_TRANSITION_MATRIX["EE+"] = '_'
	X_TRANSITION_MATRIX["E_+"] = '+'
	X_TRANSITION_MATRIX["++W"] = 'W'
	X_TRANSITION_MATRIX["O_W"] = 'W'
	X_TRANSITION_MATRIX["WW."] = '.'

	// orth

	Y_TRANSITION_MATRIX["..N"] = 'N'
	Y_TRANSITION_MATRIX["NN+"] = '_'
	Y_TRANSITION_MATRIX["N_+"] = '+'
	Y_TRANSITION_MATRIX["++S"] = 'S'
	Y_TRANSITION_MATRIX["O_S"] = 'S'
	Y_TRANSITION_MATRIX["SS."] = '.'
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

		// An FSM propagator
	
		// if I am downstream, prepare to receive
		// we must have a lattice to preserve semantic continuity of direction
		// this is borne out my the existence of dipoles

		var d,dbar int

		d = C.EAST
		dbar = (d + C.DIMENSION) % (2*C.DIMENSION)

		e,w := InferPolarity(agent,d,dbar)

		C.AGENT[agent].ID = TransformState(d,e,C.AGENT[agent].ID,w)

		C.CausalIndependence(false)

		d = C.NORTH
		dbar = (d + C.DIMENSION) % (2*C.DIMENSION)

		n,s := InferPolarity(agent,d,dbar)

		C.AGENT[agent].ID = TransformState(d,n,C.AGENT[agent].ID,s)

		C.CausalIndependence(false)
		C.CausalIndependence(false)
		
	}
}

// ****************************************************************

func InferPolarity(agent,d,dbar int) (byte,byte) {

	nf := C.AGENT[agent].Neigh[d]
	nb := C.AGENT[agent].Neigh[dbar]

	m := C.AGENT[agent].ID
	f := C.AGENT[nf].ID
	b := C.AGENT[nb].ID

	state := fmt.Sprintf("%c%c%c",f,m,b)

	var exists bool

	switch d {

	case C.EAST: _, exists = X_TRANSITION_MATRIX[state]
	case C.NORTH: _, exists = Y_TRANSITION_MATRIX[state]
	}

	if exists {
		return f,b
	}

	state = fmt.Sprintf("%c%c%c",b,m,f)

	switch d {

	case C.EAST: _, exists = X_TRANSITION_MATRIX[state]
	case C.NORTH: _, exists = Y_TRANSITION_MATRIX[state]
	}

	if exists {
		return b,f
	}

	return 'x','x'
}

// ****************************************************************

func TransformState(direction int,fwd,me,bwd byte) byte {

	state := fmt.Sprintf("%c%c%c",fwd,me,bwd)

	var newstate byte
	var exists bool

	switch direction {

	case C.EAST: newstate, exists = X_TRANSITION_MATRIX[state]
	case C.NORTH: newstate, exists = Y_TRANSITION_MATRIX[state]
	}

	if exists {
		return newstate
	} else {
		return me
	}
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


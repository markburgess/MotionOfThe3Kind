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

// ****************************************************************

func main () {

	C.MODEL_NAME = "OneDimension"

	var st [C.Ylim]string
	st[0] = "********************************************************"
	st[1] = "..........................>>>..........................|"  // X
	st[2] = "********************************************************"

	// Keep the data structures for agents global too for convenience

	C.Initialize(st,DoF)
	EquilGuideRail()
	ShowLinear(C.MAXTIME)
}

// ****************************************************************

func EquilGuideRail() {

	for i := 1; i < C.Adim; i++ {
				
		go UpdateAgent_Flow(i)
	}
}

// ****************************************************************

func UpdateAgent_Flow(agent int) {

	const ew_only = 2

	for direction := ew_only; direction < C.N; direction++ {

		neighbour := C.AGENT[agent].Neigh[direction]
		
		if neighbour != 0 {
			var breaker C.Message
			breaker.Value = C.AGENT[agent].Psi
			breaker.Phase = C.TICK
			C.CHANNEL[agent][neighbour] = breaker
		}
	}

	const PsiQuant = 1

	C.CausalIndependence(true)

	for t := 0; t < C.MAXTIME; t++ {
		
		// Every pair of agents has a private directional channel that's not overwritten by anyone else
		// Messages persist until they are read and cannot unseen

		for direction := 0; direction < C.N; direction++ {
			
			var send,recv C.Message
			
			neighbour := C.AGENT[agent].Neigh[direction]
			
			if neighbour == 0 {
				continue // wall signal
			}

			// We need to wait for a positive signal indicating a new transfer to avoid double/empty reading

			recv = C.AcceptFromChannel(neighbour,agent)

			// ****************** PROCESS *********************

			switch recv.Phase {
				
			case C.TICK:
				C.AGENT[agent].V[direction] = recv.Value
				C.AGENT[agent].Psi += EvolvePsi(agent,direction)
				send.Value = C.AGENT[agent].Psi
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
			}
		}
	}
}


// ****************************************************************
// Now we must be able to handle negative amounts too: debt, else just a diffusion problem
// ****************************************************************

func EvolvePsi(agent,direction int) float64 {

	const mass = 15.0
	const dt = 0.01

	grad := C.AGENT[agent].V[direction] - C.AGENT[agent].Psi

	deltaTheta := grad / mass *  dt // small increments

	DeltaPsi := C.AGENT[agent].Theta * dt

	C.AGENT[agent].Theta += deltaTheta

	return DeltaPsi

}

// ****************************************************************

func ShowLinear(tmax int) {

	var count float64 
	const height = 60

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		
		for y := height/2+1; y > -height/2-1; y-- {

			count = 0

			for x := 0; x < C.Xlim; x++ {
				
				observable := float64(C.AGENT[C.COORDS[x][1]].Psi)
				count += observable
				
				h := int(observable/float64(2*DoF) * float64(height)+0.5)

				if (y > 0 && y <= h) {
					//fmt.Printf("%6d",int(observable))
					fmt.Printf("%3s","+")
				} else if (y < 0 && y >= h) {
					//fmt.Printf("%6d",int(observable))
					fmt.Printf("%3s","-")
				} else {
					//fmt.Printf("%6s",".")
					fmt.Printf("%3s"," ")
				}
			}
			fmt.Println(" X")
		}

	fmt.Println("\n\nTOTAL =",count)
	
	const base_timescale = 15  // smaller is faster
	const noflicker = 10
	time.Sleep(noflicker * time.Duration(base_timescale) * time.Millisecond) // random noise

	}
}
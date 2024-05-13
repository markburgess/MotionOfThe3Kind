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
	st[0] = "*****************************************"
	st[1] = ">.......................................|"  // X
	st[2] = "*****************************************"

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
	
	// Start with an  unconditional promise to break the deadlock symmetry

	for direction := 0; direction < C.N; direction++ {

		neighbour := C.AGENT[agent].Neigh[direction]
		
		if neighbour != 0 {
			var breaker C.Message
			breaker.Value = C.AGENT[agent].Psi
			breaker.Phase = C.TICK
			C.CHANNEL[agent][neighbour] = breaker
		}
	}

	C.CausalIndependence(true)

	for t := 0; t < C.MAXTIME; t++ {
		
		for direction := 0; direction < C.N; direction++ {
			
			var send,recv C.Message
			
			neighbour := C.AGENT[agent].Neigh[direction]
			
			if neighbour == 0 {
				continue // wall signal
			}

			recv = C.AcceptFromChannel(neighbour,agent)

			// ****************** PROCESS *********************

			switch recv.Phase {
				
			case C.TICK: // me - this is foreach direction...

				C.AGENT[agent].V[direction] = recv.Value
				send.Phase = C.TAKE
				// Because there is private communication with each agent, this is now inside the loop
				send.Value = EvolveAndOfferDeltaPsi(agent,direction)
				C.AGENT[agent].Offer[direction] = send.Value
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue

			case C.TAKE: // YOU
				transfer_offer := recv.Value
				send.Value = transfer_offer
				C.AGENT[agent].Accept[direction] = transfer_offer // reserve amount
				C.AGENT[agent].Psi -= transfer_offer
				send.Phase = C.TACK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue

			case C.TACK: // me
				transfer_offer := recv.Value

				if C.AGENT[agent].Offer[direction] == transfer_offer {
					C.AGENT[agent].Psi += transfer_offer
					send.Value = transfer_offer
				} else {
					send.Value = C.NOTACCEPT
				}
				C.AGENT[agent].Accept[direction] = 0
				C.AGENT[agent].Offer[direction] = 0
				send.Phase = C.TOCK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
				
			case C.TOCK: // YOU - initiate a change / Xfer

				if recv.Value == C.NOTACCEPT {
					// Move the reserved amount back
					C.AGENT[agent].Psi += C.AGENT[agent].Accept[direction]
				}

				C.AGENT[agent].Accept[direction] = 0  // clear reservation, consider sent
				C.AGENT[agent].Offer[direction] = 0
				send.Value = C.AGENT[agent].Psi
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
			}
		}
	}
}

// ****************************************************************
// Now we must be able to handle negative amounts too: debt, else just a diffusion problem
// ****************************************************************

func EvolveAndOfferDeltaPsi(agent,direction int) int {

	const mass = 111
	const coupling = 11

	const dt = 1

	// Because this is private for each neighbour, the
	// Laplacian is just the gradient for each individual direction

	d2 := C.AGENT[agent].V[direction] - C.AGENT[agent].Psi

	// This is negative when Psi is higher than neighbours
	// Increase inv_velocity theta to increase wavelength

	dtheta := d2 / mass

	// Stability, reduce velocity more quicky (leads to quicker smaller oscillations)

	C.AGENT[agent].Theta += dtheta * dt % PERIOD  // update velocity

	deltaPsi := C.AGENT[agent].Theta * dt // displacement = velocity x time

	// C.AGENT[agent].Psi += deltaPsi deferred to xfer
	return deltaPsi / coupling

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
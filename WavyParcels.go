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
	"os"
	C "Cellibrium"
)

// **********************************************************

const DoF = 20000
const wrange = 100
const PERIOD = C.WAVELENGTH * wrange

// ****************************************************************

func main () {

	C.MODEL_NAME = "WavyParcels"

	var st [C.Ylim]string

	switch os.Args[1] {
	case "ds" : C.DOUBLE_SLIT = true
	case "ss" : C.DOUBLE_SLIT = false
	}
	
	 st[0] = "*************************************"
	 st[1] = "*........*..........................|"
	 st[2] = "*........*..........................|"
	 st[3] = "*........*..........................|"
	 st[4] = "*........*..........................|"
	 st[5] = "*........*..........................|"
	 st[6] = "*........*..........................|"
	 st[7] = "*........*..........................|"
	 st[8] = "*........*..........................|"
	 st[9] = "*........*..........................|"
	st[10] = "*........*..........................|"
	st[11] = "*........*..........................|"
	st[12] = "*........*..........................|"
	st[13] = "*........*..........................|"
	st[14] = "*........*..........................|"
	st[15] = "*........*..........................|"
	st[16] = "*........*..........................|"
	st[17] = "*........*..........................|"
	st[18] = "*........*..........................|"
	st[19] = "*........*..........................|"
	st[20] = "*........*..........................|"
	st[21] = "*........*..........................|"
	st[22] = "*........*..........................|"
	st[23] = "*........*..........................|"
	st[24] = "*........*..........................|"
	st[25] = "*........*..........................|"
	st[26] = "*........*..........................|"
	st[27] = "*........>..........................|"  // X
	st[28] = "*........>..........................|"  // X
	st[29] = "*........>..........................|"  // X
	st[30] = "*........*..........................|"
	st[31] = "*........*..........................|"
	st[32] = "*........*..........................|"
	st[33] = "*........*..........................|"
	st[34] = "*........*..........................|"
	st[35] = "*........*..........................|"
	st[36] = "*........*..........................|"
	st[37] = "*........*..........................|" // >
	st[38] = "*........*..........................|" // >
	st[39] = "*........*..........................|" // >
	st[40] = "*........*..........................|"
	st[41] = "*........*..........................|"
	st[42] = "*........*..........................|"
	st[43] = "*........*..........................|"
	st[44] = "*........*..........................|"
	st[45] = "*........*..........................|"
	st[46] = "*........*..........................|"
	st[47] = "*........>..........................|"
	st[48] = "*........>..........................|"
	st[49] = "*........>..........................|"
	st[50] = "*........*..........................|"
	st[51] = "*........*..........................|"
	st[52] = "*........*..........................|"
	st[53] = "*........*..........................|"
	st[54] = "*........*..........................|"
	st[55] = "*........*..........................|"
	st[56] = "*........*..........................|"
	st[57] = "*........*..........................|"
	st[58] = "*........*..........................|"
	st[59] = "*........*..........................|"
	st[60] = "*........*..........................|"
	st[61] = "*........*..........................|"
	st[62] = "*........*..........................|"
	st[63] = "*........*..........................|"
	st[64] = "*........*..........................|"
	st[65] = "*........*..........................|"
	st[66] = "*........*..........................|"
	st[67] = "*........*..........................|"
	st[68] = "*........*..........................|"
	st[69] = "*........*..........................|"
	st[70] = "*........*..........................|"
	st[71] = "*........*..........................|"
	st[72] = "*........*..........................|"
	st[73] = "*........*..........................|"
	st[74] = "*........*..........................|"
	st[75] = "*************************************"
	// Keep the data structures for agents global too for convenience


	C.Initialize(st,DoF)
	C.ShowState(st,1,37,76,"num")
	EquilGuideRail()
	//C.ShowState(st,C.MAXTIME,37,76,"+")
	C.ShowAffinity(st,C.MAXTIME,37,76)
	//C.ShowPhase(st,C.MAXTIME,37,76)
	//go C.MovingPromise()
	//C.ShowPosition(st,C.MAXTIME,37,76)
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

func EvolveAndOfferDeltaPsi(agent,direction int) float64 {

	const mass = 111
	const coupling = 11

	const dt = 1.0

	// Because this is private for each neighbour, the
	// Laplacian is just the gradient for each individual direction

	d2 := C.AGENT[agent].V[direction] - C.AGENT[agent].Psi

	// This is negative when Psi is higher than neighbours
	// Increase inv_velocity theta to increase wavelength

	dtheta := d2 / mass

	// Stability, reduce velocity more quicky (leads to quicker smaller oscillations)

	C.AGENT[agent].Theta += float64(int(dtheta * dt+0.5) % PERIOD)

	deltaPsi := C.AGENT[agent].Theta * dt // displacement = velocity x time

	// C.AGENT[agent].Psi += deltaPsi deferred to xfer
	return deltaPsi / coupling

}


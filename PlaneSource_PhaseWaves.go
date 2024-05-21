///////////////////////////////////////////////////////////////
//
// Cellular automaton version Linear Schrodinger-like waves
// This is a single time-deriv phase based Schrodinger-like
// with no intentional conservation, just transmitted vector
// Instead of an exterior restoring force, we now have an interior
// memory of intermediate velocity as a phase shift
//
// In this version the amplitude source starts from the left
// as a single points source which has to pass through the slits
//
// Changing arg ss to ds opens the XXX sluices ...
//
///////////////////////////////////////////////////////////////

package main

import (
//	"fmt"
	"os"
	C "Cellibrium"
)

// **********************************************************

const DoF = 10000
const wrange = 10000
const PERIOD = C.WAVELENGTH * wrange

// ****************************************************************

func main () {

	C.MODEL_NAME = "LongSourceSlits"

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
	st[27] = "*........<..........................|"  // X
	st[28] = "*........<..........................|"  // X
	st[29] = "*........<..........................|"  // X
	st[30] = "*........*..........................|"
	st[31] = "*........*..........................|"
	st[32] = "*........*..........................|"
	st[33] = "*........*..........................|"
	st[34] = "*........*..........................|"
	st[35] = "*........*..........................|"
	st[36] = "*........*..........................|"
	st[37] = "*........*..........................|"
	st[38] = "*........*..........................|"
	st[39] = "*........*..........................|"
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

	C.ShowState(st,1,37,76,"+")
	EquilGuideRail()
	//C.ShowState(st,C.MAXTIME,37,76,"+")
	C.ShowAffinity(st,C.MAXTIME,37,76)
	//go C.MovingPromise()
	//C.ShowPhase(st,C.MAXTIME,37,76)
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
			breaker.Angle = C.AGENT[agent].Theta
			breaker.Phase = C.TICK
			C.CHANNEL[agent][neighbour] = breaker
		}
	}

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

			recv = C.AcceptFromChannel(neighbour,agent)

			// ****************** PROCESS *********************
			// We put directional updates per direction instead of
			// for all at once (grad rather than Laplacian) to get
			// more accurate updates..

			C.AGENT[agent] = EvolvePsi(C.AGENT[agent],direction)

			switch recv.Phase {
				
			case C.TICK:
				C.AGENT[agent].V[direction] = recv.Value
				send.Value = C.AGENT[agent].Psi
				send.Angle = C.AGENT[agent].Theta
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
			}
		}

		// If we put a full Laplacian here, we get poor results
	}
}


// ****************************************************************

func EvolvePsi(agent C.STAgent,direction int) C.STAgent { // Laplacian

	/* The challenge is to stop Psi from growing in amplitude so that differences 
           no longer matter and the waves eventually stop propagating. It' s very
           hard to do this with small integer arithmetic .. which suggests that the smoothness
           of quantum phenomena suggest that there is plenty of room at the bottom for large numbers. */

	agent.Theta += dTheta(agent,direction)  // float64(int(dTheta(agent)+0.5) % PERIOD)
	agent.Psi += dPsi(agent)

	return agent
}

// ******************************************************************

func dTheta(agent C.STAgent,di int) float64 { // Laplacian

	var   d2 float64 = 0
	const dt = 0.01
	const mass = 2.0

	// Velocity = laplaciant gradient

	d2 += agent.V[di] - agent.Psi

	// This is negative when Psi is higher than neighbours

	dtheta := dt * d2 / float64(C.N * mass)

	return dtheta
}

// ******************************************************************

func dPsi(agent C.STAgent) float64 { // Laplacian

	const dt = 0.01

	deltaPsi := agent.Theta * dt

	return deltaPsi
}



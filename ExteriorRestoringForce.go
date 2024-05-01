///////////////////////////////////////////////////////////////
//
// Cellular automaton version of classical externalized wave dynamics
// This is a Laplacian derived state dynamic, pass the real gradients
// and no intentional conservation, just transmitted "string tension"
// providing a stigmergic restoring force...
//
///////////////////////////////////////////////////////////////

package main

import (
	"os"
	C "Cellibrium"
)

// **********************************************************

const momentum = 1
const psi = 0
const MODE = psi
const DoF = 10000
const MAXTIME = 100000

// ****************************************************************

func main () {

	C.MODEL_NAME = "ExtRestForce"

	switch os.Args[1] {
	case "ds" : C.DOUBLE_SLIT = true
	case "ss" : C.DOUBLE_SLIT = false
	}

	var st [C.Ylim]string

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
	st[27] = "*........X..........................|"  // X
	st[28] = "*........X..........................|"  // X
	st[29] = "*........X..........................|"  // X
	st[30] = "*........*..........................|"
	st[31] = "*........*..........................|"
	st[32] = "*........*..........................|"
	st[33] = "*........*..........................|"
	st[34] = "*........*..........................|"
	st[35] = "*........*..........................|"
	st[36] = "*........*..........................|"
	st[37] = ">........*..........................|"
	st[38] = ">........*..........................|"
	st[39] = ">........*..........................|"
	st[40] = "*........*..........................|"
	st[41] = "*........*..........................|"
	st[42] = "*........*..........................|"
	st[43] = "*........*..........................|"
	st[44] = "*........*..........................|"
	st[45] = "*........*..........................|"
	st[46] = "*........*..........................|"
	st[47] = "*...................................|"
	st[48] = "*...................................|"
	st[49] = "*...................................|"
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

	C.ShowState(st,1,23,65)
	EquilGuideRail()
	C.ShowState(st,MAXTIME,23,65)
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

	const PsiQuant = 1

	C.CausalIndependence(true)

	for t := 0; t < MAXTIME; t++ {
		
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
				send.Value = C.AGENT[agent].Psi
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
			}
		}
		
		// Now we have updated neighbour Psi[N] info to reevaluate

		C.AGENT[agent] = EvolvePsi(C.AGENT[agent])
		
	}
}

// ****************************************************************

func EvolvePsi(agent C.STAgent) C.STAgent { // Laplacian

	agent.PsiDot += dPsiDot(agent)
	agent.Psi += dPsi(agent)

	return agent
}

// ****************************************************************

func dPsiDot(agent C.STAgent) int { // Laplacian

	var deltaPsiDot int = 0
	const dt = 1
	const velocity = 9

	// Velocity = laplaciant gradient

	for di := 0; di < C.N; di++ {

		deltaPsiDot += agent.V[di] - agent.Psi
	}

	// This is negative when Psi is higher than neighbours
	dv := dt * deltaPsiDot / (C.N * velocity)

	return dv
}

// ****************************************************************

func dPsi(agent C.STAgent) int { // Laplacian

	var deltaPsi int = 0
	const dt = 1
	const velocity = 10

	deltaPsi = agent.PsiDot * dt
	dpsi := deltaPsi / velocity

	return dpsi
}

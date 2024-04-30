///////////////////////////////////////////////////////////////
//
// Cellular automaton version of spherical waves emanating from
// a source
// with no intentional conservation, just transmitted vector
// This one with smoother waveforms
//
///////////////////////////////////////////////////////////////

package main

import (
	C "Cellibrium"
)

// **********************************************************

const DoF = 11000
const MAXTIME = 100000

// ****************************************************************

func main () {

	C.MODEL_NAME = "Waves"

	var st [C.Ylim]string
	
	 st[0] = "*************************************"
	 st[1] = "*...................................|"
	 st[2] = "*...................................|"
	 st[3] = "*...................................|"
	 st[4] = "*...................................|"
	 st[5] = "*...................................|"
	 st[6] = "*...................................|"
	 st[7] = "*...................................|"
	 st[8] = "*...................................|"
	 st[9] = "*...................................|"
	st[10] = "*...................................|"
	st[11] = "*...................................|"
	st[12] = "*...................................|"
	st[13] = "*...................................|"
	st[14] = "*...................................|"
	st[15] = "*...................................|"
	st[16] = "*...................................|"
	st[17] = "*...................................|"
	st[18] = "*...................................|"
	st[19] = "*...................................|"
	st[20] = "*...................................|"
	st[21] = "*...................................|"
	st[22] = "*...................................|"
	st[23] = "*...................................|"
	st[24] = "*...................................|"
	st[25] = "*...................................|"
	st[26] = "*...................................|"
	st[27] = "*...................................|"
	st[28] = "*...................................|"
	st[29] = "*...................................|"
	st[30] = "*...................................|"
	st[31] = "*...................................|"
	st[32] = "*...............>>..................|"
	st[33] = "*...............>>..................|"
	st[34] = "*...................................|"
	st[35] = "*...................................|"
	st[36] = "*...................................|"
	st[37] = "*...................................|"
	st[38] = "*...................................|"
	st[39] = "*...................................|"
	st[40] = "*...................................|"
	st[41] = "*...................................|"
	st[42] = "*...................................|"
	st[43] = "*...................................|"
	st[44] = "*...................................|"
	st[45] = "*...................................|"
	st[46] = "*...................................|"
	st[47] = "*...................................|"
	st[48] = "*...................................|"
	st[49] = "*...................................|"
	st[50] = "*...................................|"
	st[51] = "*...................................|"
	st[52] = "*...................................|"
	st[53] = "*...................................|"
	st[54] = "*...................................|"
	st[55] = "*...................................|"
	st[56] = "*...................................|"
	st[57] = "*...................................|"
	st[58] = "*...................................|"
	st[59] = "*...................................|"
	st[60] = "*...................................|"
	st[61] = "*...................................|"
	st[62] = "*...................................|"
	st[63] = "*...................................|"
	st[64] = "*...................................|"
	st[65] = "*************************************"

	// Keep the data structures for agents global too for convenience

	C.Initialize(st,DoF)

	C.ShowState(st,1,37,66)
	EquilGuideRail()
	C.ShowState(st,MAXTIME,37,66)
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
		
		// Now we have updated messages from our neighbours for their Psi[N]

		C.AGENT[agent] = EvolvePsiTypeII(C.AGENT[agent])
	}
}

// ****************************************************************

func EvolvePsiTypeI(agent C.STAgent) C.STAgent { // Laplacian

	/* The challenge is to stop Psi from growing in amplitude so that differences 
           no longer matter and the waves eventually stop propagating. It' s very
hard to do this with small integer arithmetic .. which suggests that the smoothness
of quantum phenomena suggest that there is plenty of room at the bottom for large numbers.
The spin case converges over about 100 iterations with a simple two state model, so for
waves with interference */

	const affinity = 10
	const v2 = 3  // odd number 3,5,7
	var d2 int = 0
	var newagent C.STAgent = agent

	for di := 0; di < C.N; di++ {		

		d2 += (agent.V[di] - agent.Psi)
	}

	// To shorten the wavelength increase v2 - even/odd numbers play a role due to the discrete scale

	newtheta := (agent.Theta + C.WAVELENGTH/4) % C.WAVELENGTH; // (agent.Theta + d2/v2) % WAVELENGTH;
	dpsi := -C.WAVE[newtheta] * d2/v2
	newagent.Psi = agent.Psi + dpsi
	newagent.Theta = newtheta 

	return newagent
}

// ****************************************************************

func EvolvePsiTypeII(agent C.STAgent) C.STAgent { // Laplacian

	const affinity = 10
	const v2 = 5  // odd number
	var d2 int = 0
	var newagent C.STAgent = agent

	for di := 0; di < C.N; di++ {		

		d2 += (agent.V[di] - agent.Psi)
	}

	// To shorten the wavelength increase v2 - even/odd numbers play a role due to the discrete scale

	newtheta := (agent.Theta + d2/v2) % C.WAVELENGTH;

	for offset := C.WAVELENGTH; newtheta < 0; offset += C.WAVELENGTH {

		newtheta = (agent.Theta + offset + d2) % C.WAVELENGTH;
	}

	drho := C.WAVE[newtheta] * d2/v2 / C.N
	newagent.Psi = agent.Psi + drho
	newagent.Theta = newtheta 

	return newagent
}

///////////////////////////////////////////////////////////////
//
// Can we maintain particle cluster coherence with a centre of mass
// guided by a psi wave? Without a directional affinity for each
// subagent, the mass simple drifts to the edges of the wave and
// loses coherence, even if we break the psi symmetry relative to M
//
// This initial condition is far enough from the wave to be stable
// as long as each cell doesn't overselect several directions instead of one
//
///////////////////////////////////////////////////////////////

package main

import (
//	"fmt"
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

	 st[0] = "....................................."
	 st[1] = "....................................|"
	 st[2] = "....................................|"
	 st[3] = "....................................|"
	 st[4] = "....................................|"
	 st[5] = "....................................|"
	 st[6] = "....................................|"
	 st[7] = "....................................|"
	 st[8] = "....................................|"
	 st[9] = "....................................|"
	st[10] = "....................................|"
	st[11] = "....................................|"
	st[12] = "....................................|"
	st[13] = "....................................|"
	st[14] = "....................................|"
	st[15] = "....................................|"
	st[16] = "....................................|"
	st[17] = "....................................|"
	st[18] = "....................................|"
	st[19] = "....................................|"
	st[20] = "....................................|"
	st[21] = "....................................|"
	st[22] = "....................................|"
	st[23] = "....................................|"
	st[24] = "....................................|"
	st[25] = "....................................|"
	st[26] = "....................................|"
	st[27] = "....................................|"  // X
	st[28] = "....................................|"  // X
	st[29] = "....................................|"  // X
	st[30] = "....................................|"
	st[31] = "....................................|"
	st[32] = "....................................|"
	st[33] = "....................................|"
	st[34] = "....................................|"
	st[35] = "....................................|"
	st[36] = "....................................|"
	st[37] = "....................................|"
	st[38] = "....................................|"
	st[39] = "....................................|"
	st[40] = "....................................|"
	st[41] = "....................................|"
	st[42] = "....................................|"
	st[43] = "....................................|"
	st[44] = ".......m............................|"
	st[45] = ">.....mmm...........................|"
	st[46] = ">....mmmmm .........................|"
	st[47] = ">...mmmmmmm.........................|"
	st[48] = ">...mmmmmmm.........................|"
	st[49] = ">....mmmmm..........................|"
	st[50] = ">.....mmm...........................|"
	st[51] = ".......m............................|"
	st[52] = "....................................|"
	st[53] = "....................................|"
	st[54] = "....................................|"
	st[55] = "....................................|"
	st[56] = "....................................|"
	st[57] = "....................................|"
	st[58] = "....................................|"
	st[59] = "....................................|"
	st[60] = "....................................|"
	st[61] = "....................................|"
	st[62] = "....................................|"
	st[63] = "....................................|"
	st[64] = "....................................|"
	st[65] = "....................................|"
	st[66] = "....................................|"
	st[67] = "....................................|"
	st[68] = "....................................|"
	st[69] = "....................................|"
	st[70] = "....................................|"
	st[71] = "....................................|"
	st[72] = "....................................|"
	st[73] = "....................................|"
	st[74] = "....................................|"
	st[75] = "....................................."

	// Keep the data structures for agents global too for convenience

	C.Initialize(st,DoF)

	C.ShowState(st,1,37,76,"+")
	EquilGuideRail()
	C.ShowState(st,C.MAXTIME,37,76,"+")
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

		C.AGENT[agent].P[direction] = -1
		
		if neighbour != 0 {
			var breaker C.Message
			breaker.Value = C.AGENT[agent].Psi
			breaker.Phase = C.TICK
			C.CHANNEL[agent][neighbour] = breaker
		}
	}

	const PsiQuant = 1
	C.AGENT[agent].Moment = -1

	C.CausalIndependence(true)

	for t := 0; t < C.MAXTIME; t++ {
		
		// Every pair of agents has a private directional channel that's not overwritten by anyone else
		// Messages persist until they are read and cannot unseen

		for d := 0; d < C.N; d++ {

			var send,recv C.Message
			
			neighbour := C.AGENT[agent].Neigh[d]

			dbar := (d + C.N/2) % C.N
			
			if neighbour == 0 {
				continue // wall signal
			}

			// We need to wait for a positive signal indicating a new transfer to avoid double/empty reading

			recv = C.AcceptFromChannel(neighbour,agent)

			// ****************** PROCESS *********************

			switch recv.Phase {
				
			case C.TICK: // me

				C.AGENT[agent].V[d] = recv.Value  // the value here is psi
				C.AGENT[agent].M[d] = recv.MassID
				C.AGENT[agent].P[d] = recv.Moment

				C.AGENT[agent] = EvolveAgents(C.AGENT[agent],d)

				// In this phase we can choose to make an offer to accept
				// a neighbouring massID, we have to have received an update first

				if AcceptingMass(C.AGENT[agent],d,dbar) > 0 {
					send.Phase = C.TAKE
					send.Value = 1
					C.AGENT[agent].Offer[d] = send.Value
					C.AGENT[agent].Moment = dbar
					C.ConditionalChannelOffer(agent,neighbour,send)
				} else {
					send.MassID = C.AGENT[agent].MassID
					send.Moment = C.AGENT[agent].Moment
					send.Value = C.AGENT[agent].Psi
					send.Angle = C.AGENT[agent].Theta
					send.Phase = C.TICK
					C.ConditionalChannelOffer(agent,neighbour,send)
				}

				continue

			case C.TAKE: // YOU
				transfer_offer := recv.Value  // This is massID now
				send.Value = transfer_offer
				send.Phase = C.TACK
				C.AGENT[agent].MassID = 0     // assume you'll take it
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue

			case C.TACK: // ME

				transfer_offer := recv.Value
				C.AGENT[agent].MassID = transfer_offer
				send.Value = transfer_offer
				send.Phase = C.TOCK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
				
			case C.TOCK: // YOU - initiate a change / Xfer

				C.AGENT[agent].MassID = 0
				send.MassID = C.AGENT[agent].MassID
				send.Value = C.AGENT[agent].Psi
				send.Moment = C.AGENT[agent].Moment
				send.Angle = C.AGENT[agent].Theta
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
			}
		}
	}
}

// ****************************************************************

func AcceptingMass(agent C.STAgent,d,dbar int) int {

	// We look to accept some mass from d if dbar looks empty

	if agent.Psi * agent.Psi < 0.1 {
		return 0
	}

	// find max gradient behind us
	var max float64 = 0
	var dbarmax int

	for d := 0; d < C.N; d++ {

		affinity := agent.V[d] * agent.V[d]
		dbar := (d + C.N/2) % C.N

		if affinity > max {
			max = affinity
			dbarmax = dbar
		}
	}

	// If the bow wave is focussed too close to the mass, it will diffuse the mass
	agent.Moment = dbarmax

	// select the direction of motion - conservation of initial momentum

	if (agent.M[d] > 0) && (agent.MassID == 0) && (dbar == agent.Moment) {
		return 1
	}

	return 0
}

// ****************************************************************

func WillingToPassToRecipient(agent,d int) bool {

	dbar := (d + C.N/2) % C.N

	if C.AGENT[agent].MassID > 0 && C.AGENT[agent].Moment == dbar {
		return true
	}

	return false
}

// ****************************************************************

func EvolveAgents(agent C.STAgent,direction int) C.STAgent { // Laplacian

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



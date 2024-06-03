///////////////////////////////////////////////////////////////
//
// Can we maintain particle cluster coherence with a centre of mass
// guided by a psi wave? Without a directional affinity for each
// subagent, the mass simple drifts to the edges of the wave and
// loses coherence, even if we break the psi symmetry relative to M
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
	st[27] = ".........<..........................|"  // X
	st[28] = ".........<..........................|"  // X
	st[29] = ".........<..........................|"  // X
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
	st[44] = ".....m..............................|"
	st[45] = "....mmm.............................|"
	st[46] = "...mmmmm>...........................|"
	st[47] = "..mmmmmmm>..........................|"
	st[48] = "..mmmmmmm>..........................|"
	st[49] = "...mmmmm.>..........................|"
	st[50] = "....mmm.............................|"
	st[51] = ".....m..............................|"
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
		
		if neighbour != 0 {
			var breaker C.Message
			breaker.Value = C.AGENT[agent].Psi
			breaker.Phase = C.TICK
			C.CHANNEL[agent][neighbour] = breaker
		}
	}

	const PsiQuant = 1

	C.CausalIndependence(true)

	var pending_accept_offer[C.N] float64

	for t := 0; t < C.MAXTIME; t++ {
		
		// Every pair of agents has a private directional channel that's not overwritten by anyone else
		// Messages persist until they are read and cannot unseen

		for d := 0; d < C.N; d++ {

			//dbar := (d + C.N/2) % C.N
			
			var send,recv C.Message
			
			neighbour := C.AGENT[agent].Neigh[d]
			
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
				C.AGENT[agent] = EvolveAgents(C.AGENT[agent],d)

				// In this phase we can choose to make an offer to offload
				// a neighbour's mass, we have to have received an update first

				if pending_accept_offer[d] > 0 {
					send.Phase = C.TAKE
					send.Value = pending_accept_offer[d]

					C.AGENT[agent].Offer[d] = send.Value
					C.ConditionalChannelOffer(agent,neighbour,send)
				} else {
					send.MassID = C.AGENT[agent].MassID
					send.Value = C.AGENT[agent].Psi
					send.Angle = C.AGENT[agent].Theta
					send.Phase = C.TICK
					C.ConditionalChannelOffer(agent,neighbour,send)
				}

				continue

			case C.TAKE: // YOU
				transfer_offer := recv.Value  // This is massID now

				if WillingToPassToRecipient() {
					send.Value = transfer_offer
					C.AGENT[agent].Accept[d] = transfer_offer 
					C.AGENT[agent].MassID -= C.AGENT[agent].Accept[d]
					send.Phase = C.TACK
					C.ConditionalChannelOffer(agent,neighbour,send)
				} else {
					send.Phase = C.TICK                              // if not accepting, ignore
					send.Value = C.AGENT[agent].Psi                  // nothing reserved yet, so ok
					send.MassID = C.AGENT[agent].MassID
					send.Angle = C.AGENT[agent].Theta

					C.ConditionalChannelOffer(agent,neighbour,send)
				}
				continue

			case C.TACK: // ME, we are receiving a confirmation of our offer in .Offer
				   // (thank's for offer, I have reserved the amount now.. )
                                   // The actual amount you reserved to send me was

				transfer_offer := recv.Value

				if C.AGENT[agent].Offer[d] == transfer_offer {
					C.AGENT[agent].MassID += transfer_offer
					C.AGENT[agent].Accept[d] = transfer_offer
					send.Value = transfer_offer
				} else {
					// That's not what I agreed to
					send.Value = C.NOTACCEPT
				}

				C.AGENT[agent].Offer[d] = 0

				send.Phase = C.TOCK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
				
			case C.TOCK: // YOU - initiate a change / Xfer

				if recv.Value == C.NOTACCEPT {                                // my acceptance was refused
					// We're trusting the other side won't credit
					C.AGENT[agent].MassID += C.AGENT[agent].Accept[d]  // restore the reserved
					send.MassID = C.AGENT[agent].MassID
					send.Value = C.AGENT[agent].Psi
					send.Angle = C.AGENT[agent].Theta

				}

				C.AGENT[agent].Accept[d] = 0  // clear reservation, consider sent
				C.AGENT[agent].Offer[d] = 0
				
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
			}
		}

		// Now that all the directional information has been updated fairly,
		// we decide if/when to accept the movement of a mass unit

		for d := 0; d < C.N; d++ {
			dbar := (d + C.N/2) % C.N
			pending_accept_offer[d] = AcceptingMass(C.AGENT[agent],d,dbar)
			if pending_accept_offer[d] > 0 {
				break // we can only take one
			}
		}
	}
}

// ****************************************************************

func AcceptingMass(agent C.STAgent,d,dbar int) float64 {
	
	// We look to accept some mass from d if dbar looks empty

	if agent.Psi * agent.Psi < 0.1 {
		return 0
	}

	// Without a sense of direction for each subagent, the mass will just
	// dissipate to the edge of the psi field without particle coherence

	if (agent.M[d] > agent.MassID) {

		// if affinity gradient falling off ahead (dbar direction), then pull from rear

		if (agent.Psi * agent.Psi - agent.V[dbar] * agent.V[dbar]) > 0 {

			return 1
		}

	}
	return 0
}

// ****************************************************************

func WillingToPassToRecipient() bool {

	return true
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



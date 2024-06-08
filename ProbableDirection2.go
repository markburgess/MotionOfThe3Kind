///////////////////////////////////////////////////////////////
//
// Can we maintain particle cluster coherence with a centre of mass
// guided by a psi wave and probabilistic directional affinity?

// This initial condition is crucial in determining direction
// as the directional gradient is x,y symmetrical. Even with heavily
// skewed directions , the mass spreads out along the wave, since it
// travels quite slowly, depending on the relative speed / acceptance threshold.
// In this model, there is no centre of mass cohesive "force", so a body can
// just dissociate freely
//
///////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"time"
	C "Cellibrium"
)

// **********************************************************

const DoF = 10000
const wrange = 10000
const PERIOD = C.WAVELENGTH * wrange

// ****************************************************************

func main () {

	C.MODEL_NAME = "ProbableDirection"

	var st [C.Ylim]string

	 st[0] = ".........*..........................."
	 st[1] = ".........*..........................|"
	 st[2] = ".........*..........................|"
	 st[3] = ".........*..........................|"
	 st[4] = ".........*..........................|"
	 st[5] = ".........*..........................|"
	 st[6] = ".........*..........................|"
	 st[7] = ".........*..........................|"
	 st[8] = ".........*..........................|"
	 st[9] = ".........*..........................|"
	st[10] = ".........*..........................|"
	st[11] = ".........*..........................|"
	st[12] = ".........*..........................|"
	st[13] = ".........*..........................|"
	st[14] = ".........*..........................|"
	st[15] = ".........*..........................|"
	st[16] = ".........*..........................|"
	st[17] = ".........*..........................|"
	st[18] = ".........*..........................|"
	st[19] = ".........*..........................|"
	st[20] = ".........*..........................|"
	st[21] = ".........*..........................|"
	st[22] = ".........*..........................|"
	st[23] = ".........*..........................|"
	st[24] = ".........*..........................|"
	st[25] = ".........*..........................|"
	st[26] = ">...................................|"
	st[27] = ">..www.............................*|"  // X
	st[28] = ">..www.............................*|"  // X
	st[29] = ">..www.............................*|"  // X
	st[30] = ">...................................|"
	st[31] = ">...................................|"
	st[32] = ".........*..........................|"
	st[33] = ".........*..........................|"
	st[34] = ".........*..........................|"
	st[35] = ".........*..........................|"
	st[36] = ".........*..........................|"
	st[37] = ".........*..........................|"
	st[38] = ".........*..........................|"
	st[39] = ".........*..........................|"
	st[40] = ".........*..........................|"
	st[41] = ".........*..........................|"
	st[42] = ".........*..........................|"
	st[43] = ".........*..........................|"
	st[44] = ".........*..........................|"
	st[45] = ">...................................|"
	st[46] = ">..mmm.............................*|"
	st[47] = ">..mmm.............................*|"
	st[48] = ">..mmm.............................*|"
	st[49] = ">...................................|"
	st[50] = ">...................................|"
	st[51] = ".........*..........................|"
	st[52] = ".........*..........................|"
	st[53] = ".........*..........................|"
	st[54] = ".........*..........................|"
	st[55] = ".........*..........................|"
	st[56] = ".........*..........................|"
	st[57] = ".........*..........................|"
	st[58] = ".........*..........................|"
	st[59] = ".........*..........................|"
	st[60] = ".........*..........................|"
	st[61] = ".........*..........................|"
	st[62] = ".........*..........................|"
	st[63] = ".........*..........................|"
	st[64] = ".........*..........................|"
	st[65] = ".........*..........................|"
	st[66] = ".........*..........................|"
	st[67] = ".........*..........................|"
	st[68] = ".........*..........................|"
	st[69] = ".........*..........................|"
	st[70] = ".........*..........................|"
	st[71] = ".........*..........................|"
	st[72] = ".........*..........................|"
	st[73] = ".........*..........................|"
	st[74] = ".........*..........................|"
	st[75] = "...***...*..........................."


	// Keep the data structures for agents global too for convenience

	C.Initialize(st,DoF)
	C.ShowState(st,1,37,76,"+")
	EquilGuideRail()
	C.ShowState(st,C.MAXTIME,37,76,"+")
	//ShowMomentum(st,C.MAXTIME,37,76,"+")

	fmt.Println("")
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

	const mass = 1.0

	for direction := 0; direction < C.N; direction++ {

		neighbour := C.AGENT[agent].Neigh[direction]

		C.AGENT[agent].P[direction] = 0
		
		if neighbour != 0 {
			var breaker C.Message
			breaker.Value = C.AGENT[agent].Psi
			breaker.Intent = C.AGENT[agent].Intent
			breaker.Phase = C.TICK
			C.CHANNEL[agent][neighbour] = breaker
		}
	}

	C.CausalIndependence(true)

	for t := 0; t < C.MAXTIME; t++ {
		
		// Every pair of agents has a private directional channel that's not overwritten by anyone else
		// Messages persist until they are read and cannot unseen

		for direction := 0; direction < C.N ; direction++ {

			var send,recv C.Message

			d := direction % C.N
			dbar := (direction + C.N/2) % C.N
			
			neighbour := C.AGENT[agent].Neigh[d]
			
			if neighbour == 0 {
				continue // wall signal
			}

			// We need to wait for a positive signal indicating a new transfer to avoid double/empty reading

			recv = C.AcceptFromChannel(neighbour,agent)

			// ****************** PROCESS *********************

			switch recv.Phase {
				
			case C.TICK: // me

				C.AGENT[agent].V[d] = recv.Value       // the value here is psi
				C.AGENT[agent].M[d] = recv.MassID

				if C.AGENT[agent].M[d] > 0 {
					// only accept momentum from non-empty cells
					C.AGENT[agent].Intent = recv.Intent
				}

				C.AGENT[agent] = EvolveAgents(C.AGENT[agent],d)

				// In this phase we can choose to make an offer to accept
				// a neighbouring massID, we have to have received an update first

				if AcceptingMass(C.AGENT[agent],d,dbar) > 0 {
					send.Phase = C.TAKE
					send.Value = mass
					C.ConditionalChannelOffer(agent,neighbour,send)
				} else {
					send.MassID = C.AGENT[agent].MassID
					send.Value = C.AGENT[agent].Psi
					send.Intent = C.AGENT[agent].Intent
					send.Angle = C.AGENT[agent].Theta
					send.Phase = C.TICK
					C.ConditionalChannelOffer(agent,neighbour,send)
				}

			case C.TAKE: // YOU
				send.Phase = C.TACK
				C.AGENT[agent].MassID = 0     // assume I'll take it
				C.ConditionalChannelOffer(agent,neighbour,send)

			case C.TACK: // ME
				// Once we've accepted the momentum, rotate the direction clock
				C.AGENT[agent].MassID = mass
				C.AGENT[agent].Intent = C.Rotate(C.AGENT[agent].Intent)
				send.Phase = C.TOCK
				C.ConditionalChannelOffer(agent,neighbour,send)
				
			case C.TOCK: // YOU - initiate a change / Xfer

				C.AGENT[agent].MassID = 0
				send.MassID = C.AGENT[agent].MassID
				send.Intent = C.AGENT[agent].Intent
				send.Value = C.AGENT[agent].Psi
				send.Angle = C.AGENT[agent].Theta
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
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

	affinity := agent.V[d] * agent.V[d] - agent.Psi * agent.Psi

	const psi_threshold = 1.0   // should really express in dimensionless vars..?

	alignment := agent.Intent[0]
		
	if affinity > psi_threshold && alignment == dbar {

		if (agent.M[d] > 0) && (agent.MassID == 0) {
			return 1
		}
	}

	return 0
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

// ****************************************************************

func ShowMomentum(st_rows [C.Ylim]string,tmax,xlim,ylim int,mode string) {

	var fieldwidth,numwidth string

	switch mode {
	case "+": fieldwidth = fmt.Sprintf("%c%ds",'%',10)
	default:  fieldwidth = fmt.Sprintf("%c%ds",'%',8)
		numwidth = fmt.Sprintf("%c%d.1f",'%',8)
	}

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		count := 0.0
		mass_count := 0.0
		
		for y := 0; y < ylim; y++ {
			
			for x := 0; x < xlim; x++ {
				
				if !C.Blocked(st_rows,x,y) {

					mass := C.AGENT[C.COORDS[x][y]].MassID

					if mass > 0 {
						fmt.Printf(fieldwidth,Momentum(C.AGENT[C.COORDS[x][y]].Intent))
						mass_count += mass
						continue
					}

					observable := C.AGENT[C.COORDS[x][y]].Psi

					count += observable
  
					if observable != 0 {

						if mode == "+" {
							if observable > 1 {
								fmt.Printf(fieldwidth,"+")
							} else if observable < -1 {
								fmt.Printf(fieldwidth,"-")
							} else {
								fmt.Printf(fieldwidth,".")
							}

						} else {
							if observable*observable > 1 {
								fmt.Printf(numwidth,observable)
							} else {
								fmt.Printf(fieldwidth,".")
							}
						}
					} else {
						fmt.Printf(fieldwidth,".")
					}
					
				} else {
					fmt.Printf(fieldwidth," ")
				}
			}
			
			fmt.Println("")
		}

		fmt.Println("TOTAL =",count, "MASS=", mass_count)

		const noflicker = 10
		time.Sleep(noflicker * time.Duration(50) * time.Millisecond) // random noise
	}
}

// ***********************

func Momentum(p [C.MOMENTUMPROCESS]int) string {

	var s string

	for i := 0; i < C.MOMENTUMPROCESS; i++ {
		s += fmt.Sprintf("%d",p[i])
	}

	return s
}


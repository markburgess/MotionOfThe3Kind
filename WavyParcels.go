///////////////////////////////////////////////////////////////
//
// Cellular automaton version of conserved token passing.
// This is the pass the parcel protocol with conservation
// the only wave to ensure true conservation under virtual diffusion.
//
///////////////////////////////////////////////////////////////

package main

import (
	"os"
	"fmt"
	C "Cellibrium"
)

// **********************************************************

const momentum = 1
const psi = 0

const MODE = psi
const DoF = 10000

// ****************************************************************

func main () {

	C.MODEL_NAME = "ShortSourceSlits"

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

	C.ShowState(st,1,37,76,"+")
	EquilGuideRail()
	//C.ShowState(st,C.MAXTIME,37,76,"+")
	MyShowState(st,C.MAXTIME)
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
				
			case C.TICK: // me

				if recv.Value == C.CREDIT {

					C.AGENT[agent].Psi += C.AGENT[agent].Accept[direction]
					C.AGENT[agent].Accept[direction] = 0
					C.AGENT[agent].Offer[direction] = 0
					// return to update cycle

				} else { 
					C.AGENT[agent].V[direction] = recv.Value

					send.Phase = C.TAKE
					send.Value = DeltaPsi(agent)
					C.AGENT[agent].Offer[direction] = send.Value
					C.ConditionalChannelOffer(agent,neighbour,send)
					continue
					// return to update cycle
				}

				send.Value = C.AGENT[agent].Psi
				send.Phase = C.TICK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue

			case C.TAKE: // YOU
				transfer_offer := recv.Value

				send.Value = transfer_offer
				C.AGENT[agent].Accept[direction] = transfer_offer
				C.AGENT[agent].Psi -= C.AGENT[agent].Accept[direction]
				send.Phase = C.TACK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue

			case C.TACK:
				transfer_offer := recv.Value

				if C.AGENT[agent].Offer[direction] == transfer_offer {
					C.AGENT[agent].Accept[direction] = transfer_offer
					send.Value = transfer_offer
				} else {
					// That's not what I agreed to
					send.Value = C.NOTACCEPT
				}

				C.AGENT[agent].Offer[direction] = 0
				send.Phase = C.TOCK
				C.ConditionalChannelOffer(agent,neighbour,send)
				continue
				
			case C.TOCK: // YOU - initiate a change / Xfer

				if recv.Value == C.NOTACCEPT {

					C.AGENT[agent].Psi += C.AGENT[agent].Accept[direction]
					C.AGENT[agent].Psi++
					send.Value = C.AGENT[agent].Psi
				} else {
					send.Value = C.CREDIT
				}

				C.AGENT[agent].Accept[direction] = 0  // clear reservation, consider sent
				C.AGENT[agent].Offer[direction] = 0

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

func DeltaPsi(a int) int {
	
	const velocity = 9
	var d2 int = 0
	agent := C.AGENT[a]

	for di := 0; di < C.N; di++ {
		
		d2 += agent.V[di] - agent.Psi
	}
	
	// This is negative when Psi is higher than neighbours
	
	dtheta := d2 / (C.N * velocity)
	agent.Theta += dtheta
	deltaPsi := agent.Theta
	dpsi := deltaPsi / velocity
	return dpsi
	
}

// ****************************************************************

func MyShowState(st_rows [C.Ylim]string,tmax int) {
	
	var average float64 = 0

	var screen [C.Ylim]int

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		
		var total,visible int = 0,0
		var xf int = 0
		
		for y := 0; y < C.Ylim; y++ {
			
			for x := 0; x < C.Xlim; x++ {
				
				if st_rows[y][x] != '*' {

					const IsScreen = C.Xlim - 1

					observable := C.AGENT[C.COORDS[x][y]].Psi 
					xfer := C.AGENT[C.COORDS[x][y]].Xfer
					xf += xfer

					full := observable 
					full += C.AGENT[C.COORDS[x][y]].Accept[0] 
					full += C.AGENT[C.COORDS[x][y]].Accept[1] 
					full += C.AGENT[C.COORDS[x][y]].Accept[2] 
					full += C.AGENT[C.COORDS[x][y]].Accept[3]

					if x == IsScreen-1 {
						switch MODE {

						case psi: 	screen[y] += C.AGENT[C.COORDS[x][y]].Psi 
						case momentum: 	screen[y] += xfer
						}
					}

					if (x == IsScreen) {
						fmt.Printf("%20d",screen[y])
						continue
					}
  
					if observable != 0 {

						switch MODE {
						case psi: 	fmt.Printf("%11d",observable)
						case momentum: 	fmt.Printf("%11d",xfer)
						}

			
						total += full
						visible += observable
					} else {
						fmt.Printf("%11s",".")
					}
					
				} else {
					fmt.Printf("%11s"," ")
				}
			}
			
			fmt.Println("")
		}

		average += float64(total)
		fmt.Println("OBSERVABLE",visible,"TOTAL",total,"<av>=",average/float64(t),"XFER",xf)
		C.CausalIndependence(false)
	}
}

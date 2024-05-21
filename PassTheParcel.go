///////////////////////////////////////////////////////////////
//
// Cellular automaton version of conserved token passing.
// This is the pass the parcel protocol with conservation
// the only wave to ensure true conservation under virtual diffusion.
//
///////////////////////////////////////////////////////////////

package main

import (
	"fmt"
	"time"
	"math/rand"
	C "Cellibrium"
)

// **********************************************************

const momentum = 1
const psi = 0

const MODE = psi
const MATTER_LEVEL = 10000

// **********************************************************

const N = 4
const MAXPSI = 256
const MINPSI = 0

type STAgent struct {

	Psi    int
	Neigh  [N]int
	Offer  [N]int // last send
	Accept [N]int // last recv
	Xfer int
	t int
}

// **********************************************************

const EMPTY int = 0
const TICK int = 1 // default, should not be zero so we know when the channel is empty
const TOCK int = 2
const TAKE int = 3
const TACK int = 4

const CREDIT int = -1
const NOTACCEPT int = -1

// **********************************************************

type Message struct {
	Phase     int
	Value     int
	Direction int
}
// **********************************************************

const WAVE = "AAAA0000aaaa0000"
const WAVELENGTH = len(WAVE)

const Xlim = 24
const Ylim = 65
const Adim = Xlim * Ylim
const MAXTIME = 100000
const base_timescale = 15

var AGENT    [Adim]STAgent
var CHANNEL  [Adim][Adim]Message    // adjacency matrix
var ADJ      [Adim][Adim]int    // adjacency matrix
var COORDS   [Xlim][Ylim]int        // map [x,y] -> i

var S = rand.NewSource(time.Now().UnixNano())
var R = rand.New(S)

// ****************************************************************

func main () {

	C.MODEL_NAME = "PassTheParcel"

	var st [Ylim]string

	 st[0] = "************************"
	 st[1] = "*.........*............|"
	 st[2] = "*.........*............|"
	 st[3] = "*.........*............|"
	 st[4] = "*.........*............|"
	 st[5] = "*.........*............|"
	 st[6] = "*.........*............|"
	 st[7] = "*.........*............|"
	 st[8] = "*.........*............|"
	 st[9] = "*.........*............|"
	st[10] = "*.........*............|"
	st[11] = "*.........*............|"
	st[12] = "*.........*............|"
	st[13] = "*.........*............|"
	st[14] = "*.........*............|"
	st[15] = "*.........*............|"
	st[16] = "*.........*............|"
	st[17] = "*.........*............|"
	st[18] = "*.........*............|"
	st[19] = "*.........*............|"
	st[20] = "*.........*............|"
	st[21] = "*.........*............|"
	st[22] = "*.........*............|"
	st[23] = "*.........*............|"
	st[24] = "*.........*............|"
	st[25] = "*.........*............|"
	st[26] = "*......................|"
	st[27] = "*......................|"
	st[28] = "*.........*............|"
	st[29] = "*.........*............|"
	st[30] = "*.........*............|"
	st[31] = "*.........*............|"
	st[32] = ">.........*............|"
	st[33] = ">.........*............|"
	st[34] = "*.........*............|"
	st[35] = "*.........*............|"
	st[36] = "*.........*............|"
	st[37] = "*.........*............|"
	st[38] = "*......................|"
	st[39] = "*......................|"
	st[40] = "*.........*............|"
	st[41] = "*.........*............|"
	st[42] = "*.........*............|"
	st[43] = "*.........*............|"
	st[44] = "*.........*............|"
	st[45] = "*.........*............|"
	st[46] = "*.........*............|"
	st[47] = "*.........*............|"
	st[48] = "*.........*............|"
	st[49] = "*.........*............|"
	st[50] = "*.........*............|"
	st[51] = "*.........*............|"
	st[52] = "*.........*............|"
	st[53] = "*.........*............|"
	st[54] = "*.........*............|"
	st[55] = "*.........*............|"
	st[56] = "*.........*............|"
	st[57] = "*.........*............|"
	st[58] = "*.........*............|"
	st[59] = "*.........*............|"
	st[60] = "*.........*............|"
	st[61] = "*.........*............|"
	st[62] = "*.........*............|"
	st[63] = "*.........*............|"
	st[64] = "************************"

	// Keep the data structures for agents global too for convenience

	Initialize(st)

	ShowState(st,1)
	EquilGuideRail()
	ShowState(st,MAXTIME)
}

// ****************************************************************

func Initialize(st_rows [Ylim]string) {

	// Agent index begins at 1 .. < dimgraph

	var i int = 1

	// First need a lookup table from x,y -> i
	// Set up graph matrix and laplacian, and VECTOR PSI

	const d = 1.0 // standard hop distance

	// Make the adjacency matrix

	for y := 1; y < len(st_rows)-1; y++ {

		for x := 0; x < len(st_rows[y])-1; x++ {

			if st_rows[y][x] != '*' {
				COORDS[x][y] = i
				i++
			}
		}
	}

	for y := 1; y < Ylim - 1; y++ {

		for x := 0; x < Xlim - 1; x++ {

			switch st_rows[y][x] {

			case '.': InitAgentGeomAndAdj(x,y,0) 

			case '>': InitAgentGeomAndAdj(x,y,MATTER_LEVEL)

			default: InitAgentGeomAndAdj(x,y,0)
			}
		}
	}
}

//***********************************************************

func InitAgentGeomAndAdj(x,y int,amplitude int) {

	agent := COORDS[x][y]
	AGENT[agent].Psi = amplitude

	if agent == 0 {
		return
	}

	// adjacencies
	var w,e,n,s int

	if x > 0 {
		w = COORDS[x-1][y]
	} else {
		w = 0
	}

	e = COORDS[x+1][y]
	n = COORDS[x][y+1]
	s = COORDS[x][y-1]
	
	const d = 1.0

	var tot int = 0

	if n != 0 {
		ADJ[agent][n] = d
		ADJ[n][agent] = d
		tot++
	}

	if s != 0 {
		ADJ[agent][s] = d
		ADJ[s][agent] = d
		tot++
	}

	if e != 0 {
		ADJ[agent][e] = d
		ADJ[e][agent] = d
		tot++
	}

	if w != 0 {
		ADJ[agent][w] = d
		ADJ[w][agent] = d
		tot++
	}

	// Shortcut

	AGENT[agent].Neigh[0] = n
	AGENT[agent].Neigh[1] = s
	AGENT[agent].Neigh[2] = e
	AGENT[agent].Neigh[3] = w

}

// ****************************************************************

func EquilGuideRail() {

	for i := 1; i < Adim; i++ {
				
		go UpdateAgent_Flow(i)
	}
}

// ****************************************************************

func UpdateAgent_Flow(agent int) {
	
	// Start with an  unconditional promise to break the deadlock symmetry

	for direction := 0; direction < N; direction++ {

		neighbour := AGENT[agent].Neigh[direction]
		
		if neighbour != 0 {
			var breaker Message
			breaker.Value = AGENT[agent].Psi
			breaker.Phase = TICK
			breaker.Direction = direction
			CHANNEL[agent][neighbour] = breaker
		}
	}

	const PsiQuant = 1

	CausalIndependence(true)

	for t := 0; t < MAXTIME; t++ {
		
		// Every pair of agents has a private directional channel that's not overwritten by anyone else
		// Messages persist until they are read and cannot unseen

		for direction := 0; direction < N; direction++ {
			
			var send,recv Message
			
			neighbour := AGENT[agent].Neigh[direction]
			
			if neighbour == 0 {
				continue // wall signal
			}

			// We need to wait for a positive signal indicating a new transfer to avoid double/empty reading

			recv = AcceptFromChannel(neighbour,agent)

			// ****************** PROCESS *********************

			switch recv.Phase {
				
			case TICK: // me

				neighbourPsi := recv.Value
				
				// In this phase we can choose to make an offer to offload
				// a neighbour's Psi, we have to have received an update first
				// We issue a promise to take some psi, if the gradient is +ve
				
				const mingradient = PsiQuant+1 // if zero we can get trivial oscillations
				
				if GradientCapacity(neighbourPsi,AGENT[agent].Psi)  >= mingradient {
					send.Phase = TAKE
					send.Value = PsiQuant
					AGENT[agent].Offer[direction] = send.Value
					ConditionalChannelOffer(agent,neighbour,send)
					continue
				}

				send.Value = AGENT[agent].Psi
				send.Phase = TICK
				send.Direction = direction
				ConditionalChannelOffer(agent,neighbour,send)
				continue

			case TAKE: // YOU
				transfer_offer := recv.Value

				if WillingToTake(transfer_offer,agent) {                    // if we're willing to accept
					send.Value = transfer_offer
					AGENT[agent].Accept[direction] = transfer_offer     // reserve PAYMENT YOU
					AGENT[agent].Psi -= AGENT[agent].Accept[direction]  // reserve decr amount (single thr)  ** DO **
					send.Phase = TACK
					ConditionalChannelOffer(agent,neighbour,send)
				} else {
					send.Phase = TICK                                   // if not accepting, ignore
					send.Value = AGENT[agent].Psi                       // nothing reserved yet, so ok
					ConditionalChannelOffer(agent,neighbour,send)
				}
				continue

			case TACK: // ME, we are receiving a confirmation of our offer in .Offer
				   // (thank's for offer, I have reserved the amount now.. )
                                   // The actual amount you reserved to send me was

				transfer_offer := recv.Value

				if AGENT[agent].Offer[direction] == transfer_offer {
					AGENT[agent].Psi += transfer_offer
					AGENT[agent].Accept[direction] = transfer_offer // accept priv PAYMENT   ** STAGE WRITE **
					send.Value = transfer_offer
				} else {
					// That's not what I agreed to
					send.Value = NOTACCEPT
				}

				AGENT[agent].Offer[direction] = 0                       // clear or delete offer

				send.Phase = TOCK                                       // reset cycle
				ConditionalChannelOffer(agent,neighbour,send)
				continue
				
			case TOCK: // YOU - initiate a change / Xfer

				if recv.Value == NOTACCEPT {                                // my acceptance was refused
					// We're trusting the other side won't credit
					AGENT[agent].Psi += AGENT[agent].Accept[direction]  // restore the reserve amount   ** UNDO **
					send.Value = AGENT[agent].Psi                       // go back to update cycle
				}

				AGENT[agent].Accept[direction] = 0  // clear reservation, consider sent
				AGENT[agent].Offer[direction] = 0

				send.Phase = TICK
				ConditionalChannelOffer(agent,neighbour,send)
				continue
			}
		}
	}
}

// ****************************************************************

func EmptyMessage() Message {

	var m Message
	m.Value = 0
	m.Phase = 0
	return m
}

// ****************************************************************

func ConditionalChannelOffer(from,to int, mesg Message) {

	var recv Message

	for recv = CHANNEL[from][to]; recv != EmptyMessage(); recv = CHANNEL[from][to] {

		CausalIndependence(false)
	}

	AGENT[from].Xfer--
	CHANNEL[from][to] = mesg
	AGENT[to].Xfer++
}

// ****************************************************************

func AcceptFromChannel(neighbour,agent int) Message {

	var recv Message

	for recv = CHANNEL[neighbour][agent]; recv == EmptyMessage(); recv = CHANNEL[neighbour][agent] {

		CausalIndependence(false)
	}
	
	AGENT[neighbour].Xfer--
	CHANNEL[neighbour][agent] = EmptyMessage()
	AGENT[agent].Xfer++
	return recv
}

// ****************************************************************

func CausalIndependence(mode bool) {

	switch mode {
	case true: time.Sleep(time.Duration(2*base_timescale+R.Intn(50)) * time.Millisecond) // random noise
	case false: time.Sleep(time.Duration(2*base_timescale) * time.Millisecond)
	}
}

// ****************************************************************

func WillingToTake(offer,agent int) bool {

	return (AGENT[agent].Psi >= MINPSI + offer)
}

// ****************************************************************

func GradientCapacity (you, me int) int {

	return (you - me) //  && me < MAXPSI)
}

//***********************************************************

func GetRandOffset() int {

	return R.Intn(WAVELENGTH)
}

// ****************************************************************

func ShowState(st_rows [Ylim]string,tmax int) {
	
	var average float64 = 0

	var screen [Ylim]int

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		
		var total,visible int = 0,0
		var xf int = 0
		
		for y := 0; y < Ylim; y++ {
			
			for x := 0; x < Xlim; x++ {
				
				if st_rows[y][x] != '*' {

					const IsScreen = Xlim - 1

					observable := AGENT[COORDS[x][y]].Psi 
					xfer := AGENT[COORDS[x][y]].Xfer
					xf += xfer

					full := observable 
					full += AGENT[COORDS[x][y]].Accept[0] 
					full += AGENT[COORDS[x][y]].Accept[1] 
					full += AGENT[COORDS[x][y]].Accept[2] 
					full += AGENT[COORDS[x][y]].Accept[3]

					if x == IsScreen-1 {
						switch MODE {

						case psi: 	screen[y] += AGENT[COORDS[x][y]].Psi 
						case momentum: 	screen[y] += xfer
						}
					}

					if (x == IsScreen) {
						fmt.Printf("%18d",screen[y])
						continue
					}
  
					if observable != 0 {

						switch MODE {
						case psi: 	fmt.Printf("%8d",observable)
						case momentum: 	fmt.Printf("%8d",xfer)
						}

			
						total += full
						visible += observable
					} else {
						fmt.Printf("%8s",".")
					}
					
				} else {
					fmt.Printf("%8s"," ")
				}
			}
			
			fmt.Println("")
		}

		average += float64(total)
		fmt.Println("OBSERVABLE",visible,"TOTAL",total,"<av>=",average/float64(t),"XFER",xf)
		time.Sleep(time.Duration(base_timescale) * time.Millisecond) // random noise
	}
}

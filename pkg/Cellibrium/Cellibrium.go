
package Cellibrium

import (
	"os"
	"fmt"
	"time"
	"math"
	"math/rand"
)

// **********************************************************

const Xlim = 37
const Ylim = 76
const Adim = Xlim * Ylim
const base_timescale = 15  // smaller is faster
const MAXTIME = 100000

var   DOUBLE_SLIT bool = false
const WAVECHAR = "uuuu0000dddd0000"
const WAVELENGTH = len(WAVECHAR)
var   WAVE [WAVELENGTH]int
const N = 4

var MODEL_NAME string

// **********************************************************

type STAgent struct {

	Psi     int
	PsiDot  int
	Theta   int
	Neigh   [N]int
	V       [N]int
	P       [N]int
	Wave    [N][WAVELENGTH]int

	// Conservation equipment
	Offer  [N]int
	Accept [N]int
	Xfer int
	Cancel int
}

// **********************************************************

const EMPTY int = 0
const TICK int = 1 // default, should not be zero so we know when the channel is empty
const TOCK int = 2
const TAKE int = 3
const TACK int = 4

const CREDIT int = 1234567
const NOTACCEPT int = -1234567

// **********************************************************

type Message struct {

	Value   int
	Angle   int
	Phase   int
}

// **********************************************************

var AGENT    [Adim]STAgent
var CHANNEL  [Adim][Adim]Message    // adjacency matrix
var COORDS   [Xlim][Ylim]int        // map [x,y] -> i
var SCREEN   [Ylim]int
var POSITION int = 0
var FIRSTPOSITION int = 0

var S = rand.NewSource(time.Now().UnixNano())
var R = rand.New(S)

// ****************************************************************

func Initialize(st_rows [Ylim]string, DoF int) {

	// Agent index begins at 1 .. < dimgraph

	var i int = 1

	WAVE = MakeWaves(WAVECHAR)

	// First need a lookup table from x,y -> i
	// Set up graph matrix and laplacian, and VECTOR PSI

	const d = 1.0 // standard hop distance

	// Make the adjacency matrix

	for y := 1; y < len(st_rows)-1; y++ {

		for x := 0; x < len(st_rows[y])-1; x++ {

			if !Blocked(st_rows,x,y) {
				COORDS[x][y] = i
				i++
			}
		}
	}

	for y := 1; y < len(st_rows)-1; y++ {

		for x := 0; x < len(st_rows[y])-1; x++ {

			switch st_rows[y][x] {

			case '.': InitAgentGeomAndAdj(x,y,0) 

			case '>': InitAgentGeomAndAdj(x,y,DoF)

			case '<': if DOUBLE_SLIT {
				InitAgentGeomAndAdj(x,y,DoF)
			} else {
				InitAgentGeomAndAdj(x,y,0)
			}

			case '|': InitAgentGeomAndAdj(x,y,0)
				  SCREEN[y] = COORDS[x][y]

			default: InitAgentGeomAndAdj(x,y,0)
			}
		}
	}
}

//***********************************************************

func InitAgentGeomAndAdj(x,y int,amplitude int) {

	agent := COORDS[x][y]
	AGENT[agent].Psi = amplitude
	AGENT[agent].Theta = 0

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

	// Shortcut

	AGENT[agent].Neigh[0] = n
	AGENT[agent].Neigh[1] = s
	AGENT[agent].Neigh[2] = e
	AGENT[agent].Neigh[3] = w

	if amplitude > 0 {
		POSITION = agent
		FIRSTPOSITION = agent
	}
}

// ****************************************************************

func MovingPromise() {

	for t := 0; t < MAXTIME; t++ {

		CausalIndependence(true) // some noise

		Transition()

		if t % 1000 == 0 {
			POSITION = FIRSTPOSITION
		}

	}
}

// ****************************************************************

func Transition() {

	location := AGENT[POSITION] // the single privileged location promise

	//var peak int
	var selection int = -1
	var direction int = -1

	for di := 0; di < N; di++ {		

		if location.Neigh[di] == 0 {
			continue
		}

		av := (location.V[di] + location.Psi)/2

		if av != 0 {
			grad := 2*(location.V[di] - location.Psi)/av

			current := location.Psi * grad

			if current > selection {
				selection = current
				direction = di
			}
		}

	}

	if direction >= 0 {
		POSITION = location.Neigh[direction]
	}
	return

}

// ****************************************************************

func ShowState(st_rows [Ylim]string,tmax,xlim,ylim int,mode string) {
	
	var screen [Ylim]float64
	var fieldwidth,numwidth string

	switch mode {
	case "+": fieldwidth = fmt.Sprintf("%c%ds",'%',3)
	default:  fieldwidth = fmt.Sprintf("%c%ds",'%',6)
		numwidth = fmt.Sprintf("%c%dd",'%',6)
	}

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		count := 0
		
		for y := 0; y < ylim; y++ {
			
			for x := 0; x < xlim; x++ {
				
				if !Blocked(st_rows,x,y) {

					var IsScreen = xlim - 1

					observable := AGENT[COORDS[x][y]].Psi

					count += observable

					if x == IsScreen-1 {

						screen[y] += math.Sqrt(float64(observable) * float64(observable))
					}

					if (x == IsScreen) {
						fmt.Printf("%24.1f",screen[y])
						continue
					}
  
					if observable != 0 {

						if mode == "+" {
							if observable > 0 {
								fmt.Printf(fieldwidth,"+")
							} else {
								fmt.Printf(fieldwidth,"-")
							}
						} else {
							fmt.Printf(numwidth,observable)
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

		fmt.Println("TOTAL =",count)
		SaveScreenPattern("state",screen)

		const noflicker = 10
		time.Sleep(noflicker * time.Duration(base_timescale) * time.Millisecond) // random noise
	}
}

// ****************************************************************

func ShowPhase(st_rows [Ylim]string,tmax,xlim,ylim int) {
	
	var screen [Ylim]float64

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		
		for y := 0; y < ylim; y++ {
			
			for x := 0; x < xlim; x++ {
				
				if !Blocked(st_rows,x,y) {

					var IsScreen = xlim - 1

					observable := AGENT[COORDS[x][y]].Theta

					if x == IsScreen-1 {

						screen[y] += math.Sqrt(float64(observable) * float64(observable))
					}

					if (x == IsScreen) {
						fmt.Printf("%24.1f",screen[y])
						continue
					}
  
					if observable != 0 {
						fmt.Printf("%6d",observable)
					} else {
						fmt.Printf("%6s",".")
					}
					
				} else {
					fmt.Printf("%6s"," ")
				}
			}
			
			fmt.Println("")
		}

		SaveScreenPattern("phase",screen)

		const noflicker = 10

		time.Sleep(noflicker * time.Duration(base_timescale) * time.Millisecond) // random noise
	}
}

// ****************************************************************

func ShowPosition(st_rows [Ylim]string,tmax,xlim,ylim int) {

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		
		for y := 0; y < ylim; y++ {
			
			for x := 0; x < xlim; x++ {
				
				if !Blocked(st_rows,x,y) {

					if POSITION == COORDS[x][y] {
						
						fmt.Printf("%3s","X")
					} else {
						fmt.Printf("%3s",".")
					}
					
				} else {
					fmt.Printf("%3s"," ")
				}
			}
			
			fmt.Println("")
		}

		const noflicker = 10

		time.Sleep(noflicker * time.Duration(base_timescale) * time.Millisecond) // random noise
	}
}

// *************************************************************

func MakeWaves(process string) [WAVELENGTH]int {

	var wave [WAVELENGTH]int

	for pos := 0; pos < WAVELENGTH; pos++ {

		switch process[pos] {
			
		case 'u': wave[pos] = +1
			
		case 'd': wave[pos] = -1

		default: wave[pos] = 0
		}
	}

	return wave
}

//***********************************************************

func Cyc(pos int) int {

	if pos < 0 {
		return (2 * WAVELENGTH + pos) % WAVELENGTH
	} else {
	
		return pos % WAVELENGTH
	}
}

//***********************************************************

func SaveScreenPattern(name string,pattern [Ylim]float64) {

	var output string = ""
	var dosave bool = false

	for y := 0; y < Ylim; y++ {
		
		if pattern[y] > 0 {
			dosave = true
		}

		output += fmt.Sprintf("%f\n",pattern[y])
	}

	output += fmt.Sprintf("\n")

	if dosave {
		if  DOUBLE_SLIT {
			os.WriteFile(MODEL_NAME+"_DS", []byte(output), 0644)
		} else  {
			os.WriteFile(MODEL_NAME+"_SS", []byte(output), 0644)
		}
	}

}

//***********************************************************

func Blocked(st_rows [Ylim]string, x,y int) bool {

	if DOUBLE_SLIT {	
		return (st_rows[y][x] == '*')
	} else {
		return (st_rows[y][x] == 'X' || st_rows[y][x] == '*')
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

func AcceptFromChannel(neighbour,agent int) Message {

	var recv Message

	for recv = CHANNEL[neighbour][agent]; recv == EmptyMessage(); recv = CHANNEL[neighbour][agent] {

		CausalIndependence(false)
	}
	
	CHANNEL[neighbour][agent] = EmptyMessage()
	return recv
}

// ****************************************************************

func ConditionalChannelOffer(from,to int, mesg Message) {

	var recv Message

	for recv = CHANNEL[from][to]; recv != EmptyMessage(); recv = CHANNEL[from][to] {

		CausalIndependence(false)
	}

	CHANNEL[from][to] = mesg
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

	return (AGENT[agent].Psi >= 0 + offer)
}

// ****************************************************************

func GradientCapacity (you, me int) int {

	return (you - me) //  && me < MAXPSI)
}

//***********************************************************

func GetRandOffset() int {

	return R.Intn(WAVELENGTH)
}


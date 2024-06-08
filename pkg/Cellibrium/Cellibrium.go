
package Cellibrium

import (
	"os"
	"fmt"
	"time"
	"math"
	"math/rand"
)

// **********************************************************

const Xlim = 56
const Ylim = 76
const Adim = Xlim * Ylim
const base_timescale = 15  // smaller is faster
const MAXTIME = 100000

var   DOUBLE_SLIT bool = false
const WAVECHAR = "uuuu0000dddd0000"
const WAVELENGTH = len(WAVECHAR)
const MOMENTUMPROCESS = 8
var   WAVE [WAVELENGTH]float64
const N = 4

var MODEL_NAME string

// **********************************************************

type STAgent struct {

	Psi      float64    // my amplitude
	PsiDot   float64    // my velocity
	Theta    float64    // my phase
	MassID   float64    // my mass
	ID       byte       // my semantic label
	Moment   int        // direction of mass process

	Intent   [MOMENTUMPROCESS]int // directional encoding for process wave (program counter)

	// Neighbour cache

	Neigh    [N]int       // channel
	V        [N]float64   // psi
	M        [N]float64   // mass
	P        [N]int       // momentum (cache of Intent[0])
	NeighID  [N]byte      // semantic label

	// Conservation equipment

	Offer  [N]float64
	Accept [N]float64
	Xfer int
	Cancel int
}

// **********************************************************

const DIMENSION = 2
const NORTH = 0
const EAST  = 1
const SOUTH = 2
const WEST  = 3

// **********************************************************

const EMPTY int = 0
const TICK int = 1 // default, should not be zero so we know when the channel is empty
const TOCK int = 2
const TAKE int = 3
const TACK int = 4

const CREDIT int = 1234567
const NOTACCEPT = -1234567.0

// **********************************************************

type Message struct {

	Value    float64
	Angle    float64
	MassID   float64
	Moment   int
	Intent   [MOMENTUMPROCESS]int
	Phase    int      // proto phase
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

func Initialize(st_rows [Ylim]string, DoF float64) {

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

			AGENT[COORDS[x][y]].ID = st_rows[y][x]

			switch st_rows[y][x] {

			case '.': InitAgentGeomAndAdj(x,y,0) 

			case 'm': InitAgentMassGeomAndAdj(x,y,DoF/2,"SSEEEEEE") 

			case 'w': InitAgentMassGeomAndAdj(x,y,DoF/2,"NNEEEEEE") 

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

func InitAgentMassGeomAndAdj(x,y int,amplitude float64, momentum string) {

	agent := COORDS[x][y]

	InitAgentGeomAndAdj(x,y,0) 

	AGENT[agent].MassID = 1
	AGENT[agent].Intent = MakeMomentum(momentum) 
}

//***********************************************************

func InitAgentGeomAndAdj(x,y int,amplitude float64) {

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

	AGENT[agent].Neigh[NORTH] = n
	AGENT[agent].Neigh[SOUTH] = s
	AGENT[agent].Neigh[EAST] = e
	AGENT[agent].Neigh[WEST] = w

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
	var selection float64 = -1
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
	default:  fieldwidth = fmt.Sprintf("%c%ds",'%',8)
		numwidth = fmt.Sprintf("%c%d.1f",'%',8)
	}

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		count := 0.0
		mass_count := 0.0
		
		for y := 0; y < ylim; y++ {
			
			for x := 0; x < xlim; x++ {
				
				if !Blocked(st_rows,x,y) {

					var IsScreen = xlim - 1

					mass := AGENT[COORDS[x][y]].MassID

					if mass > 0 {
						fmt.Printf(fieldwidth,"M")
						mass_count += mass
						continue
					}

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
		SaveScreenPattern("state",screen)

		const noflicker = 10
		time.Sleep(noflicker * time.Duration(base_timescale) * time.Millisecond) // random noise
	}
}

// ****************************************************************

func ShowAffinity(st_rows [Ylim]string,tmax,xlim,ylim int) {
	
	var screen [Ylim]float64
	var fieldwidth,numwidth string

	fieldwidth = fmt.Sprintf("%c%ds",'%',8)
	numwidth = fmt.Sprintf("%c%d.1f",'%',8)

	for t := 1; t < tmax; t++ {
		
		fmt.Printf("\x1b[2J") // CLS
		count := 0.0
		
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

					affinity := math.Log(observable*observable)
  
					if observable > 1 {
						fmt.Printf(numwidth,affinity)
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

				// Pumping test
				// if st_rows[x][y] == '>' {
				//	AGENT[COORDS[x][y]].Psi = 44000
				// }
				
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

func MakeWaves(process string) [WAVELENGTH]float64 {

	var wave [WAVELENGTH]float64

	for pos := 0; pos < WAVELENGTH; pos++ {

		switch process[pos] {
			
		case 'u': wave[pos] = +1
			
		case 'd': wave[pos] = -1

		default: wave[pos] = 0
		}
	}

	return wave
}

// *************************************************************

func MakeMomentum(process string) [MOMENTUMPROCESS]int {

	var wave [MOMENTUMPROCESS]int

	for pos := 0; pos < MOMENTUMPROCESS; pos++ {

		switch process[pos] {
			
		case 'E': wave[pos] = EAST

		case 'N': wave[pos] = NORTH

		case 'S': wave[pos] = SOUTH

		default: wave[pos] = -1
		}
	}

	fmt.Println("WAVE",wave)
	return wave
}

//***********************************************************

func Rotate(p [MOMENTUMPROCESS]int) [MOMENTUMPROCESS]int {

	var newp [MOMENTUMPROCESS]int

	for pos := 0; pos < MOMENTUMPROCESS; pos++ {

		newp[pos] = p[(pos+1) % MOMENTUMPROCESS]
	}

	return newp
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

func WillingToTake(offer float64,agent int) bool {

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

//***********************************************************

func GetRandDirection() int {

	return R.Intn(N)
}

//***********************************************************

func GetRandXY(ptotal float64) float64 {

	return R.Float64() * ptotal
}
///////////////////////////////////////////////////////////////
//
// A bistable oscillator to figure out the wave, phase relationships
//
// In this model, we introduce a positive only energy parameter
// which can be used to prevent borrowing. Energy diffuses, amplitude
// finances...
//
///////////////////////////////////////////////////////////////

package main

import (
	"fmt"
)

// ****************************************************************

func main () {

	var psi1,psi2,theta1,theta2 float64
	var ppsi1,ppsi2 float64

	psi1 = 100

	for t := 0; t < 10000; t++ {

		dpsi1,dtheta1 := SenseNearby(ppsi1,psi1,theta1,psi2)
		dpsi2,dtheta2 := SenseNearby(ppsi2,psi2,theta2,psi1)

		// This is vulnerable to rounding accuracy, so we MUST have float

		psi1 += dpsi1
		psi2 += dpsi2
		theta1 += dtheta1
		theta2 += dtheta2

		// Artificial energy conservation, correct overshoot error
		// how can we achieve this detailed balance without cheating?

		Show(psi1,theta1,psi2,theta2)
		ppsi1 = psi1
		ppsi2 = psi2
	}
}

// ****************************************************************

func SenseNearby(prevme,me,theta,you float64) (float64,float64) { // Please give me

	// This is based on full trust, just do as we're told.
	// Unless there's some limit, this can keep borrowing from
	// other to produce runaway inflation

	const dt = 0.01
	const mass = 10  // blow up

	grad := you - me                // Your surplus relative to me tells me to balance you with neg feedback
                                        // (force direction) gradient makes directionality from neighbour alone
                                        // we assume the whole force is represented in the displacement -kx

	deltaTheta := grad / mass *  dt // small increments

	DeltaPsi := theta * dt

	return DeltaPsi, deltaTheta 
}

// ****************************************************************

func Show(psi1,theta1,psi2,theta2 float64) {

	energy2 := psi2 + theta2*theta2
	energy1 := psi1 + theta1*theta1 

	fmt.Printf("%8.2f %8.2f   ----> %f,,",psi1,theta1,psi1+psi2)
	fmt.Printf("%8.2f %8.2f (%.2f-%.2f)\n",theta2,psi2,energy1,energy2)
}
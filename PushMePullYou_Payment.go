///////////////////////////////////////////////////////////////
//
// A bistable oscillator to figure out the wave, phase relationships
//
// In this model, both sides can borrow from one another without limit
// They only see what the other seems to have, and while the sum
// of both sides is always zero, the oscillations grow unimpeded.
//
///////////////////////////////////////////////////////////////

package main

import (
	"fmt"
)

const PERIOD = 100000

// ****************************************************************

func main () {

	var psi1,psi2,theta1,theta2 float64

	psi1 = 0
	theta1=100

	for t := 0; t < 10000; t++ {

	//fmt.Printf("-----------------------------------------\n") // CLS
	//fmt.Printf(" PSI1 THETA1   dpsi1 ; dpsi2   THETA2  PSI2\n") // CLS
	//fmt.Printf("-----------------------------------------\n") // CLS
		dpsi1,dtheta1 := PleaseGiveMe(psi1,theta1,psi2)
		dpsi2,dtheta2 := PleaseGiveMe(psi2,theta2,psi1)

		//Show(psi1,dpsi1,theta1,psi2,dpsi2,theta2)

		// Exchange

		psi1 = psi1 + dpsi1
		psi2 = psi2 - dpsi1
		theta1 = theta1 + dtheta1 //% PERIOD

		//Show(psi1,dpsi1,theta1,psi2,dpsi2,theta2)

		psi2 = psi2 + dpsi2
		psi1 = psi1 - dpsi2
		theta2 = theta2 + dtheta2 //% PERIOD

		Show(psi1,dpsi1,theta1,psi2,dpsi2,theta2)
	}
}

// ****************************************************************

func PleaseGiveMe(me,theta,you float64) (float64,float64) { // Please give me

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

func Show(psi1,dpsi1,theta1,psi2,dpsi2,theta2 float64) {

	fmt.Printf("%8.2f %8.2f %8.2f   ----> %f\n",psi1,theta1,dpsi1,psi1+psi2)
	//fmt.Printf("%8d %8d %8d\n",dpsi2,theta2,psi2)
}
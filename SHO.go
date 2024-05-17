///////////////////////////////////////////////////////////////
//
// Simple harmonic oscillator
//
///////////////////////////////////////////////////////////////

package main

import (
	"fmt"
)

// ****************************************************************

func main () {

	const k = 1.00
	const mass = 1.0
	const dt = 0.01

	var Psi,PsiDot float64

	Psi = 1.0

	for t := 0; t < 10000; t++ {

		acc :=  -k * Psi
		
		deltaV := acc / mass * dt
		PsiDot += deltaV

		DeltaPsi := PsiDot * dt
		Psi += DeltaPsi

		fmt.Printf("%8.2f \n",Psi)
	}
}


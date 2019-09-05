package ldpc

import (
	"fmt"
	"math"
)

const (
	BigInfinity = 1000000.0
	Inf         = 64.0

	// These parameters are only used for the decoding function.
	maxIter  = 20   // The maximum number of iteration in the decoder
	crossErr = 0.01 // A transisient error probability. This is also fixed as a small value

	printRowInCol bool = false
	printColInRow bool = true

	printHashVector bool = false
	printOutputWord bool = true
)

//Parameters is used to determine matrix size
type Parameters struct {
	n    int
	m    int
	wc   int
	wr   int
	seed int
}

func funcF(x float64) float64 {
	if x >= BigInfinity {
		return (1.0 / BigInfinity)
	} else if x <= (1.0 / BigInfinity) {
		return BigInfinity
	} else {
		return (math.Log((math.Exp(x) + 1) / (math.Exp(x) - 1)))
	}
}

func infinityTest(x float64) float64 {
	if x >= Inf {
		return Inf
	} else if x <= -Inf {
		return -Inf
	} else {
		return x
	}
}

//PrintWord print Hash vector or Outputword using flag
func PrintWord(src []int, flag bool) {
	switch flag {
	case printHashVector:
		fmt.Println("hash vector")
	case printOutputWord:
		fmt.Println("OutputWord")
	default:
		fmt.Println("Check flag again")
	}

	for _, i := range src {
		fmt.Printf("%d", i)
	}
	fmt.Printf("\n")
}

//PrintQ print RowInCol or ColInRow using flag
func PrintQ(src [][]int, flag bool) {
	switch flag {
	case printRowInCol:
		fmt.Println("row in col")
	case printColInRow:
		fmt.Println("col in row")
	default:
		fmt.Println("Check flag again")
		return
	}

	for _, i := range src {
		for _, j := range i {
			fmt.Printf("%d ", j)
		}
		fmt.Println()
	}
}

//PrintH print parameter of matrix and seed
func PrintH(parameters Parameters) {
	fmt.Printf("The value of seed : %d\n", parameters.seed)
	fmt.Printf("The size of H is %d x %d with ", parameters.m, parameters.n)
	fmt.Printf("wc : %d and wr : %d \n", parameters.wc, parameters.wr)
}

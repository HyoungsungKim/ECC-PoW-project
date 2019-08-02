package ldpc

import (
	"fmt"
	"math"
)

const (
	BigInfinity = 1000000.0
	Inf         = 64.0
)

const (
	printHashVector = 1
	printOutputWord = 2
)

var n, m, wc, wr, seed int

var hashVector []int
var outputWord []int

//hashVector := make([]int, m)
//outputWord := make([]int, n)

var tmpHashVector [32]byte //32bytes => 256 bytes

var H [][]int
var rowInCol [][]int
var colInRow [][]int

//H 			:= make([][]int)
//rowInCol 	:= make([][]int)
//colInRow  	:= make([][]int)

// These parameters are only used for the decoding function.

var maxIter = 20    // The maximum number of iteration in the decoder
var crossErr = 0.01 // A transisient error probability. This is also fixed as a small value

var LRft []float64
var LRpt []float64
var LRrtl [][]float64
var LRqtl [][]float64

//LRft := make([]float64)
//LRpt := make([]float64)
//LRrtl := make([][]float64)
//LRqtl := male([][]float64)

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

func printWord(flag int) {
	var buffer []int
	if flag == printHashVector {
		buffer = hashVector
		fmt.Println("hash vector")
	} else if flag == printOutputWord {
		buffer = outputWord
		fmt.Println("OutputWord")
	} else {
		fmt.Println("Check flag again")
	}
	for _, i := range buffer {
		fmt.Printf("%d", i)
	}
	fmt.Printf("\n\n")
}

func printH() {
	fmt.Printf("The value of seed : %d\n", seed)
	fmt.Printf("The size of H is %d x %d with ", m, n)
	fmt.Printf("wc : %d and wr : %d \n", wc, wr)
}

func TestFunc() {
	for i := 0; i < 2; i++ {
		setDifficultyUsingLevel(3)

		//GenerateSeed(i)
		seed = i
		GenerateH()
		GenerateQ()
		GenerateHv([]byte("cdexff12fff3ffff3f3ff3fff3f3f3feeeed"))

		Decoding()
		printWord(printHashVector)
		printWord(printOutputWord)
		Decision()
	}
}

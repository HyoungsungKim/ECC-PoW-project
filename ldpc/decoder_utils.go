package ldpc

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	BigInfinity = 1000000.0
	Inf         = 64.0
)

const (
	printRowInCol bool = false
	printColInRow bool = true

	printHashVector bool = false
	printOutputWord bool = true
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

func PrintWord(flag bool) {
	var buffer []int
	switch flag {
	case printHashVector:
		buffer = hashVector
		fmt.Println("hash vector")
	case printOutputWord:
		buffer = outputWord
		fmt.Println("OutputWord")
	default:
		fmt.Println("Check flag again")
	}

	for _, i := range buffer {
		fmt.Printf("%d", i)
	}
	fmt.Printf("\n")
}

func PrintQ(flag bool) {
	var buffer [][]int
	switch flag {
	case printRowInCol:
		buffer = rowInCol
		fmt.Println("row in col")
	case printColInRow:
		buffer = colInRow
		fmt.Println("col in row")
	default:
		fmt.Println("Check flag again")
		return
	}

	for _, i := range buffer {
		for _, j := range i {
			fmt.Printf("%d ", j)
		}
		fmt.Println()
	}
}

func PrintH() {
	fmt.Printf("The value of seed : %d\n", seed)
	fmt.Printf("The size of H is %d x %d with ", m, n)
	fmt.Printf("wc : %d and wr : %d \n", wc, wr)
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func TestFunc() {
	tickerCounter := 0
	ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	var nonce int64
	nonce = 0
	SetDifficultyUsingLevel(0)

	var previousHashValue = "0ffff123fff"
	var currentBlockHeader = previousHashValue + time.Now().UTC().String()
	var currentBlockHeaderWithNonce string

	GenerateSeed([]byte(previousHashValue))
	GenerateH()
	GenerateQ()

	PrintH()
	//PrintQ(printRowInCol)
	//PrintQ(printColInRow)

	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Decoding")
	for {
		fmt.Printf("\rDecoding %s", ticker[tickerCounter])
		tickerCounter++
		tickerCounter %= len(ticker)

		currentBlockHeaderWithNonce = currentBlockHeader + strconv.FormatInt(nonce, 10)

		GenerateHv([]byte(currentBlockHeaderWithNonce))

		Decoding()
		flag := Decision()

		if !flag {
			Decoding()
			flag = Decision()
		}
		if flag {
			fmt.Printf("\nCodeword is founded with nonce = %d\n", nonce)
			break
		}
		nonce++
	}

	PrintWord(printHashVector)
	PrintWord(printOutputWord)
	fmt.Printf("\n")
}

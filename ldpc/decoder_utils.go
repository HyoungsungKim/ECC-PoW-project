package ldpc

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

const (
	BigInfinity = 1000000.0
	Inf         = 64.0
	MaxNonce    = 1<<32 - 1

	// These parameters are only used for the decoding function.
	maxIter  = 20   // The maximum number of iteration in the decoder
	crossErr = 0.01 // A transisient error probability. This is also fixed as a small value

	printRowInCol bool = false
	printColInRow bool = true

	printHashVector bool = false
	printOutputWord bool = true
)

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

//TestFunc test decoder.go functions
func TestFunc() {
	tickerCounter := 0
	ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	var LDPCNonce uint32
	var hashVector []int
	var outputWord []int

	tempPrevHash := "00000000000000000000000000000123"

	header := ethHeader{}
	copy(header.ParentHash[:], tempPrevHash)
	header.Time = uint64(time.Now().Unix())
	var currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
	var currentBlockHeaderWithNonce string

	parameters := SetDifficultyUsingLevel(1)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	PrintH(parameters)
	//PrintQ(printRowInCol)
	//PrintQ(printColInRow)

	for {
		fmt.Printf("\rDecoding %s", ticker[tickerCounter])
		tickerCounter++
		tickerCounter %= len(ticker)

		//If Nonce is bigger than MaxNonce, then update timestamp
		if LDPCNonce >= MaxNonce {
			LDPCNonce = 0
			header.Time = uint64(time.Now().Unix())
			currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
		}
		currentBlockHeaderWithNonce = currentBlockHeader + strconv.FormatUint(uint64(LDPCNonce), 10)

		hashVector = GenerateHv(parameters, []byte(currentBlockHeaderWithNonce))
		hashVector, outputWord = Decoding(parameters, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(parameters, colInRow, outputWord)

		if !flag {
			hashVector, outputWord = Decoding(parameters, hashVector, H, rowInCol, colInRow)
			flag = MakeDecision(parameters, colInRow, outputWord)
		}
		if flag {
			fmt.Printf("\nCodeword is founded with nonce = %d\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}

	/*
		fmt.Printf("LRft : \n %v\n", LRft)
		fmt.Printf("LRpt : \n %v\n", LRpt)
		fmt.Printf("LRrtl : \n %v\n", LRrtl)
		fmt.Printf("LRft : \n %v\n", LRqtl)
	*/

	PrintWord(hashVector, printHashVector)
	PrintWord(outputWord, printOutputWord)
	fmt.Printf("\n")
}

//TestRunLDPC test runLDPC function
func TestRunLDPC() {
	parameters := SetDifficultyUsingLevel(0)
	tempParentHash := "00000000000000000000000000000123"

	tempHeader := ethHeader{}
	copy(tempHeader.ParentHash[:], tempParentHash)
	tempHeader.Time = uint64(time.Now().Unix())

	runLDPC(parameters, tempHeader)
}

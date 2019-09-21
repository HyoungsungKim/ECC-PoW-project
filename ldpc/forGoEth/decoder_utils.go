package ldpc

import (
	"crypto/sha256"
	"fmt"
	"math"
	"strconv"
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

// SetParameters sets n, wc, wr, m, seed return parameters and difficulty level
func SetParameters(header ethHeader) (Parameters, int) {
	level := SearchLevel(header.Difficulty)

	parameters := Parameters{
		n:  Table[level].n,
		wc: Table[level].wc,
		wr: Table[level].wr,
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)
	parameters.seed = generateSeed(header.ParentHash)

	return parameters, level
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

//RunOptimizedLDPC is implemented to test decoding process
func RunOptimizedLDPC(block ethHeader) ([]int, []int, uint64) {
	//Need to set difficulty before running LDPC
	parameters, _ := SetParameters(block)
	LDPCNonce := generateRandomNonce()
	var hashVector []int
	var outputWord []int
	//	var LRrtl [][]float64

	var serializedHeader = string(block.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce string
	var encryptedHeaderWithNonce [32]byte
	H := generateH(parameters)
	colInRow, rowInCol := generateQ(parameters, H)

	for {
		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(LDPCNonce, 10)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = generateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ = OptimizedDecoding(block, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(block, colInRow, outputWord)

		if flag {
			fmt.Printf("Codeword is founded with nonce = %d\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}

	return hashVector, outputWord, LDPCNonce
}

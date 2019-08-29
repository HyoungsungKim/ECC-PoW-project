package ldpc

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

func TestDecodingElapseTime(t *testing.T) {
	var hashVector []int
	//var LRrtl [][]float64

	tempPrevHash := "00000000000000000000000000000000"

	header := ethHeader{}
	copy(header.ParentHash[:], tempPrevHash)
	var currentBlockHeader = string(header.ParentHash[:])

	parameters := SetDifficultyUsingLevel(0)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	hashVector = GenerateHv(parameters, []byte(currentBlockHeader))
	hashVector, _, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)

	for i := 0; i < 100000; i++ {
		_, _, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
	}
}

//TestDecodingProcess test decoder.go functions
func TestDecodingProcess(t *testing.T) {
	tickerCounter := 0
	ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	var LDPCNonce uint32
	var hashVector []int
	var outputWord []int
	//var LRrtl [][]float64

	tempPrevHash := "00000000000000000000000000000000"

	header := ethHeader{}
	copy(header.ParentHash[:], tempPrevHash)
	header.Time = uint64(time.Now().Unix())
	var currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
	var currentBlockHeaderWithNonce string

	parameters := SetDifficultyUsingLevel(0)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

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
		hashVector, outputWord, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(parameters, colInRow, outputWord)

		if !flag {
			hashVector, outputWord, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
			flag = MakeDecision(parameters, colInRow, outputWord)
		}
		if flag {
			t.Logf("\nCodeword is founded with nonce = %v\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}
}

//TestRunLDPC test runLDPC function
func TestRunLDPC(t *testing.T) {
	parameters := SetDifficultyUsingLevel(0)
	var tempParentHash [32]byte
	//tempParentHash = [0, 0, ..., 0]
	parameters.seed = GenerateSeed(tempParentHash)

	tempHeader := ethHeader{}
	tempHeader.ParentHash = tempParentHash
	tempHeader.Time = uint64(time.Now().Unix())

	RunLDPC(parameters, tempHeader)
}

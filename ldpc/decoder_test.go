package ldpc

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"
)

func TestOptimizedDecodingImplement(t *testing.T) {
	for i := 0; i < 100000; i++ {
		header := ethHeader{}

		var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
		var serializedHeaderWithNonce = serializedHeader + ""
		var encryptedHeaderWithNonce [32]byte

		var hashVector []int

		parameters := SetDifficultyUsingLevel(0)
		parameters.seed = GenerateSeed(header.ParentHash)

		H := GenerateH(parameters)
		colInRow, rowInCol := GenerateQ(parameters, H)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ := Decoding(parameters, hashVector, H, rowInCol, colInRow)

		opHeader := ethHeader{}

		var opSerializedHeader = string(opHeader.ParentHash[:]) // + ... + string(header.MixDigest)
		var opSerializedHeaderWithNonce = opSerializedHeader + ""
		var opEncryptedHeaderWithNonce [32]byte

		var opHashVector []int

		opParameters := SetDifficultyUsingLevel(0)
		opParameters.seed = GenerateSeed(header.ParentHash)

		opH := GenerateH(opParameters)
		opColInRow, opRowInCol := GenerateQ(opParameters, opH)
		opEncryptedHeaderWithNonce = sha256.Sum256([]byte(opSerializedHeaderWithNonce))

		opHashVector = GenerateHv(opParameters, opEncryptedHeaderWithNonce)

		opHashVector, opOutputWord, _ := OptimizedDecoding(opParameters, opHashVector, opH, opRowInCol, opColInRow)

		if !reflect.DeepEqual(hashVector, opHashVector) || !reflect.DeepEqual(outputWord, opOutputWord) {
			t.Errorf("Decoder hashVector :  %v\n", hashVector)
			t.Errorf("OptimezedDecoder hashVector: %v\n", opHashVector)

			t.Errorf("Decoder outputWord :  %v\n", outputWord)
			t.Errorf("OptimezedDecoder outputWord: %v\n", opOutputWord)
		}
	}
}

func TestDecodingElapseTime(t *testing.T) {
	header := ethHeader{}

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce = serializedHeader + ""
	var encryptedHeaderWithNonce [32]byte

	var hashVector []int
	//var LRrtl [][]float64

	parameters := SetDifficultyUsingLevel(1)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)
	encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)

	for i := 0; i < 10000; i++ {
		//	startTime := time.Now()
		_, _, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
		//	elapseTime := time.Since(startTime)
		//	t.Logf("%v", elapseTime)
	}
}

func TestOptimizedDecodingElapseTime(t *testing.T) {
	header := ethHeader{}

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce = serializedHeader + ""
	var encryptedHeaderWithNonce [32]byte

	var hashVector []int
	//var LRrtl [][]float64

	parameters := SetDifficultyUsingLevel(1)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)
	encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)

	for i := 0; i < 10000; i++ {
		_, _, _ = OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)
	}
}

//TestDecodingProcess test decoder.go functions
func TestDecodingProcess(t *testing.T) {
	//tickerCounter := 0
	//ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	var LDPCNonce uint32
	var hashVector []int
	var outputWord []int
	//var LRrtl [][]float64

	header := ethHeader{}

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce string
	var encryptedHeaderWithNonce [32]byte

	parameters := SetDifficultyUsingLevel(0)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	for {
		startTime := time.Now()
		//	fmt.Printf("\rDecoding %s", ticker[tickerCounter])
		//	tickerCounter++
		//	tickerCounter %= len(ticker)

		//If Nonce is bigger than MaxNonce, then update timestamp
		if LDPCNonce >= MaxNonce {
			LDPCNonce = 0
			//header.Time = uint64(time.Now().Unix())
			//currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
		}
		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(uint64(LDPCNonce), 10)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(parameters, colInRow, outputWord)

		elapseTime := time.Since(startTime)

		if LDPCNonce%10000 == 0 {
			fmt.Printf("1 cycle decoding elapse Time : %v\n", elapseTime)
			fmt.Printf("hashVector : %v\n", hashVector)
			fmt.Printf("outputWord : %v\n", outputWord)
			fmt.Printf("LDPC Nonce : %v\n", LDPCNonce)
		}

		if flag {
			t.Logf("\nCodeword is founded with nonce = %v\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}
	t.Logf("hashVector : %v", hashVector)
	t.Logf("outputWord : %v", outputWord)
	t.Logf("LDPC Nonce : %v", LDPCNonce)
}

//TestDecodingProcess test decoder.go functions
func TestOptimizedDecodingProcess(t *testing.T) {
	//tickerCounter := 0
	//ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	var LDPCNonce uint32
	var hashVector []int
	var outputWord []int
	//var LRrtl [][]float64

	header := ethHeader{}

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce string
	var encryptedHeaderWithNonce [32]byte

	parameters := SetDifficultyUsingLevel(0)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	for {
		startTime := time.Now()
		//	fmt.Printf("\rDecoding %s", ticker[tickerCounter])
		//	tickerCounter++
		//	tickerCounter %= len(ticker)

		//If Nonce is bigger than MaxNonce, then update timestamp
		if LDPCNonce >= MaxNonce {
			LDPCNonce = 0
			//header.Time = uint64(time.Now().Unix())
			//currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
		}
		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(uint64(LDPCNonce), 10)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ = OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(parameters, colInRow, outputWord)

		elapseTime := time.Since(startTime)

		if LDPCNonce%10000 == 0 {
			fmt.Printf("1 cycle decoding elapse Time : %v\n", elapseTime)
			fmt.Printf("hashVector : %v\n", hashVector)
			fmt.Printf("outputWord : %v\n", outputWord)
			fmt.Printf("LDPC Nonce : %v\n", LDPCNonce)
		}

		if flag {
			t.Logf("\nCodeword is founded with nonce = %v\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}
	t.Logf("hashVector : %v", hashVector)
	t.Logf("outputWord : %v", outputWord)
	t.Logf("LDPC Nonce : %v", LDPCNonce)
}

//TestRunLDPC test runLDPC function
func TestRunLDPC(t *testing.T) {
	parameters := SetDifficultyUsingLevel(0)
	var tempParentHash [32]byte
	//tempParentHash = [0, 0, ..., 0]
	parameters.seed = GenerateSeed(tempParentHash)

	tempHeader := ethHeader{}

	RunLDPC(parameters, tempHeader)
}

func TestVerifyDecoding(t *testing.T) {
	//parameters := SetDifficultyUsingLevel(0)

	parameters := Parameters{
		n:  32,
		wc: 3,
		wr: 8,
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)

	header := ethHeader{}
	copy(header.ParentHash[:], "00000000000000000000000000000000")
	parameters.seed = GenerateSeed(header.ParentHash)

	hashVector, outputWord, LDPCNonce := RunLDPC(parameters, header)
	verificationResult, hashVectorOfVerification, outputWordOfVerification := VerifyDecoding(parameters, outputWord, LDPCNonce, header)

	if !verificationResult {
		t.Error("Wrong outputwWord")
		t.Errorf("OutputWord of decoding     : %v", outputWord)
		t.Errorf("OutputWord of verification : %v", outputWordOfVerification)
		t.Errorf("HashVector of decoding     : %v", hashVector)
		t.Errorf("HashVector of verification : %v", hashVectorOfVerification)
	} else {
		t.Logf("OutputWord of decoding     : %v", outputWord)
		t.Logf("OutputWord of verification : %v", outputWordOfVerification)
		t.Logf("HashVector of decoding     : %v", hashVector)
		t.Logf("HashVector of verification : %v", hashVectorOfVerification)
	}
}

func TestVerifyOptimizedDecoding(t *testing.T) {

	parameters := Parameters{
		n:  32,
		wc: 3,
		wr: 8,
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)

	header := ethHeader{}
	copy(header.ParentHash[:], "00000000000000000000000000000000")
	parameters.seed = GenerateSeed(header.ParentHash)

	hashVector, outputWord, LDPCNonce := RunOptimizedLDPC(parameters, header)
	verificationResult, hashVectorOfVerification, outputWordOfVerification := VerifyOptimizedDecoding(parameters, outputWord, LDPCNonce, header)

	if !verificationResult {
		t.Error("Wrong outputwWord")
		t.Errorf("OutputWord of decoding     : %v", outputWord)
		t.Errorf("OutputWord of verification : %v", outputWordOfVerification)
		t.Errorf("HashVector of decoding     : %v", hashVector)
		t.Errorf("HashVector of verification : %v", hashVectorOfVerification)
	} else {
		t.Logf("OutputWord of decoding     : %v", outputWord)
		t.Logf("OutputWord of verification : %v", outputWordOfVerification)
		t.Logf("HashVector of decoding     : %v", hashVector)
		t.Logf("HashVector of verification : %v", hashVectorOfVerification)
	}
}

package ldpc

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"sync"
	"testing"
	"time"
)

// TestDecoding function is implemented to compare with original decoder(c++ version)
func TestDecoding(t *testing.T) {
	var opHashVector []int

	opParameters := SetDifficultyUsingLevel(0)

	opH := func(parameters Parameters) [][]int {
		var H [][]int
		var colOrder []int

		k := parameters.m / parameters.wc
		H = make([][]int, parameters.m)
		for i := range H {
			H[i] = make([]int, parameters.n)
		}

		for i := 0; i < k; i++ {
			for j := i * parameters.wr; j < (i+1)*parameters.wr; j++ {
				H[i][j] = 1
			}
		}

		for i := 1; i < parameters.wc; i++ {
			colOrder = nil
			for j := 0; j < parameters.n; j++ {
				colOrder = append(colOrder, j)
			}

			for j := 0; j < parameters.n; j++ {
				index := (colOrder[j]/parameters.wr + k*i)
				H[index][j] = 1
			}
		}
		return H
	}(opParameters)

	opColInRow, opRowInCol := GenerateQ(opParameters, opH)
	opEncryptedHeaderWithNonce := sha256.Sum256([]byte("0"))

	opHashVector = GenerateHv(opParameters, opEncryptedHeaderWithNonce)

	opHashVector, opOutputWord, _ := OptimizedDecoding(opParameters, opHashVector, opH, opRowInCol, opColInRow)

	t.Logf("OptimizedDecoder opHashVector : %v\n", opHashVector)
	t.Logf("OptimezedDecoder outputWord	: %v\n", opOutputWord)
}
func TestConcurrencyPerformance(t *testing.T) {
	header := ethHeader{}
	parameters := SetDifficultyUsingLevel(0)

	var wg sync.WaitGroup
	var outerLoopSignal = make(chan struct{})
	var innerLoopSignal = make(chan struct{})
	var goRoutineSignal = make(chan struct{})

	attemptCount := 0

outerLoop:
	for {
		select {
		// If outerLoopSignal channel is closed, then break outerLoop
		case <-outerLoopSignal:
			break outerLoop

		default:
			// Defined default to unblock select statement
		}

	innerLoop:
		for i := 0; i < runtime.NumCPU(); i++ {
			select {
			// If innerLoop signal is closed, then break innerLoop and close outerLoopSignal
			case <-innerLoopSignal:
				close(outerLoopSignal)
				break innerLoop

			default:
				// Defined default to unblock select statement
			}

			wg.Add(1)
			go func(goRoutineSignal chan struct{}) {
				defer wg.Done()
				//goRoutineNonce := generateRandomNonce()
				//fmt.Printf("Initial goroutine Nonce : %v\n", goRoutineNonce)

				var goRoutineHashVector []int

				var serializedHeader = string(header.ParentHash[:])
				var serializedHeaderWithNonce string
				var encryptedHeaderWithNonce [32]byte
				H := GenerateH(parameters)
				colInRow, rowInCol := GenerateQ(parameters, H)

				select {
				case <-goRoutineSignal:
					break

				default:
				attemptLoop:
					for attempt := 0; attempt < 5000; attempt++ {
						attemptCount++
						goRoutineNonce := generateRandomNonce()
						serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(goRoutineNonce, 10)
						encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

						goRoutineHashVector = GenerateHv(parameters, encryptedHeaderWithNonce)
						goRoutineHashVector, _, _ = OptimizedDecoding(parameters, goRoutineHashVector, H, rowInCol, colInRow)

						select {
						case <-goRoutineSignal:
							// fmt.Println("goRoutineSignal channel is already closed")
							break attemptLoop
						default:
							if attemptCount == 1000000 {
								close(goRoutineSignal)
								close(innerLoopSignal)
								fmt.Printf("Codeword is founded with nonce = %d\n", goRoutineNonce)
								break attemptLoop
							}
						}
						//goRoutineNonce++
					}
				}
			}(goRoutineSignal)
		}
		// Need to wait to prevent memory leak
		wg.Wait()
	}
}

func TestOptimizedDecodingImplement(t *testing.T) {
	var nonce uint64
	for i := 0; i < 100000; i++ {
		header := ethHeader{}

		var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
		var serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(nonce, 10)
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
		var opSerializedHeaderWithNonce = opSerializedHeader + strconv.FormatUint(nonce, 10)
		var opEncryptedHeaderWithNonce [32]byte

		var opHashVector []int

		opParameters := SetDifficultyUsingLevel(0)
		opParameters.seed = GenerateSeed(opHeader.ParentHash)

		opH := GenerateH(opParameters)
		opColInRow, opRowInCol := GenerateQ(opParameters, opH)
		opEncryptedHeaderWithNonce = sha256.Sum256([]byte(opSerializedHeaderWithNonce))

		opHashVector = GenerateHv(opParameters, opEncryptedHeaderWithNonce)

		opHashVector, opOutputWord, _ := OptimizedDecoding(opParameters, opHashVector, opH, opRowInCol, opColInRow)

		if !reflect.DeepEqual(hashVector, opHashVector) || !reflect.DeepEqual(outputWord, opOutputWord) {
			t.Errorf("Decoder hashVector		  :  %v\n", hashVector)
			t.Errorf("OptimezedDecoder hashVector : %v\n", opHashVector)

			t.Errorf("Decoder outputWord		  :  %v\n", outputWord)
			t.Errorf("OptimezedDecoder outputWord : %v\n", opOutputWord)
		}

		nonce++
		/*
			t.Logf("Decoder hashVector			: %v\n", hashVector)
			t.Logf("OptimezedDecoder hashVector : %v\n", opHashVector)

			t.Logf("Decoder outputWord			: %v\n", outputWord)
			t.Logf("OptimezedDecoder outputWord	: %v\n", opOutputWord)
		*/
	}
}

func TestDecodingElapseTime(t *testing.T) {
	header := ethHeader{}

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce = serializedHeader + ""
	var encryptedHeaderWithNonce [32]byte

	var hashVector []int
	//var LRrtl [][]float64

	parameters := SetDifficultyUsingLevel(2)
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

	parameters := SetDifficultyUsingLevel(0)
	parameters.seed = GenerateSeed(header.ParentHash)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)
	encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)

	for i := 0; i < 1000000; i++ {
		_, _, _ = OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)
	}
}

//TestDecodingProcess test decoder.go functions
func TestDecodingProcess(t *testing.T) {
	//tickerCounter := 0
	//ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	LDPCNonce := generateRandomNonce()
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

		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(LDPCNonce, 10)
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

	LDPCNonce := generateRandomNonce()
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

		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(LDPCNonce, 10)
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

func TestRunOptimizedLDPC(t *testing.T) {
	parameters := SetDifficultyUsingLevel(1)
	var tempParentHash [32]byte
	//tempParentHash = [0, 0, ..., 0]
	parameters.seed = GenerateSeed(tempParentHash)

	tempHeader := ethHeader{}

	RunOptimizedLDPC(parameters, tempHeader)
}

func TestRunOptimizedConcurrencyLDPC(t *testing.T) {
	parameters := SetDifficultyUsingLevel(0)
	var tempParentHash [32]byte
	//tempParentHash = [0, 0, ..., 0]
	parameters.seed = GenerateSeed(tempParentHash)

	tempHeader := ethHeader{}

	RunOptimizedConcurrencyLDPC(parameters, tempHeader)
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

func TestVerifyConcurrencyDecoding(t *testing.T) {

	parameters := Parameters{
		n:  32,
		wc: 3,
		wr: 8,
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)

	header := ethHeader{}
	copy(header.ParentHash[:], "00000000000000000000000000000000")
	parameters.seed = GenerateSeed(header.ParentHash)

	hashVector, outputWord, LDPCNonce := RunOptimizedConcurrencyLDPC(parameters, header)
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

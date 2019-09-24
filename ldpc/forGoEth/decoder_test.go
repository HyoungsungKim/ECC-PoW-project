package ldpc

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/Onther-Tech/go-ethereum/core/types"
)

func TestNonceDecoding(t *testing.T) {
	LDPCNonce := generateRandomNonce()
	EncodedNonce := types.EncodeNonce(LDPCNonce)
	DecodedNonce := EncodedNonce.Uint64()

	if LDPCNonce == DecodedNonce {
		t.Logf("LDPCNonce : %v\n", LDPCNonce)
		t.Logf("Decoded Nonce : %v\n", DecodedNonce)
	} else {
		t.Errorf("LDPCNonce : %v\n", LDPCNonce)
		t.Errorf("Decoded Nonce : %v\n", DecodedNonce)
	}
}

func TestOptimizedDecodingElapseTime(t *testing.T) {
	var block ethHeader
	var serializedHeader = string(block.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce = serializedHeader + ""
	var encryptedHeaderWithNonce [32]byte

	var hashVector []int
	//var LRrtl [][]float64

	block.Difficulty = big.NewInt(0)
	parameters, _ := SetParameters(block)
	parameters.seed = generateSeed(block.ParentHash)

	H := generateH(parameters)
	colInRow, rowInCol := generateQ(parameters, H)
	encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector = generateHv(parameters, encryptedHeaderWithNonce)

	for i := 0; i < 10000; i++ {
		_, _, _ = OptimizedDecoding(block, hashVector, H, rowInCol, colInRow)
	}
}

//TestDecodingProcess test decoder.go functions
func TestOptimizedDecodingProcess(t *testing.T) {
	//tickerCounter := 0
	//ticker := []string{"-", "-", "\\", "\\", "|", "|", "/", "/"}

	LDPCNonce := generateRandomNonce()
	var hashVector []int
	var outputWord []int
	//var LRrtl [][]float64

	block := ethHeader{}
	block.Difficulty = big.NewInt(0)

	var serializedHeader = string(block.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce string
	var encryptedHeaderWithNonce [32]byte

	parameters, _ := SetParameters(block)
	parameters.seed = generateSeed(block.ParentHash)

	H := generateH(parameters)
	colInRow, rowInCol := generateQ(parameters, H)

	for {
		startTime := time.Now()

		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(LDPCNonce, 10)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = generateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ = OptimizedDecoding(block, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(block, colInRow, outputWord)

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

func TestRunOptimizedLDPC(t *testing.T) {
	block := ethHeader{}
	block.Difficulty = big.NewInt(0)
	RunOptimizedLDPC(block)
}

func TestRunOptimizedConcurrencyLDPC(t *testing.T) {
	block := ethHeader{}
	block.Difficulty = big.NewInt(0)
	RunOptimizedConcurrencyLDPC(block)
}

func TestVerifyOptimizedDecoding(t *testing.T) {
	block := ethHeader{}
	block.Difficulty = big.NewInt(0)

	hashVector, outputWord, LDPCNonce := RunOptimizedLDPC(block)
	verificationResult, hashVectorOfVerification, outputWordOfVerification := VerifyOptimizedDecoding(block, LDPCNonce)

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
	block := ethHeader{}
	block.Difficulty = big.NewInt(0)

	hashVector, outputWord, LDPCNonce := RunOptimizedConcurrencyLDPC(block)
	verificationResult, hashVectorOfVerification, outputWordOfVerification := VerifyOptimizedDecoding(block, LDPCNonce)

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

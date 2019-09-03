package ldpc

import (
	"crypto/sha256"
	"io"
	"strconv"
	"testing"

	"github.com/ethereum/go-ethereum/rlp"
)

//https://github.com/ethereum/go-ethereum/blob/master/rlp/encode_test.go
//https://github.com/ethereum/go-ethereum/blob/master/rlp/decode_test.go
//Read and implement encoding and decoding test again

type MyCoolType struct {
	Name string
	a, b uint
}

func (x *MyCoolType) EncodeRLP(w io.Writer) (err error) {
	if x == nil {
		err = rlp.Encode(w, []int{0, 0})
	} else {
		err = rlp.Encode(w, []uint{x.a, x.b})
	}

	return err
}

func TestEncoder(t *testing.T) {
	var m *MyCoolType
	bytes, _ := rlp.EncodeToBytes(m)
	t.Logf("%v -> %X\n", m, bytes)

	m = &MyCoolType{Name: "foobar", a: 5, b: 6}
	bytes, _ = rlp.EncodeToBytes(m)
	t.Logf("%v -> %X", m, bytes)
}

/*
//EncodeRLP implementation from https://godoc.org/github.com/ethereum/go-ethereum/rlp#example-Encoder
func (en *extraNonce) EncodeRLP(w io.Writer) (err error) {
	if en == nil {
		err = rlp.Encode(w, extraNonce{0, "0", 0})
	} else {
		err = rlp.Encode(w, extraNonce{en.difficulty, en.outputWord, en.LDPCNonce})
	}
	return err
}

func (en *extraNonce) DecodeRLP(r io.Reader) (err error) {
	if en == nil {
		err = rlp.Decode(r, extraNonce{0, "0", 0})
	} else {
		err = rlp.Decode(r, extraNonce{en.difficulty, en.outputWord, en.LDPCNonce})
	}
	return err
}

func TestRLPEncoding(t *testing.T) {
	var en *extraNonce
	encodingResult, _ := rlp.EncodeToBytes(en)
	t.Logf("%v -> %X\n", en, encodingResult)

	en = &extraNonce{0, "123456789", 0}
	encodingResult, _ = rlp.EncodeToBytes(en)
	t.Logf("%v -> %X\n", en, encodingResult)
}

func TestRLPDecoding(t *testing.T) {
	var result *extraNonce
	en := &extraNonce{
		difficulty: 0,
		outputWord: "123456789",
		LDPCNonce:  0,
	}
	encodingResult, _ := rlp.EncodeToBytes(en)
	t.Logf("Encoding Result : %v -> %X\n", en, encodingResult)

	err := rlp.Decode(bytes.NewReader(encodingResult), &result)
	if err != nil {
		t.Errorf("Error : %v\n", err)
	}

	if en == result {
		t.Logf("Before encoding : %v\n", en)
		t.Logf("Encoding Result : %X\n", encodingResult)
		t.Logf("After decoding : %v\n", result)
	} else {
		t.Errorf("Before encoding : %v\n", en)
		t.Errorf("Encoding Result : %X\n", encodingResult)
		t.Errorf("After decoding : %v\n", result)
	}

}
*/

func TestDecodingElapseTime(t *testing.T) {
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
	hashVector, _, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)

	for i := 0; i < 200000; i++ {
		_, _, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
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
	hashVector, _, _ = OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)

	for i := 0; i < 200000; i++ {
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

	//	parameters := SetDifficultyUsingLevel(0)
	//	parameters.seed = GenerateSeed(header.ParentHash)

	parameters := Parameters{
		n:  24,
		wc: 3,
		wr: 8,
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	for {
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

	//	parameters := SetDifficultyUsingLevel(0)
	//	parameters.seed = GenerateSeed(header.ParentHash)
	parameters := Parameters{
		n:  24,
		wc: 3,
		wr: 8,
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)

	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	for {
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

		if !flag {
			hashVector, outputWord, _ = OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)
			flag = MakeDecision(parameters, colInRow, outputWord)
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

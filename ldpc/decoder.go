package ldpc

import (
	crand "crypto/rand"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
)

//GenerateRandomNonce generate 64bit random nonce with similar way of ethereum block nonce
func generateRandomNonce() uint64 {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	source := rand.New(rand.NewSource(seed.Int64()))

	return uint64(source.Int63())
}

//RunLDPC function needs more concrete implementation
//return hashVector, outputWord, LDPCNonce
func RunLDPC(parameters Parameters, header ethHeader) ([]int, []int, uint64) {
	//Need to set difficulty before running LDPC
	LDPCNonce := generateRandomNonce()
	var hashVector []int
	var outputWord []int
	//	var LRrtl [][]float64

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce string
	var encryptedHeaderWithNonce [32]byte
	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	for {
		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(LDPCNonce, 10)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ = Decoding(parameters, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(parameters, colInRow, outputWord)

		if flag {
			fmt.Printf("Codeword is founded with nonce = %d\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}

	return hashVector, outputWord, LDPCNonce
}

//RunOptimizedLDPC use OptimizedDecoding function not decoding function
func RunOptimizedLDPC(parameters Parameters, header ethHeader) ([]int, []int, uint64) {
	//Need to set difficulty before running LDPC
	LDPCNonce := generateRandomNonce()
	var hashVector []int
	var outputWord []int
	//	var LRrtl [][]float64

	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	var serializedHeaderWithNonce string
	var encryptedHeaderWithNonce [32]byte
	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)

	for {
		serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(LDPCNonce, 10)
		encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

		hashVector = GenerateHv(parameters, encryptedHeaderWithNonce)
		hashVector, outputWord, _ = OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)
		flag := MakeDecision(parameters, colInRow, outputWord)

		if flag {
			fmt.Printf("Codeword is founded with nonce = %d\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}

	return hashVector, outputWord, LDPCNonce
}

//SetDifficultyUsingLevel set matrix parameters
//Only 4 parameters are valied 0 : Very easy, 1 : Easy, 2 : Medium, 3 : hard
func SetDifficultyUsingLevel(level int) Parameters {
	//level 4 is max level
	if level > 4 {
		level = 4
	}

	parameters := Parameters{}
	if level == 0 {
		parameters.n = 16
		parameters.wc = 3
		parameters.wr = 4
	} else if level == 1 {
		parameters.n = 32
		parameters.wc = 3
		parameters.wr = 4
	} else if level == 2 {
		parameters.n = 64
		parameters.wc = 3
		parameters.wr = 4
	} else if level == 3 {
		parameters.n = 128
		parameters.wc = 3
		parameters.wr = 4
	}
	parameters.m = int(parameters.n * parameters.wc / parameters.wr)

	return parameters
}

//GenerateSeed generate seed using previous hash vector
func GenerateSeed(phv [32]byte) int {
	sum := 0
	for i := 0; i < len(phv); i++ {
		sum += int(phv[i])
	}
	return sum
}

//VerifyDecoding return bool, hashVector of verification, outputWord of verification
func VerifyDecoding(parameters Parameters, outputWord []int, LDPCNonce uint64, header ethHeader) (bool, []int, []int) {
	//VerifyDecoding function needs more concrete implementation
	//It has to be decided
	//It is right to generate H, colInRow, rowInCol or pass these using other ways
	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)
	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	serializedHeaderWithNonce := serializedHeader + strconv.FormatUint(LDPCNonce, 10)
	encryptedHeaderWithNonce := sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector := GenerateHv(parameters, encryptedHeaderWithNonce)
	hashVectorOfVerification, outputWordOfVerification, _ := Decoding(parameters, hashVector, H, rowInCol, colInRow)

	if reflect.DeepEqual(outputWord, outputWordOfVerification) && reflect.DeepEqual(hashVector, hashVectorOfVerification) {
		return true, hashVectorOfVerification, outputWordOfVerification
	}

	return false, hashVectorOfVerification, outputWordOfVerification
}

//VerifyOptimizedDecoding return bool, hashVector of verification, outputWord of verification
func VerifyOptimizedDecoding(parameters Parameters, outputWord []int, LDPCNonce uint64, header ethHeader) (bool, []int, []int) {
	//VerifyOptimizedDecoding function needs more concrete implementation
	//It has to be decided
	//It is right to generate H, colInRow, rowInCol or pass these using other ways
	H := GenerateH(parameters)
	colInRow, rowInCol := GenerateQ(parameters, H)
	var serializedHeader = string(header.ParentHash[:]) // + ... + string(header.MixDigest)
	serializedHeaderWithNonce := serializedHeader + strconv.FormatUint(LDPCNonce, 10)
	encryptedHeaderWithNonce := sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector := GenerateHv(parameters, encryptedHeaderWithNonce)
	hashVectorOfVerification, outputWordOfVerification, _ := OptimizedDecoding(parameters, hashVector, H, rowInCol, colInRow)

	if reflect.DeepEqual(outputWord, outputWordOfVerification) && reflect.DeepEqual(hashVector, hashVectorOfVerification) {
		return true, hashVectorOfVerification, outputWordOfVerification
	}

	return false, hashVectorOfVerification, outputWordOfVerification
}

//Decoding carry out LDPC decoding and return hashvector, outputWord, LRrtl
func Decoding(parameters Parameters,
	hashVector []int,
	H, rowInCol, colInRow [][]int,
) ([]int, []int, [][]float64) {
	var temp3, tempSign, sign, magnitude float64

	outputWord := make([]int, parameters.n)
	LRqtl := make([][]float64, parameters.n)
	LRrtl := make([][]float64, parameters.n)
	LRft := make([]float64, parameters.n)

	for i := 0; i < parameters.n; i++ {
		LRqtl[i] = make([]float64, parameters.m)
		LRrtl[i] = make([]float64, parameters.m)
		LRft[i] = math.Log((1-crossErr)/crossErr) * float64((hashVector[i]*2 - 1))
	}
	LRpt := make([]float64, parameters.n)

	var i, k, l, m, ind, t, mp int
	for ind = 1; ind <= maxIter; ind++ {
		for t = 0; t < parameters.n; t++ {
			for m = 0; m < parameters.wc; m++ {
				temp3 = 0
				for mp = 0; mp < parameters.wc; mp++ {
					if mp != m {
						temp3 = infinityTest(temp3 + LRrtl[t][rowInCol[mp][t]])
					}
				}
				LRqtl[t][rowInCol[m][t]] = infinityTest(LRft[t] + temp3)
			}
		}
		for k = 0; k < m; k++ {
			for l = 0; l < parameters.wr; l++ {
				temp3 = 0.0
				sign = 1
				for m = 0; m < parameters.wr; m++ {
					if m != l {
						temp3 = temp3 + funcF(math.Abs(LRqtl[colInRow[m][k]][k]))
						if LRqtl[colInRow[m][k]][k] > 0.0 {
							tempSign = 1.0
						} else {
							tempSign = -1.0
						}
						sign = sign * tempSign
					}
				}
				magnitude = funcF(temp3)
				LRrtl[colInRow[l][k]][k] = infinityTest(sign * magnitude)
			}
		}
		for m = 0; m < parameters.n; m++ {
			LRpt[m] = infinityTest(LRft[m])
			for k = 0; k < parameters.wc; k++ {
				LRpt[m] += LRrtl[m][rowInCol[k][m]]
				LRpt[m] = infinityTest(LRpt[m])
			}
		}
	}
	for i = 0; i < parameters.n; i++ {
		if LRpt[i] >= 0 {
			outputWord[i] = 1
		} else {
			outputWord[i] = 0
		}
	}

	return hashVector, outputWord, LRrtl
}

//OptimizedDecoding is 20% faster than previous decoding function when they use same nonce
//percentage can be changed because of random nonce
func OptimizedDecoding(parameters Parameters,
	hashVector []int,
	H, rowInCol, colInRow [][]int,
) ([]int, []int, [][]float64) {

	outputWord := make([]int, parameters.n)
	LRqtl := make([][]float64, parameters.n)
	LRrtl := make([][]float64, parameters.n)
	LRft := make([]float64, parameters.n)

	for i := 0; i < parameters.n; i++ {
		LRqtl[i] = make([]float64, parameters.m)
		LRrtl[i] = make([]float64, parameters.m)
		LRft[i] = math.Log((1-crossErr)/crossErr) * float64((hashVector[i]*2 - 1))
	}
	LRpt := make([]float64, parameters.n)

	for ind := 1; ind <= maxIter; ind++ {
		for t := 0; t < parameters.n; t++ {
			temp3 := 0.0

			for mp := 0; mp < parameters.wc; mp++ {
				temp3 = infinityTest(temp3 + LRrtl[t][rowInCol[mp][t]])
			}
			for m := 0; m < parameters.wc; m++ {
				temp4 := temp3
				temp4 = infinityTest(temp4 - LRrtl[t][rowInCol[m][t]])
				LRqtl[t][rowInCol[m][t]] = infinityTest(LRft[t] + temp4)
			}
		}

		for k := 0; k < parameters.wr; k++ {
			for l := 0; l < parameters.wr; l++ {
				temp3 := 0.0
				sign := 1.0
				tempSign := 0.0
				for m := 0; m < parameters.wr; m++ {
					if m != l {
						temp3 = temp3 + funcF(math.Abs(LRqtl[colInRow[m][k]][k]))
						if LRqtl[colInRow[m][k]][k] > 0.0 {
							tempSign = 1.0
						} else {
							tempSign = -1.0
						}
						sign = sign * tempSign
					}
				}
				magnitude := funcF(temp3)
				LRrtl[colInRow[l][k]][k] = infinityTest(sign * magnitude)
			}
		}

		for t := 0; t < parameters.n; t++ {
			LRpt[t] = infinityTest(LRft[t])
			for k := 0; k < parameters.wc; k++ {
				LRpt[t] += LRrtl[t][rowInCol[k][t]]
				LRpt[t] = infinityTest(LRpt[t])
			}

			if LRpt[t] >= 0 {
				outputWord[t] = 1
			} else {
				outputWord[t] = 0
			}
		}
	}

	return hashVector, outputWord, LRrtl
}

//GenerateH generate H matrix using parameters
//GenerateH Cannot be sure rand is same with original implementation of C++
func GenerateH(parameters Parameters) [][]int {
	var H [][]int
	var hSeed int64
	hSeed = int64(parameters.seed)

	var colOrder []int
	/*
		if H == nil {
			return false
		}
	*/
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

		rand.Seed(hSeed)
		rand.Shuffle(len(colOrder), func(i, j int) {
			colOrder[i], colOrder[j] = colOrder[j], colOrder[i]
		})
		hSeed--

		for j := 0; j < parameters.n; j++ {
			index := (colOrder[j]/parameters.wr + k*i)
			H[index][j] = 1
		}
	}

	return H
}

//GenerateQ generate colInRow and rowInCol matrix using H matrix
func GenerateQ(parameters Parameters, H [][]int) ([][]int, [][]int) {
	colInRow := make([][]int, parameters.wr)
	for i := 0; i < parameters.wr; i++ {
		colInRow[i] = make([]int, parameters.m)
	}

	rowInCol := make([][]int, parameters.wc)
	for i := 0; i < parameters.wc; i++ {
		rowInCol[i] = make([]int, parameters.n)
	}

	rowIndex := 0
	colIndex := 0

	for i := 0; i < parameters.m; i++ {
		for j := 0; j < parameters.n; j++ {
			if H[i][j] == 1 {
				colInRow[colIndex%parameters.wr][i] = j
				colIndex++

				rowInCol[rowIndex/parameters.n][j] = i
				rowIndex++
			}
		}
	}

	return colInRow, rowInCol
}

//GenerateHv generate hashvector
//It needs to compare with origin C++ implementation Especially when sha256 function is used
func GenerateHv(parameters Parameters, encryptedHeaderWithNonce [32]byte) []int {
	//inputSize := len(headerWithNonce)
	hashVector := make([]int, parameters.n)

	/*
		if parameters.n <= 256 {
			tmpHashVector = sha256.Sum256(headerWithNonce)
		} else {
			/*
				This section is for a case in which the size of a hash vector is larger than 256.
				This section will be implemented soon.
		}
			transform the constructed hexadecimal array into an binary array
			ex) FE01 => 11111110000 0001
	*/

	for i := 0; i < parameters.n/8; i++ {
		decimal := int(encryptedHeaderWithNonce[i])
		for j := 7; j >= 0; j-- {
			hashVector[j+8*(i)] = decimal % 2
			decimal /= 2
		}
	}

	//outputWord := hashVector[:parameters.n]
	return hashVector
}

/*
//isRegular check parameters that are valid for matrix
func isRegular(nSize, wCol, wRow int) bool {
	res := float64(nSize*wCol) / float64(wRow)
	m := math.Round(res)

	if int(m)*wRow == nSize*wCol {
		return true
	}

	return false
}

//SetDifficulty sets LDPC parameters using function parameters
//If function parameters are not valied then return err
func SetDifficulty(nSize, wCol, wRow int) (Parameters, error) {
	parameters := Parameters{}
	if isRegular(nSize, wCol, wRow) {
		parameters.n = nSize
		parameters.wc = wCol
		parameters.wr = wRow
		parameters.m = int(parameters.n * parameters.wc / parameters.wr)
		return parameters, nil
	}
	return parameters, errors.New("Wrong function parameters")
}
*/

//MakeDecision check outputWord is valid or not using colInRow
func MakeDecision(parameters Parameters, colInRow [][]int, outputWord []int) bool {
	for i := 0; i < parameters.m; i++ {
		sum := 0
		for j := 0; j < parameters.wr; j++ {
			//	fmt.Printf("i : %d, j : %d, m : %d, wr : %d \n", i, j, m, wr)
			sum = sum + outputWord[colInRow[j][i]]
		}
		if sum%2 == 1 {
			return false
		}
	}
	return true
}

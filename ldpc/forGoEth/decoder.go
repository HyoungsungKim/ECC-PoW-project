package ldpc

import (
	crand "crypto/rand"
	"crypto/sha256"
	"math"
	"math/big"
	"math/rand"
	"strconv"
)

//GenerateRandomNonce generate 64bit random nonce with similar way of ethereum block nonce
func generateRandomNonce() uint64 {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	source := rand.New(rand.NewSource(seed.Int64()))

	return uint64(source.Int63())
}

//GenerateSeed generate seed using previous hash vector
func generateSeed(phv [32]byte) int {
	sum := 0
	for i := 0; i < len(phv); i++ {
		sum += int(phv[i])
	}
	return sum
}

//VerifyOptimizedDecoding return bool, hashVector of verification, outputWord of verification
func VerifyOptimizedDecoding(block ethHeader, LDPCNonce uint64) (bool, []int, []int) {
	parameters, _ := SetParameters(block)

	H := generateH(parameters)
	colInRow, rowInCol := generateQ(parameters, H)

	var serializedHeader = string(block.ParentHash[:]) // + ... + string(header.MixDigest)
	serializedHeaderWithNonce := serializedHeader + strconv.FormatUint(LDPCNonce, 10)
	encryptedHeaderWithNonce := sha256.Sum256([]byte(serializedHeaderWithNonce))

	hashVector := generateHv(parameters, encryptedHeaderWithNonce)
	hashVectorOfVerification, outputWordOfVerification, _ := OptimizedDecoding(block, hashVector, H, rowInCol, colInRow)

	if MakeDecision(block, colInRow, outputWordOfVerification) {
		return true, hashVectorOfVerification, outputWordOfVerification
	}

	return false, hashVectorOfVerification, outputWordOfVerification
}

//OptimizedDecoding is 20% faster than previous decoding function when they use same nonce
//percentage can be changed because of random nonce
func OptimizedDecoding(block ethHeader, hashVector []int, H, rowInCol, colInRow [][]int) ([]int, []int, [][]float64) {

	parameters, _ := SetParameters(block)
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
func generateH(parameters Parameters) [][]int {
	var H [][]int
	var hSeed int64
	var colOrder []int

	hSeed = int64(parameters.seed)
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
func generateQ(parameters Parameters, H [][]int) ([][]int, [][]int) {
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
func generateHv(parameters Parameters, encryptedHeaderWithNonce [32]byte) []int {
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

//MakeDecision check outputWord is valid or not using colInRow
func MakeDecision(block ethHeader, colInRow [][]int, outputWord []int) bool {
	parameters, difficultyLevel := SetParameters(block)
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

	var numOfOnes int
	for _, val := range outputWord {
		numOfOnes += val
	}

	if numOfOnes >= Table[difficultyLevel].decisionFrom &&
		numOfOnes <= Table[difficultyLevel].decisionTo &&
		numOfOnes%Table[difficultyLevel].decisionStep == 0 {
		return true
	}

	return false
}

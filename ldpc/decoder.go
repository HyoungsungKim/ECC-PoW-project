package ldpc

import (
	"crypto/sha256"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

//runLDPC function needs more concrete implementation
func runLDPC(header ethHeader) {
	//Need to set difficulty before running LDPC
	LDPCNonce = 0

	header.Time = uint64(time.Now().Unix())
	var currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
	var currentBlockHeaderWithNonce string

	GenerateSeed(header.ParentHash)
	GenerateH()
	GenerateQ()

	for {
		//If Nonce is bigger than MaxNonce, then update timestamp
		if LDPCNonce >= MaxNonce {
			LDPCNonce = 0
			header.Time = uint64(time.Now().Unix())
			currentBlockHeader = string(header.ParentHash[:]) + strconv.FormatUint(header.Time, 10)
		}
		currentBlockHeaderWithNonce = currentBlockHeader + fmt.Sprint(LDPCNonce)

		GenerateHv([]byte(currentBlockHeaderWithNonce))
		Decoding()
		flag := Decision()

		if !flag {
			Decoding()
			flag = Decision()
		}
		if flag {
			fmt.Printf("\nCodeword is founded with nonce = %d\n", LDPCNonce)
			break
		}
		LDPCNonce++
	}
}

func Decoding() {
	var temp3, tempSign, sign, magnitude float64

	outputWord = make([]int, n)
	LRqtl = make([][]float64, n)
	LRrtl = make([][]float64, n)
	LRft = make([]float64, n)

	for i := 0; i < n; i++ {
		LRqtl[i] = make([]float64, m)
		LRrtl[i] = make([]float64, m)
		LRft[i] = math.Log((1-crossErr)/crossErr) * float64((hashVector[i]*2 - 1))
	}
	LRpt = make([]float64, n)

	var i, k, l, m, ind, t, mp int
	for ind = 1; ind <= maxIter; ind++ {
		for t = 0; t < n; t++ {
			for m = 0; m < wc; m++ {
				temp3 = 0
				for mp = 0; mp < wc; mp++ {
					if mp != m {
						temp3 = infinityTest(temp3 + LRrtl[t][rowInCol[mp][t]])
					}
				}
				LRqtl[t][rowInCol[m][t]] = infinityTest(LRft[t] + temp3)
			}
		}
		for k = 0; k < m; k++ {
			for l = 0; l < wr; l++ {
				temp3 = 0.0
				sign = 1
				for m = 0; m < wr; m++ {
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
		for m = 0; m < n; m++ {
			LRpt[m] = infinityTest(LRft[m])
			for k = 0; k < wc; k++ {
				LRpt[m] += LRrtl[m][rowInCol[k][m]]
				LRpt[m] = infinityTest(LRpt[m])
			}
		}
	}
	for i = 0; i < n; i++ {
		if LRpt[i] >= 0 {
			outputWord[i] = 1
		} else {
			outputWord[i] = 0
		}
	}
}

func GenerateSeed(phv [32]byte) int {
	sum := 0
	for i := 0; i < len(phv); i++ {
		sum += int(phv[i])
	}
	seed = sum
	return sum
}

//GenerateH Cannot be sure rand is same with original implementation of C++
func GenerateH() bool {
	var hSeed int64
	hSeed = int64(seed)

	var colOrder []int
	/*
		if H == nil {
			return false
		}
	*/
	k := m / wc
	H = make([][]int, m)
	for i := range H {
		H[i] = make([]int, n)
	}

	for i := 0; i < k; i++ {
		for j := i * wr; j < (i+1)*wr; j++ {
			H[i][j] = 1
		}
	}

	for i := 1; i < wc; i++ {
		colOrder = nil
		for j := 0; j < n; j++ {
			colOrder = append(colOrder, j)
		}

		rand.Seed(hSeed)
		rand.Shuffle(len(colOrder), func(i, j int) {
			colOrder[i], colOrder[j] = colOrder[j], colOrder[i]
		})
		hSeed--

		for j := 0; j < n; j++ {
			index := (colOrder[j]/wr + k*i)
			H[index][j] = 1
		}
	}
	return true
}

func GenerateQ() bool {
	colInRow = make([][]int, wr)
	for i := 0; i < wr; i++ {
		colInRow[i] = make([]int, m)
	}

	rowInCol = make([][]int, wc)
	for i := 0; i < wc; i++ {
		rowInCol[i] = make([]int, n)
	}

	rowIndex := 0
	colIndex := 0

	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if H[i][j] == 1 {
				colInRow[colIndex%wr][i] = j
				colIndex++

				rowInCol[rowIndex/n][j] = i
				rowIndex++
			}
		}
	}
	return true
}

//GenerateHv need to compare with origin C++ implementation
//Especially when sha256 function is used
func GenerateHv(headerWithNonce []byte) {
	//inputSize := len(headerWithNonce)
	hashVector = make([]int, n)

	if n <= 256 {
		tmpHashVector = sha256.Sum256(headerWithNonce)
	} else {
		/*
			This section is for a case in which the size of a hash vector is larger than 256.
			This section will be implemented soon.
		*/
	}

	/*
		transform the constructed hexadecimal array into an binary arry
		ex) FE01 => 11111110000 0001
	*/
	for i := 0; i < n/8; i++ {
		decimal := int(tmpHashVector[i])
		for j := 7; j >= 0; j-- {
			hashVector[j+8*(i)] = decimal % 2
			decimal /= 2
		}
	}

	outputWord = hashVector[:n]
}

func isRegular(nSize, wCol, wRow int) bool {
	res := float64(nSize*wCol) / float64(wRow)
	m := math.Round(res)

	if int(m)*wRow == nSize*wCol {
		return true
	}

	return false
}

func SetDifficulty(nSize, wCol, wRow int) bool {
	if isRegular(nSize, wCol, wRow) {
		n = nSize
		wc = wCol
		wr = wRow
		m = int(n * wc / wr)
		return true
	}
	return false
}

//SetDifficultyUsingLevel 0 : Very easy, 1 : Easy, 2 : Medium, 3 : hard
func SetDifficultyUsingLevel(level int) {
	if level == 0 {
		n = 16
		wc = 3
		wr = 4
	} else if level == 1 {
		n = 32
		wc = 3
		wr = 4
	} else if level == 2 {
		n = 64
		wc = 3
		wr = 4
	} else if level == 3 {
		n = 128
		wc = 3
		wr = 4
	}
	m = int(n * wc / wr)
}

func Decision() bool {
	for i := 0; i < m; i++ {
		sum := 0
		for j := 0; j < wr; j++ {
			//	fmt.Printf("i : %d, j : %d, m : %d, wr : %d \n", i, j, m, wr)
			sum = sum + outputWord[colInRow[j][i]]
		}
		if sum%2 == 1 {
			return false
		}
	}
	return true
}

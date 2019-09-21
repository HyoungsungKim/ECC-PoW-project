package ldpc

import (
	"crypto/sha256"
	"fmt"
	"runtime"
	"strconv"
	"sync"
)

//RunOptimizedConcurrencyLDPC use goroutine for mining block
func RunOptimizedConcurrencyLDPC(block ethHeader) ([]int, []int, uint64) {
	//Need to set difficulty before running LDPC
	// Number of goroutines : 500, Number of attempts : 50000 Not bad

	var LDPCNonce uint64
	var hashVector []int
	var outputWord []int

	var wg sync.WaitGroup
	var outerLoopSignal = make(chan struct{})
	var innerLoopSignal = make(chan struct{})
	var goRoutineSignal = make(chan struct{})
	parameters, _ := SetParameters(block)

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
				var goRoutineOutputWord []int

				var serializedHeader = string(block.ParentHash[:])
				var serializedHeaderWithNonce string
				var encryptedHeaderWithNonce [32]byte
				H := generateH(parameters)
				colInRow, rowInCol := generateQ(parameters, H)

				select {
				case <-goRoutineSignal:
					break

				default:
				attemptLoop:
					for attempt := 0; attempt < 5000; attempt++ {
						goRoutineNonce := generateRandomNonce()
						serializedHeaderWithNonce = serializedHeader + strconv.FormatUint(goRoutineNonce, 10)
						encryptedHeaderWithNonce = sha256.Sum256([]byte(serializedHeaderWithNonce))

						goRoutineHashVector = generateHv(parameters, encryptedHeaderWithNonce)
						goRoutineHashVector, goRoutineOutputWord, _ = OptimizedDecoding(block, goRoutineHashVector, H, rowInCol, colInRow)
						flag := MakeDecision(block, colInRow, goRoutineOutputWord)

						select {
						case <-goRoutineSignal:
							// fmt.Println("goRoutineSignal channel is already closed")
							break attemptLoop
						default:
							if flag {
								close(goRoutineSignal)
								close(innerLoopSignal)
								fmt.Printf("Codeword is founded with nonce = %d\n", goRoutineNonce)
								LDPCNonce = goRoutineNonce
								hashVector = goRoutineHashVector
								outputWord = goRoutineOutputWord
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

	return hashVector, outputWord, LDPCNonce
}

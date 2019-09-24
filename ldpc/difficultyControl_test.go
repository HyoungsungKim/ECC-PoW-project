package ldpc

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestTablePrint(t *testing.T) {
	for i := range Table {
		fmt.Printf("level : %v, n : %v, wc : %v, wr : %v, decisionFrom : %v, decisionTo : %v, decisionStep : %v, miningProb : %v \n", Table[i].level, Table[i].n, Table[i].wc, Table[i].wr, Table[i].decisionFrom, Table[i].decisionTo, Table[i].decisionStep, Table[i].miningProb)
	}
}

func TestPrintReciprocal(t *testing.T) {
	for i := range Table {
		val := 1 / Table[i].miningProb
		bigInt := FloatToBigInt(val)
		fmt.Printf("Reciprocal of miningProb : %v \t big Int : %v\n", val, bigInt)
	}
}

func TestConversionFunc(t *testing.T) {
	for i := range Table {
		difficulty := ProbToDifficulty(Table[i].miningProb)
		miningProb := DifficultyToProb(difficulty)

		// Consider only integer part.
		fmt.Printf("Difficulty : %v \t MiningProb : %v\t, probability compare : %v \n", difficulty, miningProb, math.Abs(miningProb-Table[i].miningProb) < 1)
	}
}

func TestDifficultyChange(t *testing.T) {
	currentLevel := 0
	currentBlock := ethHeader{}
	// Parent block's timestamp is 0
	// compare elapse time(timestamp) and parent block's timestamp(0)

	for i := 0; i < 7; i++ {
		currentBlock.Difficulty = ProbToDifficulty(Table[currentLevel].miningProb)
		fmt.Printf("Current Difficulty : %v\n", currentBlock.Difficulty)

		startTime := time.Now()

		parameters := Parameters{
			n:  Table[currentLevel].n,
			wc: Table[currentLevel].wc,
			wr: Table[currentLevel].wr,
		}
		parameters.m = int(parameters.n * parameters.wc / parameters.wr)
		parameters.seed = GenerateSeed(currentBlock.ParentHash)

		RunOptimizedConcurrencyLDPC(parameters, currentBlock)
		timeStamp := uint64(time.Since(startTime).Seconds())
		fmt.Printf("Block generation time : %v\n", timeStamp)

		difficultyCalculator := MakeLDPCDifficultyCalculator()
		nextDifficulty := difficultyCalculator(timeStamp, &currentBlock)
		nextLevel := SearchLevel(nextDifficulty)

		fmt.Printf("Current prob : %v, Next Level : %v,  Next difficulty : %v, Next difficulty from table : %v\n\n", Table[currentLevel].miningProb, Table[nextLevel].level, nextDifficulty, ProbToDifficulty(Table[nextLevel].miningProb))

		// currentBlock.ParentHash = outputWord conversion from []int to [32]byte
		currentLevel = nextLevel
	}
}

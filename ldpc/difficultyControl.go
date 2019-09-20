package ldpc

import (
	"fmt"
	"math/big"
)

/*
	https://ethereum.stackexchange.com/questions/5913/how-does-the-ethereum-homestead-difficulty-adjustment-algorithm-work?noredirect=1&lq=1
	https://github.com/ethereum/EIPs/issues/100

	Ethereum difficulty adjustment
	 algorithm:
	diff = (parent_diff +
	         (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // 9), -99))
	        ) + 2^(periodCount - 2)

	LDPC difficulty adjustment
	algorithm:
	diff = (parent_diff +
			(parent_diff / 256 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // BlockGenerationTime), -99)))

	Why 8?
	This number is sensitivity of blockgeneration time
	If this number is high, difficulty is not changed much when block generation time is different from goal of block generation time
	But if this number is low, difficulty is changed much when block generatin time is different  from goal of block generation time

*/

// "github.com/ethereum/go-ethereum/consensus/ethash/consensus.go"
// Some weird constants to avoid constant memory allocs for them.
var (
	expDiffPeriod = big.NewInt(100000)
	big1          = big.NewInt(1)
	big2          = big.NewInt(2)
	big9          = big.NewInt(9)
	big10         = big.NewInt(10)
	bigMinus99    = big.NewInt(-99)

	MinimumDifficulty   = ProbToDifficulty(Table[0].miningProb)
	BlockGenerationTime = big.NewInt(100)
)

// MakeLDPCDifficultyCalculator calculate difficulty using difficulty table
func MakeLDPCDifficultyCalculator() func(time uint64, parent *ethHeader) *big.Int {
	return func(time uint64, parent *ethHeader) *big.Int {
		bigTime := new(big.Int).SetUint64(time)
		bigParentTime := new(big.Int).SetUint64(parent.Time)

		// holds intermediate values to make the algo easier to read & audit
		x := new(big.Int)
		y := new(big.Int)

		// (2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // BlockGenerationTime
		x.Sub(bigTime, bigParentTime)
		fmt.Printf("block_timestamp - parent_timestamp : %v\n", x)

		x.Div(x, BlockGenerationTime)
		fmt.Printf("(block_timestamp - parent_timestamp) / BlockGenerationTime : %v\n", x)
		/*
			Don't consider uncle block in difficulty test
				if parent.UncleHash == EmptyUncleHash {
					x.Sub(big1, x)
				} else {
					x.Sub(big2, x)
				}
		*/
		x.Sub(big1, x)
		fmt.Printf("1 - (block_timestamp - parent_timestamp) / BlockGenerationTime : %v\n", x)

		// max((2 if len(parent_uncles) else 1) - (block_timestamp - parent_timestamp) // 9, -99)
		if x.Cmp(bigMinus99) < 0 {
			x.Set(bigMinus99)
		}
		fmt.Printf("max(1 - (block_timestamp - parent_timestamp) / BlockGenerationTime, -99) : %v\n", x)

		// parent_diff + (parent_diff / 8 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) // BlockGenerationTime), -99))
		y.Div(parent.Difficulty, big.NewInt(8))
		fmt.Printf("parent.Difficulty / 8 : %v\n", y)

		x.Mul(y, x)
		fmt.Printf("parent.Difficulty / 8 * max(1 - (block_timestamp - parent_timestamp) / BlockGenerationTime, -99) : %v\n", x)

		x.Add(parent.Difficulty, x)
		fmt.Printf("parent.Difficulty - parent.Difficulty / 8 * max(1 - (block_timestamp - parent_timestamp) / BlockGenerationTime, -99) : %v\n", x)

		// minimum difficulty can ever be (before exponential factor)
		if x.Cmp(MinimumDifficulty) < 0 {
			x.Set(MinimumDifficulty)
		}
		fmt.Printf("x : %v, Minimum difficulty : %v\n", x, MinimumDifficulty)

		return x
	}
}

// SearchNextLevel return next level by using currentDifficulty of header
// Type of Ethereum difficulty is *bit.Int so arg is *big.int
func SearchNextLevel(nextDifficulty *big.Int) int {
	// foo := MakeLDPCDifficultyCalculator()
	// Next level := SearchNextLevel(foo(currentBlock's time stamp, parentBlock))

	var currentProb = DifficultyToProb(nextDifficulty)
	var nextLevel int

	for i := range Table {
		if currentProb > Table[i].miningProb {
			nextLevel = Table[i].level
			break
		}
	}

	return nextLevel
}

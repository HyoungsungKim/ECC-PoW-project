# Error Correction Code Prove of Work(ECCPoW)

- [***LDPC-pseudo-code***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/Ecc-PoW-pseudo-code) : Pseudo code of ECCPoW using C++ which are cloned from https://github.com/paaabx3/ECCPoW
- [***ECC-PoW-pseudo-blockchain***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/LDPC-pseudo-code) : Implement pseudo ECCPoW blockchain using python.
  - Blockchain source code is based on https://github.com/dvf/blockchain
- [***ldpc***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/ldpc) : Porting LDPC C++ to LDPC golang  

Writer : HyoungSung Kim 

Github : https://github.com/hyoungsungkim

Email : rktkek456@gmail.com / hyoungsung@gist.ac.kr

## LDPC decoder porting to Go version report

### 2019.08.22 Finish Porting to Go

- Finish porting to go

### 2019.08.28 Remove global variables

- Index errors are happened when LDPC is tested in go-ethereum using go routine
  - Remove global variables
  - Because of global variables, Critical section is violated
- Remove useless return
- Add comments to each function

### 2019.08.29 Add test file

- Add `decoder_test.go`
  - Calculate elapse time of decoding
  - Test LDPC Process
  - Test `RunLDPC()` function
  - Test LDPC verification
- Implement a function for verifying LDPC decoder
  - Correct return of few functions to pass information to the decoder verification function
- `GenerateHV()` function is corrected
  - Before correcting, serialized string was passed
  - But now, encrypted(sha256) string is passed 

### 2019.09.01 Decoding function optimized

- `OptimizedDecoding` function is implemented
  - When every condition is same, It is 20% faster than previous decoder (Different up to seed)

Previous Implementation

```go
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
```

Optimized Implementation

```go
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
```

### 2019.09.05 LDPCNonce generating way is changed

- Now LDPCNonce is started from Random number, not 0
  - It is same way with go-ethereum
  - `crand` is `crypto/rand`
- Previous LDPCNonce is uint32, however now LDPCNonce is uint64

go-ethereum

```go
// go-ethereum/consensus/ethash/sealer.go
func (ethash *Ethash) Seal(chain consensus.ChainReader, block *types.Block, result chan<- *types.Block, stop <- chan struct{}) error {
    .
    .
    if etheash.lock.Lock() {
        seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
        .
        .
        ethash.rand = rand.New(rand.NewSource(seed.Int64()))
    }
    .
    .
    for i := 0; i < threads; i++ {
        // nonce is started from random number
        go func(id int, nonce int64) {
            ...
        }(i, uint64(ethash.rand.Int63()))
    }
}
```

LDPC Decoder

```go
// decoder.go
func generateRandomNonce() uint64 {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	source := rand.New(rand.NewSource(seed.Int64()))

	return uint64(source.Int63())
}
// LDPCNonce := generateRandomNonce()
```

### 2019.09.18  Concurrency mining is implemented, Now LDPC Nonce is not incremented

- ***Now concurrency mining is possible!***
- Now LDPCNonce is random generation number

#### Basic Architecture

- Results which are derived by writing down in `default` might be same with this architecture
- But I think It is easier to read than writing down in `default`

```go
func RunOptimizedConcurrencyLDPC(...) {
    // ...
    var outerLoopSignal = make(chan struct{})
    var innerLoopSignal = make(chan struct{})
    var goRoutinerSignal = make(chan struct{})

    outerLoop:
    for {
        select{
        case <-OuterLoopSignal:
            break outerLoop
            
        default:
            // Empty default to unblock select statement
        }

        innerLoop:
        for i:= 0; i < numberOfGorouine; i++ {
            select {
            case <-innerLoopSignal:
                close(outerLoop)
                break innerLoop
                
            default:
                // Empty default to unblock select statement
            }

            go func(goRoutineSingal chan struct{})  {
                // ... 
                select {
                case <-goRoutineSignal:
                    break
                    
                default:
                    attemptLoop:
                    for attempt := 0; attempt < numberOfAttempt; attempt++ {
                        goRoutineNonce := generateRandomNonce()
                        // ...                    
                        select {
                        case <-goRotineSignal:
                            break attemptLoop
                        default:
                            // decoding
                        }
                    }
                }

            }(goRoutineSignale)
        }
        // Need to wait to prevent memory leak casued by unfinished goroutine
        wg.Wait()
    }
}
```

#### Result

- Before concurrency mining
  - In the lowest difficulty, it takes more than 600s
  - Usually more than 200s
- After concurrency mining
  - Tested 21 times, Only 1 test over 600s
  - Minimum is 9s,
  - Results(sec) : 9, 10, 17, 24, 30, 33, 40, 42, 45, 60, 60, 116, 143, 160 169, 210, 214, 214, 218, 220,  more than 600

#### Issues

- What is the number of optimal goroutines?
  - Why It is important?
    - Because, It there are too many goroutine, Overhead is happened.
    - It takes a time in `wg.Wait()`
    - Too many goroutine is slower because of scheduling
  - How about using constant?
    - We can get a better result if we use a constant which is derived by test
    - But it is dependent on system
    - It can be worse in different system
- What is the number of optimal attempts?
  - Why it is important?
    - Too many attempts make overhead in `wg.Wait()`
    - Too low attempts let goroutine meaningless
      - If attempts finish too early, single goroutine can be faster than multi goroutine(overhead)
  - Attempts can be higher
    - 1 attempt takes less than 1ms(Different up to difficulty)

#### Now LDPCNonce is not incremented

- Duplicated decoding can be happen when we increase LDPCNonce
  - When the number of attempts is high and distance of random generation number is close, duplication can be happen
  - For example
    - Attempt is 10,000
    - First goroutine's random generation number : 1
    - Second goroutine's random generation number : 5001
    - Then 5001~10000 are duplicated
  - It is very rare because range of random generation number is 0 ~ 2^64 -1
- Every attempt use different LDPCNonce which is generated randomly
  - Heuristically it is faster(Need more test)

### 2019.09.20 Difficulty change is implemented


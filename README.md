# Error Correction Code Proof of Work(ECCPoW)

- [***LDPC-pseudo-code***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/Ecc-PoW-pseudo-code) : Pseudo code of ECCPoW using C++ which are cloned from https://github.com/paaabx3/ECCPoW
- [***ECC-PoW-pseudo-blockchain***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/LDPC-pseudo-code) : Implement pseudo ECCPoW blockchain using python.
  - Blockchain source code is based on https://github.com/dvf/blockchain
- [***ldpc***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/ldpc) : Porting LDPC C++ to LDPC golang  
- [***ECC-PoW go-ethereum***](https://github.com/cryptoecc/ETH-ECC/tree/master/consensus/eccpow) (Organization repository)

Writer : HyoungSung Kim 

Github : https://github.com/hyoungsungkim

Email : rktkek456@gmail.com / hyoungsung@gist.ac.kr

## Concept image

### Error correction code PoW

![Error correction code PoW](img/eccpow.png)

### Error correction code PoW with ethash

![Error correction code PoW with ethash](img/eccpow-with-ethash.png)

## LDPC decoder porting to Go version report

### 2019.08.22 Finish Porting to Go

- Finish porting to go

### 2019.08.28 Remove global variables

- Index errors are happened when LDPC is tested with go routine in go-ethereum
  - Remove global variables
  - Because of global variables, Critical section is violated
- Remove useless return
- Add comments to each function

### 2019.08.29 Add test files

- Add `decoder_test.go`
  - Calculate elapse time of decoding
  - Test LDPC Process
  - Test `RunLDPC()` function
  - Test LDPC verification
- Implement a function for verifying LDPC decoder
  - Amend return of few functions for decoder verification function
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

- `if mp != m` is same with add every thing and subtract specific value 

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
- Now LDPCNonce is random generation number in concurrency mining

#### Basic Architecture

- ***This go routines algorithms are fixed for mining in geth. It works only in puppeth***
- In puppeth, user can use only 1 thread. It is inefficient in testing.(Currently ECCPoW block generation time is not fixed)
- However in geth, user can change the number of threads for mining.

```go
// consensus/eccpow/algorithm.go
func RunOptimizedConcurrencyLDPC(...) {
    // ...
    
    // make empty structure to close go routine
    var outerLoopSignal = make(chan struct{})
    var innerLoopSignal = make(chan struct{})
    var goRoutinerSignal = make(chan struct{})

    outerLoop:
    // Repeate generating go routine until valid nounce is found
    for {
        select{
	// If outerLoopSignal channel is closed, then break outerLoop and stop mining
        case <-OuterLoopSignal:
            break outerLoop
            
        default:
            // Empty default to unblock select statement
        }

        innerLoop:
	// Generate go routines for concurrency mining
        for i:= 0; i < numberOfGorouine; i++ {
            select {
	    // innerLoopSignal is propagated from go routine.
	    // If innerLoopSignal is closed, then propagate close signal by closing outerLoopSignal and break innerLoop
            case <-innerLoopSignal:
                close(outerLoopSignal)
                break innerLoop
                
            default:
                // Empty default to unblock select statement
            }

	    wg.Add(1)	
            go func(goRoutineSingal chan struct{})  {
	    	defer wg.Done()
                // ... 
                select {		
                case <-goRoutineSignal:
                    break
                    
                default:
                    attemptLoop:
                    for attempt := 0; attempt < numberOfAttempt; attempt++ {
                        goRoutineNonce := generateRandomNonce()
                        // ...                                         
			
			flag = MakeDecision(header, colInRow, goRoutineOutputWord)
			select {
			// If one of go routine found nounce, then that go routine close chanel.
			// As a result, other go routines recognize that do not need to keep mining. So stop it.
			case <-goRoutineSignal:
				// fmt.Println("goRoutineSignal channel is already closed")
				break attemptLoop
		    	default:
				// If valid nounce is found, then close goRoutineSignal to let other go routine stop mining
				// also close innerLoopSignal to propagate close signal to outer loop
				if flag {
				    close(goRoutineSignal)
				    close(innerLoopSignal)
				    hashVector = goRoutineHashVector
				    outputWord = goRoutineOutputWord
				    LDPCNonce = goRoutineNonce
				    digest = seed
				    break attemptLoop
				}		    
			}
		    }
                }
            }(goRoutineSignale)
        }
        // Need to wait to prevent memory leak
        wg.Wait()
    }
}
```

#### Result

- Before concurrency mining
  - In the lowest difficulty, it takes more than 600s
  - Usually more than 200s
- After concurrency mining
  - Tested 21 times, Only 1 test took more than 600s
  - Minimum is 9s,
  - Results(sec) : 9, 10, 17, 24, 30, 33, 40, 42, 45, 60, 60, 116, 143, 160 169, 210, 214, 214, 218, 220,  more than 600

#### Issues

- What is the number of optimal goroutines?
  - Why is it important?
    - Because, if there are too many goroutine, Overhead is happened.
    - It takes a time in `wg.Wait()`
    - Too many goroutine is slower because of scheduling
  - How about using constant?
    - We can get a better result if we use a constant which is derived by test
    - But it depends on system
    - It can be worse in different system
- What is the number of optimal attempts?
  - Why is it important?
    - Too many attempts make overhead in `wg.Wait()`
    - Too low attempts let goroutine meaningless
      - If attempts finish too early, single goroutine can be faster than multi goroutine(overhead)

#### Now LDPCNonce is not incremented(in concurrency mining)

- Duplicated decoding can be happen when we increase LDPCNonce
  - When the number of attempts is high and distance of random generation number is close, duplication can be happen
  - For example
    - Attempt is 10,000
    - First goroutine's random generation number : 1
    - Second goroutine's random generation number : 5001
    - Then 5001~10000 are duplicated
  - It is very rare because range of random generation number is 0 ~ 2^64 -1
- Every attempt use different LDPCNonce which is generated randomly
  - Empirically it is faster(Need more tests)

### 2019.09.20 Difficulty change is implemented

- Difficulty is reciprocal of mining success probability in ethereum
  - `Target <= 2^256 / Difficulty`
  - It means `Difficulty <= 2^256 / Target`
  - Difficulty : 블록 생성에 필요한 시도 횟수의 상한(Upper bound of the number of try to generate block)
  - Difficulty / Block generation time : Hashrate
- Therefore, we can convert probability of Table to Difficulty of Header

#### Difficulty of Ethereum

https://github.com/ethereum/EIPs/issues/100

```
algorithm:
diff = (parent_diff + (parent_diff / 2048 * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) / 9), -99))) + 2^(periodCount - 2)
```

Analysis

- Basic concept is adding or subtracting to difficulty of parent
- (timestamp - parent.timestamp) / 9
  - When block generation takes [0,9) sec, difficulty is increased
  - When block generation takes [9, 18) sec, difficulty is not changed
  - When block generation takes over 18 sec, then difficulty is decreased
- 2^(periodCount - 2) : For ice age
- 2048 : I will define it as sensitivity
  - Because when this number becomes higher, it is robust to difficulty change(difficulty is changed little by little)
  - When this number becomes smaller, it is weak to difficulty change(difficulty is changed rapidly )
  - On-line방식으로 수정 됨

#### Difficulty of LDPC decoder

```
algorithm:
diff = (parent_diff + (parent_diff / sensitivity * max((2 if len(parent.uncles) else 1) - ((timestamp - parent.timestamp) / Block_generation_time), -99)))
```

- Everything is same but we don't use ice age and need parameter tuning(sensitivity, Block_generation_time)

[source code](https://github.com/HyoungsungKim/ECC-PoW-project/blob/master/ldpc/difficultyControl.go)

### 2019.09.23 Porting to go-ethereum

https://github.com/HyoungsungKim/go-ethereum/tree/fix-ldpc-eccpow-1.9/consensus/eccpow

- LDPC decoder is ported in go-ethereum
- Need test for parameter tuning

### 2019.09.24 Implement decoding performance test function

- Test 1,000,000 try
  - Single goroutine : 126.385s
  - Multi goroutine : 24.595s

### 2019.09.27 Implement mixhash verification  

- However, There is a problem in verifying codeword
- Need to check arguments of verify function  

### 2019.09.28 Fix verifySeal and implement unit test

- Thread(go routine) local variable was problem (H matrix)
- Fix it as thread share variable.
- It can be shared with other threads. Because it is not changed in thread(goroutine)
- Implement unit test for verification.
- Currently, approximate block generation time is 100 ~ 120 sec

### 2019.10.04 Start mining monitoring
[Reference](https://steemit.com/kr-dev/@modolee/puppeth-aws-private-ethereum-network-poa)

### 2019.10.05 Testing

- Block generation time is too irregular
  - Need more nodes...?
  - Image is temparaily deleted for journal paper publication

![Testing](img/Testing.png)

### 2019.10.06 Testing 2

- Add more nodes
- Image is temparaily deleted for journal paper publication

![Testing2](img/Testing2.png)

### 2019.10.14

- Fix difficulty calculation detail

### 2019.10.21

- Set block generation time as 60

### 2019.10.24

- Change architecture for certification test
- ECCPoW -> ethash + ldpc decoder

![Testing3](img/eccpow-with-ethash-test.png)

### 2020.7

- Fix concurrency error for geth deployment.
- Last concurrency mining works only in puppeth deployment.
	- It makes network errors. so fixed.

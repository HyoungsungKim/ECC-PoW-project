# Error Correction Code Prove of Work(ECCPoW)

[***LDPC-pseudo-code***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/Ecc-PoW-pseudo-code) : Pseudo code of ECCPoW using C++ which are cloned from https://github.com/paaabx3/ECCPoW
[***Ecc-PoW-pseudo-blockchain***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/LDPC-pseudo-code) : Implement pseudo ECCPoW blockchain using python. Blockchain source code is based on https://github.com/dvf/blockchain
[***ldpc***](https://github.com/HyoungsungKim/ECC-PoW-project/tree/master/ldpc) : Porting LDPC C++ to LDPC golang  

Writer : HyoungSung Kim 

Github : https://github.com/hyoungsungkim

Email : rktkek456@gmail.com / hyoungsung@gist.ac.kr

# LDPC decoder porting to go version report

2019.08.22

- Finish porting to go

2019.08.28 

- Index errors are happened when LDPC is tested in go-ethereum using go routine
  - Remove global variables
  - Because of global variables, Critical section is violated
- Remove useless return
- Add comments to each function

2019.08.29

- Add `decoder_test.go`
  - Calculate elapse time of decoding
  - Test LDPC Process
  - Test `RunLDPC()` function
  - Test LDPC verification
- Implement the function for verifying LDPC decoder
  - Correct return of few function to pass information to decoder verification function
- `GenerateHV()` function is corrected
  - Before correcting, serialized string was passed
  - But now, encrypted(sha256) string is passed 

2019.09.01

- `OptimizedDecoding` function is implemented
  - It is 20% faster than previous decoder (Different up to seed)

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


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


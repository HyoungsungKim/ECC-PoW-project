# Error Correction Code Prove of Work(ECCPoW)

Writer HyoungSung Kim 

Github : https://github.com/hyoungsungkim

Email : hyoungsung@gist.ac.kr / rktkek456@gmil.com

## Decoder

### Input of Decoder

- Input : $F_b$ , difficulty, error

#### What is $F_b$?

Definition of $F_b$ : 

***Difficulty of $F_{b1}$ will be fixed for test***

#### Error

##### How to generate Error?

Denote error as `e` which is `error pattern` + `sameLength1DVector`

Error patterns  are denoted as $ep_n$
$$
e = ep_n +  sameLength1DVector
$$
if error can be corrected
$$
F_b*e = zero\ matrix
$$
and $*$ ***means matrix multiplication*** not dot production or element-wise production. 

##### Error pattern

Example : 1 bit error
$$
ep_1 = \underbrace{[1,0,0, ... ,0]}_{256},
$$

$$
ep_2 = \underbrace{[0,2,0, ... ,0]}_{256},
$$

$$
.
$$

$$
.
$$

$$
ep_{256} = \underbrace{[0,0,0, ... ,1]}_{256}
$$

All of error patterns have to pass decoder without any problems

###### Number of n bit error pattern in $2^{256}$

$$
\binom{256}{n}
$$

### How decoder works

#### How many errors can be corrected

If `t error pattern` can be covered.
$$
ep_1 + ep_2 + ep_3 + ... + ep_t = \binom{256}{1} + \binom {256}{2} + \binom {256}{3} + ... + \binom {256}{t}
$$
And ***It is defined as `# of word error corrected`***

####  Codeword

denote codeword as c, message as m
$$
c = G * m
$$

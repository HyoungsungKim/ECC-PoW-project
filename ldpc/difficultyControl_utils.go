package ldpc

import (
	"math/big"
)

// Table : level, n, wc, wr, decisionFrom, decisionTo, decisionStep, _, miningProb
/*
	How to decision
	1. The number of ones in outputword % decision_step == 0
	2. The number of ones in outputword exists between decision_from to decision_to

	How to change difficulty level
	- Reciprocal of difficulty is probability
	- Therefore We can define difficulty is reciprocal of probability
	- Find close probability

*/
type difficulty struct {
	level        int
	n            int
	wc           int
	wr           int
	decisionFrom int
	decisionTo   int
	decisionStep int
	_            float32
	miningProb   float64
}

// FloatToBigInt convert float64 to big integer
func FloatToBigInt(val float64) *big.Int {
	// float64 -> bit float -> big int
	bigFloat := big.NewFloat(val)
	bigInt := new(big.Int)
	bigFloat.Int(bigInt)

	return bigInt
}

// BigIntToFloat convert big int to float64
func BigIntToFloat(val *big.Int) float64 {
	// big int -> bit float -> float64
	bigFloat := new(big.Float).SetInt(val)
	floatVal, _ := bigFloat.Float64()

	return floatVal
}

// DifficultyToProb convert difficulty to probability of table
func DifficultyToProb(difficulty *big.Int) float64 {
	//big Int -> 1/bigInt -> float64
	prob := 1 / BigIntToFloat(difficulty)
	return prob
}

// ProbToDifficulty convert probability to difficulty of header
func ProbToDifficulty(miningProb float64) *big.Int {
	// float64 -> 1/float64 -> big Int
	difficulty := FloatToBigInt(1 / miningProb)
	return difficulty
}

// Table is difficulty table slice
var Table = []difficulty{
	{0, 32, 3, 4, 10, 22, 2, 0.329111, 3.077970e-05},
	{1, 32, 3, 4, 10, 22, 2, 0.329111, 3.077970e-05},
	{2, 32, 3, 4, 10, 16, 2, 0.329111, 2.023220e-05},
	{3, 32, 3, 4, 16, 16, 1, 0.329111, 9.684650e-06},
	{4, 32, 3, 4, 14, 14, 1, 0.329111, 6.784080e-06},
	{5, 36, 3, 4, 12, 24, 2, 0.329111, 4.830240e-06},
	{6, 36, 3, 4, 12, 18, 2, 0.369449, 3.125970e-06},
	{7, 32, 3, 4, 12, 12, 1, 0.369449, 2.862890e-06},
	{8, 44, 3, 4, 14, 30, 2, 0.369449, 1.637790e-06},
	{9, 36, 3, 4, 18, 18, 1, 0.369449, 1.421700e-06},
	{10, 36, 3, 4, 16, 16, 1, 0.369449, 1.051350e-06},
	{11, 44, 3, 4, 14, 22, 2, 0.411046, 1.029740e-06},
	{12, 40, 3, 4, 12, 28, 2, 0.411046, 7.570880e-07},
	{13, 36, 3, 4, 14, 14, 1, 0.411046, 4.865630e-07},
	{14, 40, 3, 4, 12, 20, 2, 0.411046, 4.813320e-07},
	{15, 44, 3, 4, 22, 22, 1, 0.411046, 4.216920e-07},
	{16, 44, 3, 4, 20, 20, 1, 0.411046, 3.350070e-07},
	{17, 48, 3, 4, 14, 34, 2, 0.452453, 2.677070e-07},
	{18, 40, 3, 4, 20, 20, 1, 0.452453, 2.055750e-07},
	{19, 44, 3, 4, 18, 18, 1, 0.452453, 1.788400e-07},
	{20, 48, 3, 4, 14, 24, 2, 0.452453, 1.664080e-07},
	{21, 40, 3, 4, 18, 18, 1, 0.452453, 1.583110e-07},
	{22, 40, 3, 4, 16, 16, 1, 0.452453, 7.917230e-08},
	{23, 44, 3, 4, 16, 16, 1, 0.498513, 7.103820e-08},
	{24, 48, 3, 4, 24, 24, 1, 0.498513, 6.510890e-08},
	{25, 48, 3, 4, 22, 22, 1, 0.498513, 5.300760e-08},
	{26, 52, 3, 4, 14, 40, 2, 0.498513, 4.266600e-08},
	{27, 48, 3, 4, 20, 20, 1, 0.498513, 2.990510e-08},
	{28, 40, 3, 4, 14, 14, 1, 0.498513, 2.927380e-08},
	{29, 52, 3, 4, 14, 26, 2, 0.498513, 2.626790e-08},
	{30, 60, 3, 4, 18, 42, 2, 0.498513, 1.485240e-08},
	{31, 48, 3, 4, 18, 18, 1, 0.546238, 1.267290e-08},
	{32, 52, 3, 4, 26, 26, 1, 0.546238, 9.891110e-09},
	{33, 60, 3, 4, 18, 30, 2, 0.546238, 9.019200e-09},
	{34, 48, 3, 4, 16, 32, 1, 0.546238, 8.762650e-09},
	{35, 52, 3, 4, 24, 24, 1, 0.546238, 8.213140e-09},
	{36, 56, 3, 4, 16, 42, 2, 0.546238, 6.658250e-09},
	{37, 52, 3, 4, 22, 22, 1, 0.546238, 4.856960e-09},
	{38, 48, 3, 4, 16, 16, 1, 0.546238, 4.381330e-09},
	{39, 56, 3, 4, 16, 28, 2, 0.546238, 4.068000e-09},
	{40, 60, 3, 4, 30, 30, 1, 0.546238, 3.186040e-09},
	{41, 60, 3, 4, 28, 28, 1, 0.578290, 2.725470e-09},
	{42, 64, 3, 4, 18, 46, 2, 0.578290, 2.410890e-09},
	{43, 52, 3, 4, 20, 20, 1, 0.578290, 2.181360e-09},
	{44, 60, 3, 4, 26, 26, 1, 0.578290, 1.737940e-09},
	{45, 52, 3, 4, 18, 34, 1, 0.578290, 1.595330e-09},
	{46, 56, 3, 4, 28, 28, 1, 0.578290, 1.481830e-09},
	{47, 64, 3, 4, 18, 32, 2, 0.578290, 1.454780e-09},
	{48, 56, 3, 4, 26, 26, 1, 0.578290, 1.250550e-09},
	{49, 60, 3, 4, 24, 24, 1, 0.578290, 8.614860e-10},
	{50, 52, 3, 4, 18, 18, 1, 0.578290, 7.976650e-10},
	{51, 56, 3, 4, 24, 24, 1, 0.628015, 7.700380e-10},
	{52, 60, 3, 4, 22, 38, 1, 0.628015, 6.978800e-10},
	{53, 52, 3, 4, 16, 36, 1, 0.628015, 5.069080e-10},
	{54, 64, 3, 4, 32, 32, 1, 0.628015, 4.986660e-10},
	{55, 64, 3, 4, 30, 30, 1, 0.628015, 4.315180e-10},
	{56, 68, 3, 4, 18, 50, 2, 0.628015, 3.848530e-10},
	{57, 56, 3, 4, 22, 22, 1, 0.628015, 3.643130e-10},
	{58, 60, 3, 4, 22, 22, 1, 0.628015, 3.489400e-10},
	{59, 64, 3, 4, 28, 28, 1, 0.628015, 2.836780e-10},
	{60, 56, 3, 4, 20, 36, 1, 0.628015, 2.809120e-10},
	{61, 52, 3, 4, 16, 16, 1, 0.666500, 2.534540e-10},
	{62, 60, 3, 4, 20, 40, 1, 0.666500, 2.427110e-10},
	{63, 68, 3, 4, 18, 34, 2, 0.666500, 2.309280e-10},
	{64, 64, 3, 4, 26, 26, 1, 0.666500, 1.466250e-10},
	{65, 56, 3, 4, 20, 20, 1, 0.666500, 1.404560e-10},
	{66, 76, 3, 4, 22, 54, 2, 0.666500, 1.375500e-10},
	{67, 60, 3, 4, 20, 20, 1, 0.666500, 1.213550e-10},
	{68, 56, 3, 4, 18, 38, 1, 0.666500, 9.340240e-11},
	{69, 76, 3, 4, 22, 38, 2, 0.666500, 8.174200e-11},
	{70, 68, 3, 4, 34, 34, 1, 0.666500, 7.700290e-11},
	{71, 68, 3, 4, 32, 32, 1, 0.666500, 6.729690e-11},
	{72, 64, 3, 4, 24, 24, 1, 0.706860, 6.217280e-11},
	{73, 72, 3, 4, 18, 56, 2, 0.706860, 6.056200e-11},
	{74, 56, 3, 4, 18, 18, 1, 0.706860, 4.670120e-11},
	{75, 68, 3, 4, 30, 30, 1, 0.706860, 4.543980e-11},
	{76, 64, 3, 4, 22, 42, 1, 0.706860, 4.517330e-11},
	{77, 72, 3, 4, 18, 36, 2, 0.706860, 3.615450e-11},
	{78, 76, 3, 4, 38, 38, 1, 0.706860, 2.593400e-11},
	{79, 68, 3, 4, 28, 28, 1, 0.706860, 2.438720e-11},
	{80, 76, 3, 4, 36, 36, 1, 0.706860, 2.303460e-11},
	{81, 64, 3, 4, 22, 22, 1, 0.706860, 2.258660e-11},
	{82, 80, 3, 4, 22, 58, 2, 0.706860, 2.229400e-11},
	{83, 76, 3, 4, 34, 34, 1, 0.706860, 1.626350e-11},
	{84, 64, 3, 4, 20, 40, 1, 0.706860, 1.465310e-11},
	{85, 80, 3, 4, 22, 40, 2, 0.763542, 1.319160e-11},
	{86, 72, 3, 4, 36, 36, 1, 0.763542, 1.174900e-11},
	{87, 68, 3, 4, 26, 26, 1, 0.763542, 1.078820e-11},
	{88, 72, 3, 4, 34, 34, 1, 0.763542, 1.035690e-11},
	{89, 76, 3, 4, 32, 32, 1, 0.763542, 9.311370e-12},
	{90, 68, 3, 4, 24, 44, 1, 0.763542, 8.173020e-12},
	{91, 64, 3, 4, 20, 20, 1, 0.763542, 7.326570e-12},
	{92, 72, 3, 4, 32, 32, 1, 0.763542, 7.160350e-12},
	{93, 76, 3, 4, 30, 30, 1, 0.763542, 4.440960e-12},
	{94, 80, 3, 4, 40, 40, 1, 0.763542, 4.089220e-12},
	{95, 68, 3, 4, 24, 24, 1, 0.763542, 4.086510e-12},
	{96, 72, 3, 4, 30, 30, 1, 0.763542, 3.975570e-12},
	{97, 80, 3, 4, 38, 38, 1, 0.763542, 3.656430e-12},
	{98, 76, 3, 4, 28, 48, 1, 0.763542, 3.634830e-12},
	{99, 84, 3, 4, 22, 62, 2, 0.763542, 3.566880e-12},
	{100, 68, 3, 4, 22, 46, 1, 0.783762, 2.750600e-12},
	{101, 80, 3, 4, 36, 36, 1, 0.783762, 2.630600e-12},
	{102, 84, 3, 4, 22, 42, 2, 0.783762, 2.102180e-12},
	{103, 72, 3, 4, 28, 28, 1, 0.783762, 1.828850e-12},
	{104, 76, 3, 4, 28, 28, 1, 0.783762, 1.817420e-12},
	{105, 80, 3, 4, 34, 34, 1, 0.783762, 1.548670e-12},
	{106, 72, 3, 4, 26, 46, 1, 0.783762, 1.441670e-12},
	{107, 68, 3, 4, 22, 22, 1, 0.783762, 1.375300e-12},
	{108, 76, 3, 4, 26, 50, 1, 0.783762, 1.314800e-12},
	{109, 92, 3, 4, 24, 68, 2, 0.783762, 1.296220e-12},
	{110, 68, 3, 4, 20, 48, 1, 0.783762, 8.516070e-13},
	{111, 80, 3, 4, 32, 32, 1, 0.783762, 7.636740e-13},
	{112, 92, 3, 4, 24, 46, 2, 0.783762, 7.585130e-13},
	{113, 72, 3, 4, 26, 26, 1, 0.824961, 7.208340e-13},
	{114, 76, 3, 4, 26, 26, 1, 0.824961, 6.573980e-13},
	{115, 80, 3, 4, 30, 50, 1, 0.824961, 6.475910e-13},
	{116, 84, 3, 4, 42, 42, 1, 0.824961, 6.374900e-13},
	{117, 84, 3, 4, 40, 40, 1, 0.824961, 5.734350e-13},
	{118, 88, 3, 4, 22, 66, 2, 0.824961, 5.640630e-13},
	{119, 72, 3, 4, 24, 48, 1, 0.824961, 5.032200e-13},
	{120, 76, 3, 4, 24, 52, 1, 0.824961, 4.325430e-13},
	{121, 68, 3, 4, 20, 20, 1, 0.824961, 4.258030e-13},
	{122, 84, 3, 4, 38, 38, 1, 0.824961, 4.195890e-13},
	{123, 88, 3, 4, 22, 44, 2, 0.824961, 3.312130e-13},
	{124, 80, 3, 4, 30, 30, 1, 0.824961, 3.237950e-13},
	{125, 84, 3, 4, 36, 36, 1, 0.824961, 2.533670e-13},
	{126, 72, 3, 4, 24, 24, 1, 0.824961, 2.516100e-13},
	{127, 80, 3, 4, 28, 52, 1, 0.824961, 2.424700e-13},
	{128, 92, 3, 4, 46, 46, 1, 0.865704, 2.208090e-13},
	{129, 76, 3, 4, 24, 24, 1, 0.865704, 2.162710e-13},
	{130, 96, 3, 4, 26, 72, 2, 0.865704, 2.099840e-13},
	{131, 92, 3, 4, 44, 44, 1, 0.865704, 2.006580e-13},
	{132, 72, 3, 4, 22, 50, 1, 0.865704, 1.605310e-13},
	{133, 92, 3, 4, 42, 42, 1, 0.865704, 1.511670e-13},
	{134, 84, 3, 4, 34, 34, 1, 0.865704, 1.288510e-13},
	{135, 96, 3, 4, 26, 48, 2, 0.865704, 1.224810e-13},
	{136, 80, 3, 4, 28, 28, 1, 0.865704, 1.212350e-13},
	{137, 88, 3, 4, 44, 44, 1, 0.865704, 9.836240e-14},
	{138, 92, 3, 4, 40, 40, 1, 0.865704, 9.542830e-14},
	{139, 88, 3, 4, 42, 42, 1, 0.865704, 8.895450e-14},
	{140, 80, 3, 4, 26, 54, 1, 0.865704, 8.227740e-14},
	{141, 72, 3, 4, 22, 22, 1, 0.865704, 8.026550e-14},
	{142, 88, 3, 4, 40, 40, 1, 0.865704, 6.609120e-14},
	{143, 84, 3, 4, 32, 32, 1, 0.865704, 5.648410e-14},
	{144, 92, 3, 4, 38, 38, 1, 0.908949, 5.127400e-14},
	{145, 72, 3, 4, 20, 52, 1, 0.908949, 4.822030e-14},
	{146, 84, 3, 4, 30, 54, 1, 0.908949, 4.372420e-14},
	{147, 80, 3, 4, 26, 26, 1, 0.908949, 4.113870e-14},
	{148, 88, 3, 4, 38, 38, 1, 0.908949, 4.084490e-14},
	{149, 96, 3, 4, 48, 48, 1, 0.908949, 3.497970e-14},
	{150, 100, 3, 4, 26, 76, 2, 0.908949, 3.365420e-14},
	{151, 96, 3, 4, 46, 46, 1, 0.908949, 3.192740e-14},
	{152, 80, 3, 4, 24, 56, 1, 0.908949, 2.593890e-14},
	{153, 96, 3, 4, 44, 44, 1, 0.908949, 2.435890e-14},
	{154, 72, 3, 4, 20, 20, 1, 0.908949, 2.411020e-14},
	{155, 92, 3, 4, 36, 36, 1, 0.908949, 2.388460e-14},
	{156, 84, 3, 4, 30, 30, 1, 0.908949, 2.186210e-14},
	{157, 88, 3, 4, 36, 36, 1, 0.908949, 2.137330e-14},
	{158, 92, 3, 4, 34, 58, 1, 0.908949, 1.967320e-14},
	{159, 100, 3, 4, 26, 50, 2, 0.908949, 1.957080e-14},
	{160, 96, 3, 4, 42, 42, 1, 0.908949, 1.568040e-14},
	{161, 84, 3, 4, 28, 56, 1, 0.908949, 1.529960e-14},
	{162, 80, 3, 4, 24, 24, 1, 0.954202, 1.296950e-14},
	{163, 108, 3, 4, 28, 82, 2, 0.954202, 1.237200e-14},
	{164, 92, 3, 4, 34, 34, 1, 0.954202, 9.836600e-15},
	{165, 88, 3, 4, 34, 34, 1, 0.954202, 9.667440e-15},
	{166, 96, 3, 4, 40, 40, 1, 0.954202, 8.634600e-15},
	{167, 88, 3, 4, 32, 56, 1, 0.954202, 7.725050e-15},
	{168, 84, 3, 4, 28, 28, 1, 0.954202, 7.649800e-15},
	{169, 92, 3, 4, 32, 60, 1, 0.954202, 7.305750e-15},
	{170, 108, 3, 4, 28, 54, 2, 0.954202, 7.154920e-15},
	{171, 100, 3, 4, 50, 50, 1, 0.954202, 5.487490e-15},
	{172, 104, 3, 4, 26, 78, 2, 0.954202, 5.340690e-15},
	{173, 100, 3, 4, 48, 48, 1, 0.954202, 5.028760e-15},
	{174, 84, 3, 4, 26, 58, 1, 0.954202, 4.951260e-15},
	{175, 96, 3, 4, 38, 38, 1, 0.954202, 4.134930e-15},
	{176, 100, 3, 4, 46, 46, 1, 0.954202, 3.881360e-15},
	{177, 88, 3, 4, 32, 32, 1, 0.954202, 3.862530e-15},
	{178, 92, 3, 4, 32, 32, 1, 0.954202, 3.652870e-15},
	{179, 96, 3, 4, 36, 60, 1, 0.954202, 3.505650e-15},
	{180, 104, 3, 4, 26, 52, 2, 0.993877, 3.096920e-15},
	{181, 88, 3, 4, 30, 58, 1, 0.993877, 2.785770e-15},
	{182, 100, 3, 4, 44, 44, 1, 0.993877, 2.543880e-15},
	{183, 92, 3, 4, 30, 62, 1, 0.993877, 2.493880e-15},
	{184, 84, 3, 4, 26, 26, 1, 0.993877, 2.475630e-15},
	{185, 112, 3, 4, 28, 86, 2, 0.993877, 2.003890e-15},
	{186, 108, 3, 4, 54, 54, 1, 0.993877, 1.937830e-15},
	{187, 108, 3, 4, 52, 52, 1, 0.993877, 1.788370e-15},
	{188, 96, 3, 4, 36, 36, 1, 0.993877, 1.752820e-15},
	{189, 84, 3, 4, 24, 60, 1, 0.993877, 1.514630e-15},
	{190, 100, 3, 4, 42, 42, 1, 0.993877, 1.433180e-15},
	{191, 108, 3, 4, 50, 50, 1, 0.993877, 1.408830e-15},
	{192, 88, 3, 4, 30, 30, 1, 0.993877, 1.392880e-15},
	{193, 96, 3, 4, 34, 62, 1, 0.993877, 1.339400e-15},
	{194, 92, 3, 4, 30, 30, 1, 0.993877, 1.246940e-15},
	{195, 112, 3, 4, 28, 56, 2, 0.993877, 1.155930e-15},
	{196, 108, 3, 4, 48, 48, 1, 0.993877, 9.534230e-16},
	{197, 88, 3, 4, 28, 60, 1, 0.993877, 9.258750e-16},
	{198, 104, 3, 4, 52, 52, 1, 1.035782, 8.531480e-16},
	{199, 92, 3, 4, 28, 64, 1, 1.035782, 7.972240e-16},
	{200, 104, 3, 4, 50, 50, 1, 1.035782, 7.847010e-16},
	{201, 84, 3, 4, 24, 24, 1, 1.035782, 7.573130e-16},
	{202, 100, 3, 4, 40, 40, 1, 1.035782, 7.043820e-16},
	{203, 96, 3, 4, 34, 34, 1, 1.035782, 6.696990e-16},
	{204, 100, 3, 4, 38, 62, 1, 1.035782, 6.138180e-16},
	{205, 104, 3, 4, 48, 48, 1, 1.035782, 6.121340e-16},
	{206, 108, 3, 4, 46, 46, 1, 1.035782, 5.596910e-16},
	{207, 96, 3, 4, 32, 64, 1, 1.035782, 4.694950e-16},
	{208, 88, 3, 4, 28, 28, 1, 1.035782, 4.629380e-16},
	{209, 104, 3, 4, 46, 46, 1, 1.035782, 4.079260e-16},
	{210, 92, 3, 4, 28, 28, 1, 1.035782, 3.986120e-16},
	{211, 116, 3, 4, 30, 86, 2, 1.035782, 3.215820e-16},
	{212, 112, 3, 4, 56, 56, 1, 1.035782, 3.079780e-16},
	{213, 100, 3, 4, 38, 38, 1, 1.035782, 3.069090e-16},
	{214, 88, 3, 4, 26, 62, 1, 1.035782, 2.893730e-16},
	{215, 108, 3, 4, 44, 44, 1, 1.035782, 2.884350e-16},
	{216, 112, 3, 4, 54, 54, 1, 1.035782, 2.851100e-16},
	{217, 92, 3, 4, 26, 66, 1, 1.035782, 2.429760e-16},
	{218, 100, 3, 4, 36, 64, 1, 1.083752, 2.410500e-16},
	{219, 104, 3, 4, 44, 44, 1, 1.083752, 2.347600e-16},
	{220, 96, 3, 4, 32, 32, 1, 1.083752, 2.347480e-16},
	{221, 112, 3, 4, 52, 52, 1, 1.083752, 2.266460e-16},
	{222, 116, 3, 4, 30, 58, 2, 1.083752, 1.850560e-16},
	{223, 112, 3, 4, 50, 50, 1, 1.083752, 1.555940e-16},
	{224, 96, 3, 4, 30, 66, 1, 1.083752, 1.536060e-16},
	{225, 88, 3, 4, 26, 26, 1, 1.083752, 1.446860e-16},
	{226, 108, 3, 4, 42, 42, 1, 1.083752, 1.322410e-16},
	{227, 92, 3, 4, 26, 26, 1, 1.083752, 1.214880e-16},
	{228, 100, 3, 4, 36, 36, 1, 1.083752, 1.205250e-16},
	{229, 124, 3, 4, 32, 94, 2, 1.083752, 1.192380e-16},
	{230, 104, 3, 4, 42, 42, 1, 1.083752, 1.182330e-16},
	{231, 108, 3, 4, 40, 68, 1, 1.083752, 1.093900e-16},
	{232, 112, 3, 4, 48, 48, 1, 1.083752, 9.304980e-17},
	{233, 100, 3, 4, 34, 66, 1, 1.083752, 8.672750e-17},
	{234, 88, 3, 4, 24, 64, 1, 1.083752, 8.671340e-17},
	{235, 96, 3, 4, 30, 30, 1, 1.083752, 7.680300e-17},
	{236, 124, 3, 4, 32, 62, 2, 1.083752, 6.831000e-17},
	{237, 108, 3, 4, 40, 40, 1, 1.083752, 5.469510e-17},
	{238, 104, 3, 4, 40, 40, 1, 1.083752, 5.287900e-17},
	{239, 120, 3, 4, 30, 90, 2, 1.122176, 5.116510e-17},
	{240, 112, 3, 4, 46, 46, 1, 1.122176, 4.900210e-17},
	{241, 116, 3, 4, 58, 58, 1, 1.122176, 4.853020e-17},
	{242, 96, 3, 4, 28, 68, 1, 1.122176, 4.769400e-17},
	{243, 116, 3, 4, 56, 56, 1, 1.122176, 4.505610e-17},
	{244, 100, 3, 4, 34, 34, 1, 1.122176, 4.336370e-17},
	{245, 88, 3, 4, 24, 24, 1, 1.122176, 4.335670e-17},
	{246, 104, 3, 4, 38, 66, 1, 1.122176, 4.264440e-17},
	{247, 108, 3, 4, 38, 70, 1, 1.122176, 4.139190e-17},
	{248, 116, 3, 4, 54, 54, 1, 1.122176, 3.611940e-17},
	{249, 120, 3, 4, 30, 60, 2, 1.122176, 2.937590e-17},
	{250, 100, 3, 4, 32, 68, 1, 1.122176, 2.904870e-17},
	{251, 116, 3, 4, 52, 52, 1, 1.122176, 2.512890e-17},
	{252, 96, 3, 4, 28, 28, 1, 1.122176, 2.384700e-17},
	{253, 112, 3, 4, 44, 44, 1, 1.122176, 2.300220e-17},
	{254, 104, 3, 4, 38, 38, 1, 1.122176, 2.132220e-17},
	{255, 108, 3, 4, 38, 38, 1, 1.122176, 2.069600e-17},
	{256, 112, 3, 4, 42, 70, 1, 1.122176, 1.949710e-17},
	{257, 128, 3, 4, 32, 98, 2, 1.122176, 1.931040e-17},
	{258, 124, 3, 4, 62, 62, 1, 1.122176, 1.738220e-17},
	{259, 124, 3, 4, 60, 60, 1, 1.122176, 1.622110e-17},
	{260, 104, 3, 4, 36, 68, 1, 1.165058, 1.573970e-17},
	{261, 116, 3, 4, 50, 50, 1, 1.165058, 1.529120e-17},
	{262, 108, 3, 4, 36, 72, 1, 1.165058, 1.452810e-17},
	{263, 100, 3, 4, 32, 32, 1, 1.165058, 1.452430e-17},
	{264, 124, 3, 4, 58, 58, 1, 1.165058, 1.320160e-17},
	{265, 128, 3, 4, 32, 64, 2, 1.165058, 1.103980e-17},
	{266, 112, 3, 4, 42, 42, 1, 1.165058, 9.748530e-18},
	{267, 124, 3, 4, 56, 56, 1, 1.165058, 9.408340e-18},
	{268, 100, 3, 4, 30, 70, 1, 1.165058, 9.198900e-18},
	{269, 116, 3, 4, 48, 48, 1, 1.165058, 8.218710e-18},
	{270, 104, 3, 4, 36, 36, 1, 1.165058, 7.869870e-18},
	{271, 120, 3, 4, 60, 60, 1, 1.165058, 7.586620e-18},
	{272, 112, 3, 4, 40, 72, 1, 1.165058, 7.557840e-18},
	{273, 108, 3, 4, 36, 36, 1, 1.165058, 7.264050e-18},
	{274, 120, 3, 4, 58, 58, 1, 1.165058, 7.062330e-18},
	{275, 124, 3, 4, 54, 54, 1, 1.165058, 5.908890e-18},
	{276, 120, 3, 4, 56, 56, 1, 1.165058, 5.705970e-18},
	{277, 104, 3, 4, 34, 70, 1, 1.165058, 5.397190e-18},
	{278, 108, 3, 4, 34, 74, 1, 1.165058, 4.794130e-18},
	{279, 100, 3, 4, 30, 30, 1, 1.165058, 4.599450e-18},
	{280, 120, 3, 4, 54, 54, 1, 1.165058, 4.019430e-18},
	{281, 116, 3, 4, 46, 46, 1, 1.165058, 3.945320e-18},
	{282, 112, 3, 4, 40, 40, 1, 1.165058, 3.778920e-18},
	{283, 116, 3, 4, 44, 72, 1, 1.231412, 3.423190e-18},
	{284, 124, 3, 4, 52, 52, 1, 1.231412, 3.297090e-18},
	{285, 100, 3, 4, 28, 72, 1, 1.231412, 2.795780e-18},
	{286, 128, 3, 4, 64, 64, 1, 1.231412, 2.769110e-18},
	{287, 112, 3, 4, 38, 74, 1, 1.231412, 2.714430e-18},
	{288, 104, 3, 4, 34, 34, 1, 1.231412, 2.698600e-18},
	{289, 128, 3, 4, 62, 62, 1, 1.231412, 2.590140e-18},
	{290, 120, 3, 4, 52, 52, 1, 1.231412, 2.486040e-18},
	{291, 108, 3, 4, 34, 34, 1, 1.231412, 2.397060e-18},
	{292, 128, 3, 4, 60, 60, 1, 1.231412, 2.122380e-18},
	{293, 104, 3, 4, 32, 72, 1, 1.231412, 1.744360e-18},
	{294, 116, 3, 4, 44, 44, 1, 1.231412, 1.711600e-18},
	{295, 124, 3, 4, 50, 50, 1, 1.231412, 1.649820e-18},
	{296, 128, 3, 4, 58, 58, 1, 1.231412, 1.529130e-18},
	{297, 108, 3, 4, 32, 76, 1, 1.231412, 1.506960e-18},
	{298, 100, 3, 4, 28, 28, 1, 1.231412, 1.397890e-18},
	{299, 120, 3, 4, 50, 50, 1, 1.231412, 1.362180e-18},
	{300, 116, 3, 4, 42, 74, 1, 1.231412, 1.358390e-18},
	{301, 112, 3, 4, 38, 38, 1, 1.231412, 1.357210e-18},
	{302, 128, 3, 4, 56, 56, 1, 1.231412, 9.742990e-19},
	{303, 112, 3, 4, 36, 76, 1, 1.231412, 9.147070e-19},
	{304, 104, 3, 4, 32, 32, 1, 1.231412, 8.721800e-19},
	{305, 108, 3, 4, 32, 32, 1, 1.231412, 7.534800e-19},
	{306, 124, 3, 4, 48, 48, 1, 1.273354, 7.478050e-19},
	{307, 116, 3, 4, 42, 42, 1, 1.273354, 6.791960e-19},
	{308, 120, 3, 4, 48, 48, 1, 1.273354, 6.679650e-19},
	{309, 124, 3, 4, 46, 78, 1, 1.273354, 6.204930e-19},
	{310, 120, 3, 4, 46, 74, 1, 1.273354, 5.926870e-19},
	{311, 128, 3, 4, 54, 54, 1, 1.273354, 5.530780e-19},
	{312, 104, 3, 4, 30, 74, 1, 1.273354, 5.388680e-19},
	{313, 116, 3, 4, 40, 76, 1, 1.273354, 4.990110e-19},
	{314, 112, 3, 4, 36, 36, 1, 1.273354, 4.573530e-19},
	{315, 108, 3, 4, 30, 78, 1, 1.273354, 4.570100e-19},
	{316, 124, 3, 4, 46, 46, 1, 1.273354, 3.102460e-19},
	{317, 120, 3, 4, 46, 46, 1, 1.273354, 2.963440e-19},
	{318, 112, 3, 4, 34, 78, 1, 1.273354, 2.927770e-19},
	{319, 128, 3, 4, 52, 52, 1, 1.273354, 2.821290e-19},
	{320, 104, 3, 4, 30, 30, 1, 1.273354, 2.694340e-19},
	{321, 116, 3, 4, 40, 40, 1, 1.273354, 2.495060e-19},
	{322, 120, 3, 4, 44, 76, 1, 1.273354, 2.405750e-19},
	{323, 124, 3, 4, 44, 80, 1, 1.273354, 2.381080e-19},
	{324, 108, 3, 4, 30, 30, 1, 1.273354, 2.285050e-19},
	{325, 116, 3, 4, 38, 78, 1, 1.273354, 1.717170e-19},
	{326, 104, 3, 4, 28, 76, 1, 1.273354, 1.613020e-19},
	{327, 112, 3, 4, 34, 34, 1, 1.273354, 1.463890e-19},
	{328, 128, 3, 4, 50, 50, 1, 1.273354, 1.305320e-19},
	{329, 120, 3, 4, 44, 44, 1, 1.273354, 1.202870e-19},
	{330, 124, 3, 4, 44, 44, 1, 1.273354, 1.190540e-19},
	{331, 128, 3, 4, 48, 80, 1, 1.307736, 1.106170e-19},
	{332, 120, 3, 4, 42, 78, 1, 1.307736, 9.035000e-20},
	{333, 112, 3, 4, 32, 80, 1, 1.307736, 9.008150e-20},
	{334, 116, 3, 4, 38, 38, 1, 1.307736, 8.585840e-20},
	{335, 124, 3, 4, 42, 82, 1, 1.307736, 8.539640e-20},
	{336, 104, 3, 4, 28, 28, 1, 1.307736, 8.065090e-20},
	{337, 116, 3, 4, 36, 80, 1, 1.307736, 5.599310e-20},
	{338, 128, 3, 4, 48, 48, 1, 1.307736, 5.530870e-20},
	{339, 120, 3, 4, 42, 42, 1, 1.307736, 4.517500e-20},
	{340, 112, 3, 4, 32, 32, 1, 1.307736, 4.504080e-20},
	{341, 128, 3, 4, 46, 82, 1, 1.307736, 4.334810e-20},
	{342, 124, 3, 4, 42, 42, 1, 1.307736, 4.269820e-20},
	{343, 120, 3, 4, 40, 80, 1, 1.307736, 3.174430e-20},
	{344, 124, 3, 4, 40, 84, 1, 1.307736, 2.891760e-20},
	{345, 116, 3, 4, 36, 36, 1, 1.307736, 2.799660e-20},
	{346, 112, 3, 4, 30, 82, 1, 1.307736, 2.695640e-20},
	{347, 128, 3, 4, 46, 46, 1, 1.307736, 2.167410e-20},
	{348, 116, 3, 4, 34, 82, 1, 1.307736, 1.749650e-20},
	{349, 120, 3, 4, 40, 40, 1, 1.307736, 1.587220e-20},
	{350, 128, 3, 4, 44, 84, 1, 1.307736, 1.586440e-20},
	{351, 124, 3, 4, 40, 40, 1, 1.307736, 1.445880e-20},
	{352, 112, 3, 4, 30, 30, 1, 1.307736, 1.347820e-20},
	{353, 120, 3, 4, 38, 82, 1, 1.307736, 1.054790e-20},
	{354, 124, 3, 4, 38, 86, 1, 1.307736, 9.338320e-21},
	{355, 116, 3, 4, 34, 34, 1, 1.345734, 8.748240e-21},
	{356, 128, 3, 4, 44, 44, 1, 1.345734, 7.932220e-21},
	{357, 128, 3, 4, 42, 86, 1, 1.345734, 5.474680e-21},
	{358, 116, 3, 4, 32, 84, 1, 1.345734, 5.297010e-21},
	{359, 120, 3, 4, 38, 38, 1, 1.345734, 5.273960e-21},
	{360, 124, 3, 4, 38, 38, 1, 1.345734, 4.669160e-21},
	{361, 120, 3, 4, 36, 84, 1, 1.345734, 3.349800e-21},
	{362, 124, 3, 4, 36, 88, 1, 1.345734, 2.903950e-21},
	{363, 128, 3, 4, 42, 42, 1, 1.345734, 2.737340e-21},
	{364, 116, 3, 4, 32, 32, 1, 1.345734, 2.648510e-21},
	{365, 128, 3, 4, 40, 88, 1, 1.345734, 1.798290e-21},
	{366, 120, 3, 4, 36, 36, 1, 1.345734, 1.674900e-21},
	{367, 124, 3, 4, 36, 36, 1, 1.345734, 1.451970e-21},
	{368, 120, 3, 4, 34, 86, 1, 1.345734, 1.027330e-21},
	{369, 128, 3, 4, 40, 40, 1, 1.345734, 8.991430e-22},
	{370, 124, 3, 4, 34, 90, 1, 1.345734, 8.779540e-22},
	{371, 128, 3, 4, 38, 90, 1, 1.345734, 5.674390e-22},
	{372, 120, 3, 4, 34, 34, 1, 1.345734, 5.136640e-22},
	{373, 124, 3, 4, 34, 34, 1, 1.345734, 4.389770e-22},
	{374, 120, 3, 4, 32, 88, 1, 1.345734, 3.073590e-22},
	{375, 128, 3, 4, 38, 38, 1, 1.345734, 2.837200e-22},
	{376, 128, 3, 4, 36, 92, 1, 1.345734, 1.735640e-22},
	{377, 120, 3, 4, 32, 32, 1, 1.345734, 1.536800e-22},
	{378, 128, 3, 4, 36, 36, 1, 1.345734, 8.678180e-23},
	{379, 128, 3, 4, 34, 94, 1, 1.345734, 5.192020e-23},
	{380, 128, 3, 4, 34, 34, 1, 1.345734, 2.600000e-23},
}

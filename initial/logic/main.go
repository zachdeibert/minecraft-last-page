package logic

import (
	"fmt"
	"math"

	"../gpu"
)

// Run the program
func Run() {
	sieve, e := gpu.CreateSieve()
	if e != nil {
		panic(e)
	}
	defer sieve.Close()
	out, e := createOutputWriter(outputFile)
	if e != nil {
		panic(e)
	}
	defer out.Close()
	for i := minStringLen; i <= maxStringLen; i++ {
		fmt.Printf("\rStatus (len=%d):  0.00%%", i)
		prefixLen := i - baseLen - suffixLen
		total := math.Pow(float64(len(alphabet)), float64(prefixLen))
		var done float64 = 0
		if e = sieve.Configure(prefixLen, baseLen, suffixLen, outBuf); e != nil {
			panic(e)
		}
		for p := createPrefixer(prefixLen); p.hasNext(); {
			res, e := sieve.Run(p.next())
			if e != nil {
				panic(e)
			}
			out.Write(res)
			done++
			fmt.Printf("\rStatus (len=%d): % 5.2f%%", i, done/total*100)
		}
	}
}

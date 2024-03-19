package main

import (
	gcdlcm "justbeyourselfandenjoy/RSA/GCD_LCM"
	eratosthenessieve "justbeyourselfandenjoy/RSA/eratosthenes_sieve"
	eulerssieve "justbeyourselfandenjoy/RSA/eulers_sieve"
	"justbeyourselfandenjoy/RSA/factor"
	fastexp "justbeyourselfandenjoy/RSA/fast_exp"
	"justbeyourselfandenjoy/RSA/primality"
	"justbeyourselfandenjoy/RSA/rsa"
)

func main() {
	gcdlcm.GCDLCMRun()
	fastexp.FastExtRun()
	eratosthenessieve.EratosthenesSieveRun()
	eulerssieve.EulersSieveRun()
	factor.FactorNumbersRun()
	primality.PrimalityTestRun()
	rsa.RSARun()
}

package main

import (
	gcdlcm "justbeyourselfandenjoy/RSA/GCD_LCM"
	eratosthenessieve "justbeyourselfandenjoy/RSA/eratosthenes_sieve"
	fastexp "justbeyourselfandenjoy/RSA/fast_exp"
)

func main() {
	gcdlcm.GCDLCMRun()
	fastexp.FastExtRun()
	eratosthenessieve.EratosthenesSieveRun()
}

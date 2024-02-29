package main

import (
	gcdlcm "justbeyourselfandenjoy/RSA/GCD_LCM"
	eratosthenessieve "justbeyourselfandenjoy/RSA/eratosthenes_sieve"
	eulerssieve "justbeyourselfandenjoy/RSA/eulers_sieve"
	fastexp "justbeyourselfandenjoy/RSA/fast_exp"
)

func main() {
	gcdlcm.GCDLCMRun()
	fastexp.FastExtRun()
	eratosthenessieve.EratosthenesSieveRun()
	eulerssieve.EulersSieveRun()
}

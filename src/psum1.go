// Example go program which sums a bunch of prime numbers
// This is done allegedly in parallel

package main
import "fmt"
import "runtime"

func sum(start int64, end int64, ch chan int64) {
	var s int64 = 0
	var i int64 = start
	for ; i < end; i++ {
		if(isprime(i)) {
			s+= i
		}
	}
	ch <- s
}

func isprime(n int64) bool {
	if(n % 2 == 0) {
		return false;
	}
	var i int64 = 3
	for ; i < (n/2) ; i++ {
		if(n % i == 0) {
			return false
		}
	}
	return true
}


func main() {
	NPROCS := 2
	runtime.GOMAXPROCS(NPROCS)
	var n int64 = 300000

	ch1 := make(chan int64)
	ch2 := make(chan int64)
	
	go sum(0, n/2, ch1)
	go sum(n/2, n, ch2)
	
	sum1 := <-ch1
	sum2 := <-ch2

	result := sum1 + sum2
	
	fmt.Printf("%d \n",result)
		
}
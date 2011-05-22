package main
import "fmt"

func sum(start int64, end int64) int64 {
	var s int64 = 0
	var i int64 = start
	for ; i < end; i++ {
		if(isprime(i)) {
			s+= i
		}
	}
	return s
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
	result := sum(0, 300000)
	fmt.Printf("%d \n",result)
}
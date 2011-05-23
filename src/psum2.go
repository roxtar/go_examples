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
	NPROCS := 4
	runtime.GOMAXPROCS(NPROCS)
	var n int64 = 100000
	var block_size int64 = n/int64(NPROCS)

	ch := make(chan int64)

	for i:=0; i < NPROCS; i++ {
		go sum(int64(i)*block_size, int64(i+1)*block_size, ch)
	}

	var result int64  = 0
	
	for i:=0; i < NPROCS; i++ {
		result += <- ch 
	}

	fmt.Printf("%d \n",result)
		
}
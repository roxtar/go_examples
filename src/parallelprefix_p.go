package main;

import (
	"fmt"
	"runtime"
)

const N = 30000000
const NPROCS = 4
var a [N]int64

func doInit() {
	for i :=0; i<len(a); i++ {
		a[i] = int64(i)
	}
}

func doParallelPrefixMain() {
	var ch [NPROCS]chan int64
	var finish chan bool
	
	//Initialize all channels
	finish = make(chan bool, NPROCS)
	for i:= 0; i<NPROCS; i++ {
		ch[i] = make(chan int64, NPROCS)
	}
	
	//Launch all threads
	for i:=0; i<NPROCS; i++ {
		go threadParallelPrefix(i, ch[:], finish)
	}
	
	//Equivalent to a finish
	for i:=0; i<NPROCS; i++ {
		<- finish
	}
	fmt.Printf("%v\n", a[N-1])
	
}

func threadParallelPrefix(pid int, ch []chan int64, finish chan bool) {
	//First do a local parallel prefix computation
	var startIndex int64 = int64(pid) * N/NPROCS
	var endIndex int64 = startIndex + (N/NPROCS)
	var toBeAdded int64 = 0

	for i:=startIndex+1; i<endIndex; i++ {		
		a[i] = a[i] + a[i-1]
	}
	
	for i:= pid+1; i<NPROCS; i++ {		
		ch[i] <- a[i-1]
	}
	
	for i:=0; i < pid; i++ {
		 toBeAdded += <- ch[pid]
	}
			
	for j:=startIndex; j<endIndex; j++ {
		a[j] += toBeAdded
	}

	 
	finish <- true
	
}  
func printArray() {
for i:=0; i<len(a); i++ {
	fmt.Println("a[",i,"] = ",a[i]);
}
}

func main() {
	runtime.GOMAXPROCS(int(NPROCS))
	doInit()
	doParallelPrefixMain()
}
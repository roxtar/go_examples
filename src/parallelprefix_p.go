package main;

import (
	"fmt"
	"runtime"
)

const N int64 = 30000000
const NPROCS int64 = 2
var a [N]int64

func doInit() {
	for i :=0; i<len(a); i++ {
		a[i] = 1
	}
}

func doParallelPrefixMain() {
	var ch [NPROCS]chan int64
	var finish chan bool
	
	//Initialize all channels
	finish = make(chan bool, NPROCS)
	for i:= int64(0); i<NPROCS; i++ {
		ch[i] = make(chan int64, NPROCS)
	}
	
	//Launch all threads
	for i:=int64(0); i<NPROCS; i++ {
		go threadParallelPrefix(i, ch, finish)
	}
	
	//Equivalent to a reduce
	for i:=int64(0); i<NPROCS; i++ {
		<- finish
	}
	
}

func threadParallelPrefix(pid int64, ch [NPROCS]chan int64, finish chan bool) {
	//First do a local parallel prefix computation
	var startIndex int64 = pid*N/NPROCS
	var endIndex int64 = startIndex + (N/NPROCS)
	var localSum int64 = a[startIndex]
	for i:=startIndex+1; i<endIndex; i++ {
		localSum += a[i]
		a[i] = localSum
	}
	
	for i:= pid+1; i<NPROCS; i++ {
		
		ch[i] <- localSum
	}
	
	for i := int64(0); i<pid; i++ {
		toBeAdded := <- ch[pid]
		for j:=startIndex; j<endIndex; j++ {
			a[j] += toBeAdded
		}
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
//	printArray();
// 	fmt.Println("Hello world!!")
}
// Sequential version of prefix sum
package main
import ("fmt")

func prefix_sum(a[] int64) {
	for i:=1; i < len(a); i++ {
		a[i] = a[i] + a[i-1] 
	}
}

func main() {
	const n = 300000
	a := make([] int64, n)
	for i, _ := range a {
		a[i] = int64(i)
	}
	prefix_sum(a)
	fmt.Printf("%v \n", a[n-1])
}
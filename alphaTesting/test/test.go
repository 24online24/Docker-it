package main

import "fmt"

func main() {
	ch := make(chan int)
	go numbers(ch)
	for x := range ch {
		fmt.Println(x)
	}
}

func numbers(ch chan int) {
	x := 0
	for {
		ch <- x
		x++
		fmt.Print("Mau!")
	}
}

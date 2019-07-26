package main

import (
	"fmt"
	"time"
)

func gen() <-chan int {
	c := make(chan int)
	go func() {
		defer close(c)	// must close, otherwise got deadlock because this channel will readonly
		for i := 1; i <= 100; i++ {
			c <- i
		}
	}()
	return c
}

func do(n int, i chan int) {
	time.Sleep(10 * time.Millisecond)

	fmt.Println("got: ", n)
	<- i
}

func main() {
	nums := gen()
	i := make(chan int, 10)
	for num := range nums {
		go func(){
			i <- num
		}()
		go do(num, i)
	}
	time.Sleep(2 * time.Second)
}

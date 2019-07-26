package main

import (
	"fmt"
	"time"
)

func spread(main <-chan int, output1, output2, output3 chan int) {
	for number := range main {
		output1 <- number
		output2 <- number
		output3 <- number
	}
}

//func printChl(chn chan int) {
//	for i := range chn {
//		fmt.Println(i)
//	}
//}

//func main() {
//	main := make(chan int)
//	a := make(chan int)
//	b := make(chan int)
//	c := make(chan int)
//	go spread(main, a, b, c)
//	go printChl(a)
//	go printChl(b)
//	go printChl(c)
//	for i := 1; i <= 10; i++ {
//		main <- i
//		time.Sleep(1 * time.Second)
//	}
//	main <- 1
//	//close(main)
//	_, _ = fmt.Scanln()
// }

func main() {
	source := make(chan int)
	a := make(chan int)
	b := make(chan int)
	c := make(chan int)
	output := make(chan int)
	go spread(source, a, b, c)

	go centralize(a, output)
	go centralize(b, output)
	go centralize(c, output)
	//go printChannel(a, func(i int) int {
	//	return i * 2
	//})
	//go printChannel(b, func(i int) int {
	//	return i * 3
	//})
	//go printChannel(c, func(i int) int {
	//	return i * 4
	//})
	go fanIn(output)

	for i := 1; i <= 2; i++ {
		source <- i
		//time.Sleep(1 * time.Second)
	}

	time.Sleep(2 * time.Second)
	close(output)
	time.Sleep(2 * time.Second)
}

// Refactor
//func printChannel(chn chan int, f func(int) int) {
//	for i := range chn {
//		value := f(i)
//		fmt.Println(value)
//	}
//}

func centralize(in, out chan int) {
	for number := range in {
		out <- number
	}
}

func fanIn(output chan int) int {
	var total int
	for num := range output {
		total += num
		fmt.Println("Total =", total)

	}
	//close(output)
	fmt.Println("Final total = ", total)
	return total
}

// Does not work, since the message cannot be consumed
//func fanIn(a, b, c chan int) int {
//	var total int
//	select {
//	case i := <-a:
//		total = total + i
//	case i := <-b:
//		total = total + i
//	case i := <-c:
//		total = total + i
//	}
//	fmt.Println("Total =", total)
//	return total
//}

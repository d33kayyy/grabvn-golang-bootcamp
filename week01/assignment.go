package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func isOperator(text string) (op string, err error) {
	switch text {
	case "/", "*", "+", "-":
		op = text
		return op, nil
	}
	return op, errors.New(fmt.Sprintf("Invalid operator: '%s'", text))
}

func isFloat(text string) (value float64, err error) {
	value, err = strconv.ParseFloat(text, 64)
	return
}

func validate(text string) (a float64, b float64, op string, err error) {
	words := strings.Fields(text)
	if len(words) != 3 {
		err = errors.New(fmt.Sprintf("Invalid input: '%s'", text))
		return
	}
	a, err = isFloat(words[0])
	if err != nil {
		return
	}
	op, err = isOperator(words[1])
	if err != nil {
		return
	}
	b, err = isFloat(words[2])
	if err != nil {
		return
	}
	return
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func eval(text string) (err error) {
	a, b, op, err := validate(text)
	if err != nil {
		return err
	}
	result := 0.0
	switch op {
	case "+":
		result = a + b
	case "-":
		result = a - b
	case "*":
		result = a * b
	case "/":
		if b == 0 {
			return errors.New("Cannot divide by zero")
		}
		result = a / b
	}
	fmt.Printf("%d %s %d = %s",
		int(a), op, int(b), formatFloat(result),
	)
	fmt.Println()
	return
}

func main() {
	//fmt.Println(isOperator("+"))
	//fmt.Println(isOperator("/"))
	//fmt.Println(isOperator("-"))
	//fmt.Println(isOperator("*"))
	//fmt.Println(isOperator("x"))

	//fmt.Println(validate("1 + 2"))
	//fmt.Println(validate("3 / 4"))
	//fmt.Println(validate("10 * 2"))
	//fmt.Println(validate("12 + 24"))
	//fmt.Println(validate("12 ~ 24"))
	//fmt.Println(validate("1 2 3 + 24"))

	//_ = eval("1 + 2")
	//_ = eval("3 / 4")
	//_ = eval("10 * 2")
	//_ = eval("1 + 2")
	//_ = eval("12 + 24")
	//fmt.Println(eval("12 ~ 24"))
	//fmt.Println(eval("12 / 0"))
	//fmt.Println(eval("1 2 3 + 24"))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("> ")
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			fmt.Println("Goodbye :)")
			break
		}
		err := eval(text)
		if err != nil {
			fmt.Println("Got error: " + err.Error())
		}
		fmt.Print("> ")
	}
}

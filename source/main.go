// program to demonstrate that init is called before main
package main

import "fmt"

func init() {
	fmt.Println("init is called")
}

func main() {
	fmt.Println("main is called")
	fmt.Printf(fmt.Sprint(len(leftHalf)))

	for index, element := range leftHalf {
		fmt.Println(index, element)
	}
}

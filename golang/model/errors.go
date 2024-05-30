package model

import "fmt"

type Errors []error

func (es Errors) Print() {
	fmt.Println("Errors:")
	for _, e := range es {
		fmt.Println(e)
	}
}

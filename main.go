package main

import (
	"fmt"

	"github.com/climbcomp/climbcomp-go/climbcomp"
)

func main() {
	msg := fmt.Sprintf("Climbcomp API v%v", climbcomp.VERSION)
	fmt.Println(msg)
}

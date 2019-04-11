package main

import (
	"fmt"

	"github.com/Mot93/personalfinance/financecalc"
)

// TODO: manage the panic lanched by the db

func main() {
	result, err := financecalc.Total()
	if err == nil {
		fmt.Printf("Totale = %v\n", result)
	} else {
		fmt.Printf("ERROR = %v\n", err)
	}
}

package main

import (
	"github.com/rosbit/go-deno"
	"os"
	"fmt"
	"encoding/json"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf("Usage: %s <deno-exe> <js-file> <funcName>[ <args>...]\n", os.Args[0])
		return
	}

	// 1. init deno
	ctx, err := deno.NewDeno(os.Args[1], os.Args[2])
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	defer ctx.Quit()
	// fmt.Printf("init ok\n")

	// 2. prepare arguments and variables
	funcName := os.Args[3]
	args := make([]interface{}, len(os.Args) - 4)
	for i,j:=4,0; i<len(os.Args); i,j=i+1,j+1 {
		if err = json.Unmarshal([]byte(os.Args[i]), &args[j]); err != nil {
			args[j] = os.Args[i]
		}
	}

	// 3. call funcName with arguments
	res, err := ctx.CallFunc(funcName, args...)

	if err != nil {
		fmt.Printf("failed to call %s: %v\n", funcName, err)
		return
	}
	fmt.Printf("res: %v\n", res)

	// --- xxx ---
	var getA func(map[string]interface{})interface{}
	if err = ctx.BindFunc("getA", &getA); err != nil {
		fmt.Printf("failed to BindFunc add: %s\n", err)
		return
	}
	fmt.Printf("getA({a: 10}): %v\n", getA(map[string]interface{}{"a": 10}))
}


package main

import (
	"fmt"

	_ "github.com/tarantool/go-tarantool"
	// _ "github.com/fl00r/go-tarantool-1.6"
)

func main() {
	fmt.Println("Hello 1")
	// opts := tarantool.Opts{User: "admin", Pass: "pass"}
	// conn, err := tarantool.Connect("127.0.0.1:3301", opts)
	// if err != nil {
	// 	panic(err)
	// }
	// defer conn.Close()
	// resp, err := conn.Insert(999, []interface{}{99999, "BB"})
	// if err != nil {
	// 	fmt.Println("Error", err)
	// 	fmt.Println("Code", resp.Code)
	// }
	fmt.Println("Hello 2")
}

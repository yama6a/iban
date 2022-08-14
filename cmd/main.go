package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("hi")
	for i := 1; i <= 10; i++ {
		fmt.Println(fmt.Sprintf("%02d/10 I'm still alive...", i))
		time.Sleep(time.Second)
	}
	fmt.Println("bye")
}

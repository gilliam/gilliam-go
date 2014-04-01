package main

import "github.com/gilliam/gilliam-go"
import "fmt"

func main() {
    client := gilliam.New()
    client.FormationInstances("scheduler")
	fmt.Printf("Hello, world\n");
}

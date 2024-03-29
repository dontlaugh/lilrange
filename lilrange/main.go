package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/dontlaugh/lilrange"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("must provide one argument: a lilrange string")
	}

	r, err := lilrange.Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	now := time.Now().In(time.UTC)

	fmt.Printf("now (nano precision): %v\n", now)
	fmt.Printf("start:                %v\n", r.Start)
	fmt.Printf("end:                  %v\n", r.End)
	fmt.Printf("duration:             %v\n", r.Duration)
	fmt.Printf("within range now?:    %v\n", r.Within(now))
	fmt.Println("\nThe next range's values")
	next := r.Next()
	fmt.Printf("next start:           %v\n", next.Start)
	fmt.Printf("next end:             %v\n", next.End)
	fmt.Printf("duration (same):      %v\n", next.Duration)
	fmt.Printf("until next start:     %v\n", time.Until(next.Start).Truncate(time.Second))
}

// +build ignore

package main

import (
	"fmt"

	"github.com/wmentor/ljatom"
)

func main() {

	for msg := range ljatom.Read() {

		fmt.Println("Time: " + msg.Created.String())
		fmt.Println("Journal: " + msg.Journal)
		fmt.Println("Post URL: " + msg.Url)
		fmt.Println("Post title: " + msg.Title)
		fmt.Print("Post body: " + msg.Content + "\n\n\n")

	}
}

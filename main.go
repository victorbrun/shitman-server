package main

import "fmt"

func main() {
	//http.HandleFunc("/ws", handleConnection)
	//log.Println("Starting server on :8080")
	//log.Fatal(http.ListenAndServe(":8080", nil))

	deck := NewDeck()
	fmt.Println(deck)
	hand := deck.Draw(20)
	fmt.Println(deck)
	fmt.Println(hand)
}

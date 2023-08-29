package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Guesser interface {
	Guess(done <-chan bool, turn chan<- bool) int
}

type RandomGuesser struct {
	Min int
	Max int
}

func (rg *RandomGuesser) Guess(done <-chan bool, turn chan<- bool) int {
	select {
	case <-done:
		return -1
	default:
		guess := rand.Intn(rg.Max-rg.Min+1) + rg.Min
		turn <- true // Signal it's the other guesser's turn
		return guess
	}
}

type MethodicalGuesser struct {
	Min          int
	Max          int
	CurrentGuess int
}

func (mg *MethodicalGuesser) Guess(done <-chan bool, turn chan<- bool) int {
	select {
	case <-done:
		return -1
	default:
		if mg.CurrentGuess == 0 {
			mg.CurrentGuess = mg.Min
		} else {
			mg.CurrentGuess++
		}
		turn <- true // Signal it's the other guesser's turn
		return mg.CurrentGuess
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())

	Min, Max := 1, 10
	correctValue := rand.Intn(Max-Min+1) + Min

	fmt.Println("The correct value is:", correctValue)

	randomGuesser := &RandomGuesser{Min, Max}
	methodicalGuesser := &MethodicalGuesser{Min, Max, 0}

	guesses := make(chan string)
	done := make(chan bool)
	turn := make(chan bool)

	// Goroutine for RandomGuesser
	go func() {
		for {
			<-turn // Wait for its turn
			guess := randomGuesser.Guess(done, turn)
			if guess != -1 {
				guesses <- fmt.Sprintf("RandomGuesser guessed: %d", guess)
			}
		}
	}()

	// Goroutine for MethodicalGuesser
	go func() {
		for {
			<-turn // Wait for its turn
			guess := methodicalGuesser.Guess(done, turn)
			if guess != -1 {
				guesses <- fmt.Sprintf("MethodicalGuesser guessed: %d", guess)
			}
		}
	}()

	turn <- true // Start the guessing by signaling the first guesser

	var winner string
	for {
		message := <-guesses
		fmt.Println(message)

		if message == fmt.Sprintf("RandomGuesser guessed: %d", correctValue) {
			winner = "RandomGuesser"
			break
		} else if message == fmt.Sprintf("MethodicalGuesser guessed: %d", correctValue) {
			winner = "MethodicalGuesser"
			break
		}
	}

	close(done) // Signal to goroutines to stop
	close(turn) // Close the turn channel
	fmt.Println("The winner is:", winner)
}

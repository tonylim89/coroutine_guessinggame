# Guessing Game in Go

In this game, we have two guessers: a RandomGuesser and a MethodicalGuesser. They take turns trying to guess a random number, with the first one to guess correctly being the winner.

## Structs and Interface
```
type Guesser interface {
	Guess(done <-chan bool, turn chan<- bool) int
}

type RandomGuesser struct {
	Min int
	Max int
}

type MethodicalGuesser struct {
	Min          int
	Max          int
	CurrentGuess int
}
```

Guesser: An interface that any guesser should implement. It mandates a Guess function.

RandomGuesser: Represents a guesser that makes random guesses.

MethodicalGuesser: Represents a guesser that guesses in a methodical manner, starting from the minimum value and increasing by one with each guess.

## Guessing Strategies

### RandomGuesser

```
func (rg *RandomGuesser) Guess(done <-chan bool, turn chan<- bool) int {
	select {
	case <-done:
		return -1
	default:
		guess := rand.Intn(rg.Max-rg.Min+1) + rg.Min
		turn <- true
		return guess
	}
}
```
This function tries to make a guess:

If the game is already finished (done channel is closed), it returns -1.

Otherwise, it makes a random guess between Min and Max.

### MethodicalGuesser

```
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
		turn <- true
		return mg.CurrentGuess
	}
}
```
This function tries to make a methodical guess:

If the game is already finished (done channel is closed), it returns -1.

Otherwise, it guesses the next number methodically, starting from Min.

## Main Program Execution

```
func main() {
	rand.Seed(time.Now().UnixNano())

	Min, Max := 1, 10
	correctValue := rand.Intn(Max-Min+1) + Min
  //...
}
```

Here we:

Seed the random number generator.

Define the minimum and maximum values for guessing.

Determine the correct value to be guessed.

## Game Loop and Concurrency

To manage the game, we utilize goroutines channels:

```
guesses := make(chan string)
done := make(chan bool)
turn := make(chan bool)
```

guesses channel: Carries messages about the guesses made.

done channel: Signifies when the game is over.

turn channel: Used to coordinate turns between guessers.

Goroutines take turns making guesses, and their guesses are sent to the guesses channel. The main thread then checks these guesses to determine if there's a winner.

## Thoughts
Initially I felt that it was redundant to use coroutines if we are still going to run things in sequence, however, the biggest use is that each thread could be pre-processed up till the waiting part and this would speed up things immensely where sequences are ran one after the other instead of starting the processing when it reaches the point of hand over.

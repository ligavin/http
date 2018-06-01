package main

import (
	"fmt"
	"math/rand"
	"time"
	"os"
	"http/comm"
)

func main() {
	rand.Seed(42) // Try changing this number!
	answers := []string{
		"It is certain",
		"It is decidedly so",
		"Without a doubt",
		"Yes definitely",
		"You may rely on it",
		"As I see it yes",
		"Most likely",
		"Outlook good",
		"Yes",
		"Signs point to yes",
		"Reply hazy try again",
		"Ask again later",
		"Better not tell you now",
		"Cannot predict now",
		"Concentrate and ask again",
		"Don't count on it",
		"My reply is no",
		"My sources say no",
		"Outlook not so good",
		"Very doubtful",
	}
	fmt.Println("Magic 8-Ball says:", answers[rand.Intn(len(answers))])
	// Output: Magic 8-Ball says: As I see it yes

	for i := 0; i < 10; i++{
		fmt.Printf("time:%d，rand:%d,pid:%d\n", time.Now().UnixNano(), comm.GetRandUint32(),os.Getpid())
	}

}
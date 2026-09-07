// 1nder
// CS361

package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	eMax      = 3
	rMax      = 9
	totalElf  = 10
	sleepTime = 5 * time.Second
)

// Santa is a struct that holds the waiting elves and reindeer
type Santa struct {
	sleep chan int
}

// Door checks the door for Santa and wakes him up when necessary
type Door struct {
	elfChan      chan *Elf
	reindeerChan chan *Reindeer
	santa        *Santa
}

// Reindeer is a struct that goes to Santa and waits
type Reindeer struct {
	id   int
	door *Door
	wait chan int
}

// Elf is a struct that goes to Santa and waits
type Elf struct {
	id   int
	door *Door
	wait chan int
}

// Main run function of Santa
func (santa *Santa) run(door *Door, wg *sync.WaitGroup) {
	defer wg.Done()
	var c int
	d := 0
	for {
		d++
		fmt.Printf("\n------------Sleep Cycle: %d-------------\n", d)
		fmt.Print("\nSanta went to sleep.\n\n")
		// Santa goes to sleep
		c = <-santa.sleep
		time.Sleep(sleepTime)
		// If the door sends the message 1, then there are 9 reindeer
		// else if the door sends the message 2, then there are 3 or more elves
		if c == 1 {
			santa.deliverPresents(door)
		} else if c == 2 {
			santa.talkToElves(door)
		}
	}
}

// Santa delivers presents with the reindeer
func (santa *Santa) deliverPresents(door *Door) {
	var r *Reindeer
	fmt.Print("\nSanta is delivering presents with the Reindeer.\n\n")
	time.Sleep(sleepTime)
	// Santa sends a message to the 9 reindeer, saying they can leave him
	for i := 0; i < rMax; i++ {
		r = <-door.reindeerChan
		r.wait <- 1
	}
}

// Santa talks to 3 of the elves
func (santa *Santa) talkToElves(door *Door) {
	var e *Elf
	fmt.Print("\nSanta is speaking with 3 of the elves.\n\n")
	time.Sleep(sleepTime)
	// Santa sends a message to the 3 elves, saying they can leave him
	for i := 0; i < eMax; i++ {
		e = <-door.elfChan
		e.wait <- 1
	}
}

// Thread that is constantly checking the reindeer and elf count at door
// If there are 9 reindeer or at least 3 elves, the door wakes up Santa
func (door *Door) run() {
	for {
		// Checks if the number of reindeer in the channel is 9
		if len(door.reindeerChan) == rMax {
			// Sends a message to Santa to wake up
			door.santa.sleep <- 1
			// Checks if the number of elves in the channel is 3 or more
		} else if len(door.elfChan) >= eMax {
			// Sends a message to Santa to wake up
			door.santa.sleep <- 2
		}
	}
}

// Main function of the Reindeer
func (reindeer *Reindeer) run() {
	for {
		fmt.Printf("[Reindeer] %d went to Santa.\n", reindeer.id)
		// Reindeer sends message to Santa saying they arrived
		reindeer.door.reindeerChan <- reindeer
		// Reindeer waits to recieve a message from Santa so they can leave
		<-reindeer.wait
		fmt.Printf("[Reindeer] %d returned from Santa.\n", reindeer.id)
		time.Sleep(sleepTime)
	}
}

// Main function of the Elves
func (elf *Elf) run() {
	for {
		fmt.Printf("[Elf] %d went to Santa.\n", elf.id)
		// Elf sends message to Santa saying they arrived
		elf.door.elfChan <- elf
		// Elf waits to recieve a message from Santa so they can leave
		<-elf.wait
		fmt.Printf("[Elf] %d returned from Santa.\n", elf.id)
		time.Sleep(sleepTime)
	}
}

func main() {
	elfChan := make(chan *Elf, totalElf)
	reindeerChan := make(chan *Reindeer, rMax)

	santa := Santa{sleep: make(chan int)}

	door := Door{santa: &santa, elfChan: elfChan, reindeerChan: reindeerChan}
	go door.run()

	for i := 0; i < rMax; i++ {
		r := Reindeer{id: i, door: &door, wait: make(chan int)}
		go r.run()
	}

	for i := 0; i < totalElf; i++ {
		e := Elf{id: i, door: &door, wait: make(chan int)}
		go e.run()
	}

	// I needed this wait group so that Santa's thread could run without
	// the main function ending. Otherwise, if I did go santa.run()
	// main() would end and the program would stop immediately
	var wg sync.WaitGroup
	wg.Add(1)
	go santa.run(&door, &wg)
	wg.Wait()
}

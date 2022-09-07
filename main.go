package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numOfPhilosophers = 5
const timesNeededToEat = 3
const eatingTime = 2
const thinkingTime = 1

type philosopher struct {
	philoChannel chan int
	id           int
	ownFork      *fork
	rightFork    *fork
	timesEaten   int
	doneEating   bool
	isThinking   bool
}

type fork struct {
	inUse bool
	id    int
}

func main() {
	//Create forks
	forks := make([]fork, numOfPhilosophers)
	for i := range forks {
		forks[i].id = i
		forks[i].inUse = false
	}
	//create philosophers
	philos := make([]philosopher, numOfPhilosophers)
	for j := range philos {
		philos[j].id = j
		philos[j].ownFork = &forks[j]
		philos[j].rightFork = &forks[(j+1)%numOfPhilosophers]
		philos[j].timesEaten = 0
		philos[j].doneEating = false
		philos[j].philoChannel = make(chan (int), 500)
	}
	//Signals all channels, they are philos are ready to act
	for p := range philos {
		philoPointer := &philos[p]
		philoPointer.philoChannel <- p
	}

	//Repeat routine.
	for {
		//If everyone ate 3 times, break
		if philos[0].doneEating && philos[1].doneEating && philos[2].doneEating && philos[3].doneEating && philos[4].doneEating {
			break
		}
		//Check channel message
		select {
		case msg0 := <-philos[0].philoChannel:
			philoPointer0 := &philos[0]
			go func() {
				time.Sleep(1111 * time.Millisecond)
				philoPointer0.act(msg0)
			}()
		case msg1 := <-philos[1].philoChannel:
			philoPointer1 := &philos[1]
			go func() {
				time.Sleep(789 * time.Millisecond)
				philoPointer1.act(msg1)
			}()
		case msg2 := <-philos[2].philoChannel:
			philoPointer2 := &philos[2]
			go func() {
				time.Sleep(653 * time.Millisecond)
				philoPointer2.act(msg2)
			}()
		case msg3 := <-philos[3].philoChannel:
			philoPointer3 := &philos[3]
			go func() {
				time.Sleep(526 * time.Millisecond)
				philoPointer3.act(msg3)
			}()
		case msg4 := <-philos[4].philoChannel:
			philoPointer4 := &philos[4]
			go func() {
				time.Sleep(367 * time.Millisecond)
				philoPointer4.act(msg4)
			}()
		}
	}
	fmt.Println("Everyone is full and ate:", timesNeededToEat, " times")
}

func (p *philosopher) act(msg int) {
	(*p).ownFork.inUse = true
	//Eating
	if (*p).rightFork.inUse == false && (*p).timesEaten < timesNeededToEat {
		go func() { //IKKE SIKKER PÅ OM DEN SKAL GO HER
			(*p).isThinking = false
			(*p).rightFork.inUse = true
			(*p).timesEaten++
			fmt.Println("philo: ", (*p).id, " is eating. Times eaten:", (*p).timesEaten)
			if (*p).timesEaten == 3 {
				(*p).doneEating = true
				fmt.Println("philo: ", (*p).id, " is is full of food")
			}
			randomPause(eatingTime)
			(*p).ownFork.inUse = false
			(*p).rightFork.inUse = false
			//fmt.Println("philo: ", (*p).id, " is done eating")
		}()
		//Thinking
	} else {
		go func() { //IKKE SIKKER PÅ OM DEN SKAL GO HER
			if !(*p).isThinking {
				fmt.Println("philo: ", (*p).id, " is thinking.")
			}
			(*p).ownFork.inUse = false
			randomPause(thinkingTime)
			(*p).isThinking = true
		}()
	}
	(*p).philoChannel <- (*p).id
}

func randomPause(max int) { //Fra https://github.com/iokhamafe/Golang/blob/master/diningphilosophers.go
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max*1000)+100))
}

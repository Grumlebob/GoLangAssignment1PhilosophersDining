package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numOfPhilosophers = 5
const timesNeededToEat = 3
const eatingTime = 2

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
	forkChannel  chan bool //fork ready = true. Not ready = empty
	forkRunAgain chan bool //Repeat search for own fork
}

func main() {
	//Create forks
	forks := make([]fork, numOfPhilosophers)
	for i := range forks {
		forks[i].forkChannel = make(chan (bool), 1)
		forks[i].forkRunAgain = make(chan (bool), 1)
		forks[i].forkChannel <- true
		forks[i].forkRunAgain <- true
	}
	//create philosophers
	philos := make([]philosopher, numOfPhilosophers)
	for j := range philos {
		philos[j].id = j
		philos[j].ownFork = &forks[j]
		philos[j].rightFork = &forks[(j+1)%numOfPhilosophers]
		philos[j].timesEaten = 0
		philos[j].doneEating = false
		philos[j].philoChannel = make(chan (int), 1)
		philos[j].philoChannel <- j
	}

	//Repeat routine.
	for {
		//If everyone ate 3 times exit loop
		if philos[0].doneEating && philos[1].doneEating && philos[2].doneEating && philos[3].doneEating && philos[4].doneEating {
			break
		}
		//Start Fork and Philo routines. If fork is ready, signal philosopher.
		select {
		case <-forks[0].forkRunAgain:
			philoPointer := &philos[0]
			go philoPointer.goRoutineForForks()
		case <-forks[1].forkRunAgain:
			philoPointer := &philos[1]
			go philoPointer.goRoutineForForks()
		case <-forks[2].forkRunAgain:
			philoPointer := &philos[2]
			go philoPointer.goRoutineForForks()
		case <-forks[3].forkRunAgain:
			philoPointer := &philos[3]
			go philoPointer.goRoutineForForks()
		case <-forks[4].forkRunAgain:
			philoPointer := &philos[4]
			go philoPointer.goRoutineForForks()
		case <-philos[0].philoChannel:
			philoPointer0 := &philos[0]
			go func() {
				time.Sleep(413 * time.Millisecond)
				philoPointer0.GoRoutinePhiloEatOrThink()
			}()
		case <-philos[1].philoChannel:
			philoPointer1 := &philos[1]
			go func() {
				time.Sleep(348 * time.Millisecond)
				philoPointer1.GoRoutinePhiloEatOrThink()
			}()
		case <-philos[2].philoChannel:
			philoPointer2 := &philos[2]
			go func() {
				time.Sleep(175 * time.Millisecond)
				philoPointer2.GoRoutinePhiloEatOrThink()
			}()
		case <-philos[3].philoChannel:
			philoPointer3 := &philos[3]
			go func() {
				time.Sleep(526 * time.Millisecond)
				philoPointer3.GoRoutinePhiloEatOrThink()
			}()
		case <-philos[4].philoChannel:
			philoPointer4 := &philos[4]
			go func() {
				time.Sleep(367 * time.Millisecond)
				philoPointer4.GoRoutinePhiloEatOrThink()
			}()
		}
	}
	fmt.Println("Everyone is full and ate:", timesNeededToEat, "times")
}

func randomPause(max int) { //Fra https://github.com/iokhamafe/Golang/blob/master/diningphilosophers.go
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max*1000)+50))
}

func (p *philosopher) goRoutineForForks() {
	waitingTime := make(chan bool, 1)
	go func() {
		randomPause(1)
		waitingTime <- true
	}()

	select {
	case <-(*p).ownFork.forkChannel: //Own fork is ready - announce to philosopher to attempt to eat.
		(*p).ownFork.forkChannel <- true
		(*p).philoChannel <- (*p).id
	case <-waitingTime: //Own fork not ready, try again.
		(*p).ownFork.forkRunAgain <- true
	}

}

func (p *philosopher) GoRoutinePhiloEatOrThink() {
	waitingTime := make(chan bool, 1)
	go func() {
		randomPause(2)
		waitingTime <- true
	}()
	//We got a channel msg that our own fork is ready, therefore check if other fork is ready, if so eat otherwise think.
	select {
	case <-(*p).rightFork.forkChannel:
		<-(*p).ownFork.forkChannel //pickup both forks.
		(*p).eat()
		(*p).returnForks()
	case <-waitingTime:
		(*p).think()
		(*p).ownFork.forkChannel <- true
	}
}

func (p *philosopher) returnForks() {
	//Make both forks ready, and signal fork routine. It is the go fork routine, that signals the philosophers back.
	(*p).ownFork.forkChannel <- true
	(*p).rightFork.forkChannel <- true
	(*p).ownFork.forkRunAgain <- true
	(*p).rightFork.forkRunAgain <- true
}

func (p *philosopher) eat() {
	//If we ate 3 times, just return.
	if (*p).doneEating {
		return
	}
	(*p).isThinking = false
	(*p).timesEaten++
	fmt.Println("philo: ", (*p).id, " is eating. Times eaten:", (*p).timesEaten)
	if (*p).timesEaten >= timesNeededToEat {
		(*p).doneEating = true
		fmt.Println("philo: ", (*p).id, " is full of food")
	}
	randomPause(eatingTime)
	fmt.Println("philo: ", (*p).id, " is done eating")
}

func (p *philosopher) think() {
	//Only print thinking once, untill next time we eat.
	if !(*p).isThinking {
		fmt.Println("philo: ", (*p).id, " is now thinking untill next bite.")
	}
	(*p).isThinking = true
}

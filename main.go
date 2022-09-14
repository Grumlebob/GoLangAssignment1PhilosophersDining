package main

import (
	"fmt"
	"math/rand"
	"time"
)

const numOfPhilosophers = 5
const timesNeededToEat = 100

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

/*
Why it doesn't deadlock:
Each channel only has a limit of 1 msg.
Each fork only sends 1 msg at max, meaning if that fork gets picked up,
it won't be able to send a message, aswell as being kept by
only a single philosopher, due to him also having only a single channel per philosopher with a limit of 1.
Due to passing forks by reference, we ensure we are talking to 1 channel per fork,
and same as philosophers.

The expected results are intertwined, aswell as only a max of 2 philosophers eating concurrently, due to
the limit of 5 forks.

The expected output is different each time, due to a random time delay, to simulate it taking time to pick up a fork
and even longer to eat.

We expect the program to terminate, when every philosopher has eaten %timesNeededToEat% times.
*/
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
	}

	//Repeat routine.
	for {
		//If everyone ate 3 times exit loop
		if philos[0].doneEating &&
			philos[1].doneEating &&
			philos[2].doneEating &&
			philos[3].doneEating &&
			philos[4].doneEating {
			break
		}
		//Start Fork and Philo routines. If fork is ready, signal philosopher.
		//Explictly not in loop, to make it visually easy to see that
		//every fork and philosopher has their own go routine.
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
			go philoPointer0.GoRoutinePhiloEatOrThink()
		case <-philos[1].philoChannel:
			philoPointer1 := &philos[1]
			go philoPointer1.GoRoutinePhiloEatOrThink()
		case <-philos[2].philoChannel:
			philoPointer2 := &philos[2]
			go philoPointer2.GoRoutinePhiloEatOrThink()
		case <-philos[3].philoChannel:
			philoPointer3 := &philos[3]
			go philoPointer3.GoRoutinePhiloEatOrThink()
		case <-philos[4].philoChannel:
			philoPointer4 := &philos[4]
			go philoPointer4.GoRoutinePhiloEatOrThink()
		}
	}
	fmt.Println("Everyone is full and ate:", timesNeededToEat, "times")
}

// Pause X * 0-40 milliseconds
func randomPause(max int) {
	time.Sleep(time.Millisecond * time.Duration(rand.Intn(max*40)))
}

func (p *philosopher) goRoutineForForks() {
	waitingTime := make(chan bool, 1)
	//Own go routine, so "waiting for waitingTime isn't itself blocked"
	go func() {
		randomPause(1)
		waitingTime <- true
	}()

	select {
	case <-(*p).ownFork.forkChannel: //Own fork is ready - announce to philosopher to attempt to eat.
		(*p).philoChannel <- (*p).id
	case <-waitingTime: //Own fork not ready, try again.
		(*p).ownFork.forkRunAgain <- true
	}

}

func (p *philosopher) GoRoutinePhiloEatOrThink() {
	waitingTime := make(chan bool, 1)
	//Own go routine, so "waiting for waitingTime isn't itself blocked"
	go func() {
		randomPause(3)
		waitingTime <- true
	}()
	//We got a channel msg that our own fork is ready, therefore check if other fork is ready, if so eat otherwise think.
	select {
	case <-(*p).rightFork.forkChannel: //We grab the other fork
		(*p).eat()
		(*p).returnForks()
	case <-waitingTime: //They put their own fork back.
		(*p).think()
		(*p).ownFork.forkChannel <- true
		(*p).ownFork.forkRunAgain <- true
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
	//If we ate 3 times, just return, otherwise print eat and status.
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
	//It takes some time to eat, wherein indirectly the keep both forks unavaible for other philosophers.
	randomPause(2)
	fmt.Println("philo: ", (*p).id, " is done eating")
}

func (p *philosopher) think() {
	//Only print thinking once, until next time we eat.
	if !(*p).isThinking {
		fmt.Println("philo: ", (*p).id, " is now thinking until next bite.")
	}
	(*p).isThinking = true
}

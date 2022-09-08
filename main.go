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
	inUse       bool
	id          int
	forkChannel chan bool //ready = true. Not ready = false /empty
}

func main() {
	//Create forks
	forks := make([]fork, numOfPhilosophers)
	for i := range forks {
		forks[i].id = i
		forks[i].inUse = false
		forks[i].forkChannel = make(chan (bool), 1)
		forks[i].forkChannel <- true
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
	//Signals all channels, that philos are ready to act
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
		/* OLD WORKING ROUTINE, WITH NO FORK ROUTINE.
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
		*/
		// NEW ROUTINE, WITH FORK ROUTINE.
		select {
		case <-philos[0].philoChannel:
			philoPointer0 := &philos[0]
			go philoPointer0.goRoutineForForks()
			go func() {
				time.Sleep(1111 * time.Millisecond)
				philoPointer0.attemptToGetForks()
			}()
		case <-philos[1].philoChannel:
			philoPointer1 := &philos[1]
			go philoPointer1.goRoutineForForks()
			go func() {
				time.Sleep(789 * time.Millisecond)
			}()
		case <-philos[2].philoChannel:
			philoPointer2 := &philos[2]
			go philoPointer2.goRoutineForForks()
			go func() {
				time.Sleep(653 * time.Millisecond)
			}()
		case <-philos[3].philoChannel:
			philoPointer3 := &philos[3]
			go philoPointer3.goRoutineForForks()
			go func() {
				time.Sleep(526 * time.Millisecond)
			}()
		case <-philos[4].philoChannel:
			philoPointer4 := &philos[4]
			go philoPointer4.goRoutineForForks()
			go func() {
				time.Sleep(367 * time.Millisecond)
			}()
		}

	}
	fmt.Println("Everyone is full and ate:", timesNeededToEat, " times")
}

func (p *philosopher) act(msg int) {
	(*p).ownFork.inUse = true
	//Eating
	if (*p).rightFork.inUse == false && (*p).timesEaten < timesNeededToEat {
		go func() {
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
		go func() {
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

func (p *philosopher) goRoutineForForks() {
	//KIGGER PÅ EGEN KNIV. NÅR DEN ER LEDIG, SIGNALER EGEN PHILOSOPHER.
	waitingTime := make(chan bool, 1)
	go func() {
		randomPause(1)
		waitingTime <- true
	}()

	select {
	case <-(*p).ownFork.forkChannel: //Own fork is ready - announce to philosopher to attempt to eat.
		(*p).ownFork.forkChannel <- true
		(*p).philoChannel <- (*p).id
	case <-waitingTime: //Own fork ikke ready - prøv igen.
		(*p).goRoutineForForks()
	}
}

func (p *philosopher) attemptToGetForks() {
	//KIGGER PÅ SIDEMANDS KNIV. HVIS LEDIG SÅ SPIS. ELLERS TÆNK OG LÆG EGEN KNIV LEDIG.
	<-(*p).ownFork.forkChannel
	waitingTime := make(chan bool, 1)
	go func() {
		randomPause(1)
		waitingTime <- true
	}()
	//Vi holder egen fork i op til 1 sek.

	//Hvis vi kan tage naboens, holder vi begge længere.
	//Hvis ikke så lægger vi den ned. Og begynder at tænke, mens andre kan snuppe.
	select {
	case <-(*p).rightFork.forkChannel: //Hvis vores nabo er ledige
		p.eat()
		p.returnForks() //p eller (*p) her?
	case <-waitingTime: //hvis nabo fork ikke er ledige, så tænker vi
		p.think()
		(*p).ownFork.forkChannel <- true
	}

}

func (p *philosopher) returnForks() {
	(*p).ownFork.forkChannel <- true
	(*p).rightFork.forkChannel <- true
	(*p).philoChannel <- (*p).id //Indikere philo kan tage et valg. Prøv at få ind i fork routine.
}

func (p *philosopher) eat() {
	(*p).isThinking = false
	(*p).timesEaten++
	fmt.Println("philo: ", (*p).id, " is eating. Times eaten:", (*p).timesEaten)
	if (*p).timesEaten == 3 {
		(*p).doneEating = true
		fmt.Println("philo: ", (*p).id, " is is full of food")
	}
	randomPause(eatingTime)
}

func (p *philosopher) think() {
	if !(*p).isThinking {
		fmt.Println("philo: ", (*p).id, " is thinking.")
	}
	(*p).isThinking = true
}

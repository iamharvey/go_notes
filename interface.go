package main

import "fmt"

type Robot interface {
	Forward(speed, dur int) int
	Backward(speed, dur int) int
	TurnLeft(degree int) int
	TurnRight(degree int) int
}

type Flybot interface {
	Forward(speed, dur int) int
	Backward(speed, dur int) int
	TurnLeft(degree int) int
	TurnRight(degree int) int
}

type MyRobot struct {
	Robot
}

func (bot *MyRobot) Forward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters forward ... \n", d)
	return d
}

func (bot *MyRobot) Backward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters backward ... \n", d)
	return d
}

func (bot *MyRobot) TurnLeft(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

func (bot *MyRobot) TurnRight(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

type MyRobot2 struct {}

func (bot *MyRobot2) Forward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters forward ... \n", d)
	return d
}

func (bot *MyRobot2) Backward(speed, dur int) (meter int) {
	d := speed * dur
	fmt.Printf("bot has moved %d meters backward ... \n", d)
	return d
}

func (bot *MyRobot2) TurnLeft(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

func (bot *MyRobot2) TurnRight(deg int) int {
	fmt.Printf("bot has turnd %d degree left ... \n", deg)
	return deg
}

func getBot() Robot {
	var b4 *MyRobot
	return b4
}

func main() {
	var bot MyRobot
	bot.Forward(1, 2)
	bot.TurnLeft(30)
	bot.Backward(1, 2)
	bot.TurnRight(45)

	fmt.Println("\nstarting the second robot ... ")
	var bot2 Robot = &MyRobot2{}
	bot2.Forward(1, 3)
	bot2.TurnLeft(25)
	bot2.Backward(2, 2)
	bot2.TurnRight(30)

	fmt.Println("\nstarting two agents, one on the ground, one in the sky ... ")
	var agent1 Robot = &MyRobot{}
	var agent2 Flybot = &MyRobot{}
	agent1.Forward(1, 2)
	agent1.TurnLeft(30)
	agent2.Backward(1, 2)
	agent2.TurnRight(45)

	//var bot3 Robot
	//fmt.Printf("%v, %T", bot3, bot3)
	//bot3.TurnRight(30)
	//the above code gives a runtime error

	fmt.Println()

	// empty interface
	var i interface{}
	fmt.Printf("%v, %T, %v\n", i, i, i == nil)

	i = 123
	_, ok := i.(int)
	fmt.Printf("%v, %T, %v\n", i, i, ok==true)

	i = "hello interface"
	_, ok = i.(string)
	fmt.Printf("%v, %T, %v\n", i, i, ok==true)

	// nil and empty
	fmt.Printf("%v, %T, %v", getBot(), getBot(), getBot() == nil)


}



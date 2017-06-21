package main

import (
	"fmt"
	"mylibs"
)

func HandlerTest(ptr *mylibs.LibClient, cmd int16, msg *mylibs.NetMsg) error {
	return nil
}

func main() {
	c := mylibs.NewLibClient(HandlerTest)
	err := c.Connect("127.0.0.1:9001")
	if err != nil {
		fmt.Println(err)
		return
	}
	c.Run()
	for {
		var msg string
		fmt.Print("input data:")
		fmt.Scanf("%s", &msg)

		dt := mylibs.NewNetMsg()
		dt.Cmd = 10
		dt.Data = []byte(msg)
		c.Send(dt)
	}
}

package main

import (
	"log"
	"mylibs"
	"time"
)

func HandlerTest1(cmd int16, msg *mylibs.NetMsg) error {
	log.Printf("cmd:%d,data:%s", cmd, string(msg.Data))
	return nil

}

func HandlerTest2(cmd int16, msg *mylibs.NetMsg) error {
	log.Printf("cmd:%d,data:%s", cmd, string(msg.Data))
	return nil
}

func main() {
	invoker := mylibs.NewLibInvokerHandler()
	invoker.Register(int16(10), HandlerTest1)
	invoker.Register(int16(11), HandlerTest2)

	code := mylibs.NewLibCodec()
	s := mylibs.NewLibServer(":9001", code, invoker)
	s.Run()

	time.Sleep(3)
}

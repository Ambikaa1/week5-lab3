package main

import (
	"flag"
	"fmt"
	"net"
	"net/rpc"
	"pairbroker/stubs"
)

type Factory struct{}

func (f *Factory) Multiply(pair stubs.Pair, response *stubs.JobReport) (err error) {
	response.Result = pair.X * pair.Y
	fmt.Println(response.Result)
	return
}

//TODO: Define a Multiply function to be accessed via RPC.
//Check the previous weeks' examples to figure out how to do this.

func main() {
	pAddr := flag.String("ip", "127.0.0.1:8050", "IP and port to listen on")            //8050 is how peeps will contact factory
	brokerAddr := flag.String("broker", "127.0.0.1:8030", "Address of broker instance") //address to send to
	flag.Parse()                                                                        //reads command lines
	rpc.Register(&Factory{})                                                            //added factory to my "phonebook"
	listener, _ := net.Listen("tcp", *pAddr)                                            //listens for tcp connection on 8050 from broker
	defer listener.Close()                                                              //closes the listener

	client, _ := rpc.Dial("tcp", *brokerAddr) //calling the broker back
	defer client.Close()
	//TODO: You'll need to set up the RPC server, and subscribe to the running broker instance.

	request := stubs.Subscription{Topic: "multiply", FactoryAddress: *pAddr, Callback: "Factory.Multiply"}
	response := new(stubs.StatusReport)
	client.Call(stubs.Subscribe, request, response)
	rpc.Accept(listener) //creates the connection
}

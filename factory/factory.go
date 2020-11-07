package main

import (
	"flag"
	"net"
	"net/rpc"
	"pairbroker/stubs"
)

type Factory struct{}

var number chan int

func divide(client *rpc.Client) {
	for {
		x := <-number
		y := <-number
		z := <-number
		//fmt.Println("hey")
		request := stubs.PublishRequest{Topic: "divide", Triplet: stubs.Triplet{X: x, Y: y, Z: z}}
		response := new(stubs.StatusReport)
		client.Call(stubs.Publish, request, response)
		//fmt.Println("it works")
	}
}

func (f *Factory) Multiply(triple stubs.Triplet, response *stubs.JobReport) (err error) { //hello
	response.Result = triple.X * triple.Y * triple.Z
	number <- response.Result
	//fmt.Println(response.Result)
	return
}

func (f *Factory) Divide(triplet stubs.Triplet, response *stubs.JobReport) (err error) {
	first := triplet.X / triplet.Y
	second := triplet.Y / triplet.Z
	response.Result = first / (second + 1)
	//fmt.Println(response.Result)
	return
}

//TODO: Define a Multiply function to be accessed via RPC.
//Check the previous weeks' examples to figure out how to do this.

func main() {
	number = make(chan int)
	pAddr := flag.String("ip", "127.0.0.1:8050", "IP and port to listen on")            //8050 is how peeps will contact factory
	brokerAddr := flag.String("broker", "127.0.0.1:8030", "Address of broker instance") //address to send to
	flag.Parse()                                                                        //reads command lines
	rpc.Register(&Factory{})                                                            //added factory to my "phonebook"
	listener, _ := net.Listen("tcp", *pAddr)                                            //listens for tcp connection on 8050 from broker
	defer listener.Close()                                                              //closes the listener

	client, _ := rpc.Dial("tcp", *brokerAddr) //calling the broker back
	defer client.Close()
	//TODO: You'll need to set up the RPC server, and subscribe to the running broker instance.

	requestM := stubs.Subscription{Topic: "multiply", FactoryAddress: *pAddr, Callback: "Factory.Multiply"}
	responseM := new(stubs.StatusReport)
	client.Call(stubs.Subscribe, requestM, responseM)

	chanreq := stubs.ChannelRequest{Topic: "divide", Buffer: 10}
	res := new(stubs.StatusReport)
	client.Call(stubs.CreateChannel, chanreq, res)

	requestD := stubs.Subscription{Topic: "divide", FactoryAddress: *pAddr, Callback: "Factory.Divide"}
	responseD := new(stubs.StatusReport)
	client.Call(stubs.Subscribe, requestD, responseD)

	go divide(client)

	rpc.Accept(listener) //creates the connection
}

package stubs

var CreateChannel = "Broker.CreateChannel"
var Publish = "Broker.Publish"
var Subscribe = "Broker.Subscribe"
var Multiply = "Factory.Multiply"

type Pair struct {
	X int
	Y int
}

type Triplet struct {
	X int
	Y int
	Z int
}

type PublishRequest struct {
	Topic   string
	Triplet Triplet
}

type ChannelRequest struct {
	Topic  string
	Buffer int
}

type Subscription struct {
	Topic          string
	FactoryAddress string
	Callback       string
}

type JobReport struct {
	Result int
}

type StatusReport struct {
	Message string
}

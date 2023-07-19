package iface

type Iproducer interface {
	PushMessageToQueue(topic string, message []byte) error
}

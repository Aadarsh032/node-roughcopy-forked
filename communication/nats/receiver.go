package nats

import (
	"github.com/mysterium/node/communication"
	"github.com/nats-io/go-nats"

	"fmt"
	log "github.com/cihub/seelog"
	"github.com/mysterium/node/communication/nats_discovery"
)

const RECEIVER_LOG_PREFIX = "[NATS.Receiver] "

func NewReceiver(address *nats_discovery.NatsAddress) *receiverNats {
	return &receiverNats{
		connection:   address.GetConnection(),
		codec:        communication.NewCodecJSON(),
		messageTopic: address.GetTopic() + ".",
	}
}

type receiverNats struct {
	connection   *nats.Conn
	codec        communication.Codec
	messageTopic string
}

func (receiver *receiverNats) Receive(handler communication.MessageHandler) error {

	messageType := string(handler.GetMessageType())

	messageHandler := func(msg *nats.Msg) {
		log.Debug(RECEIVER_LOG_PREFIX, fmt.Sprintf("Message '%s' received: %s", messageType, msg.Data))
		messagePtr := handler.NewMessage()
		err := receiver.codec.Unpack(msg.Data, messagePtr)
		if err != nil {
			err = fmt.Errorf("Failed to unpack message '%s'. %s", messageType, err)
			log.Error(RECEIVER_LOG_PREFIX, err)
			return
		}

		err = handler.Handle(messagePtr)
		if err != nil {
			err = fmt.Errorf("Failed to process message '%s'. %s", messageType, err)
			log.Error(RECEIVER_LOG_PREFIX, err)
			return
		}
	}

	_, err := receiver.connection.Subscribe(receiver.messageTopic+messageType, messageHandler)
	if err != nil {
		err = fmt.Errorf("Failed subscribe message '%s'. %s", messageType, err)
		return err
	}

	return nil
}

func (receiver *receiverNats) Respond(consumer communication.RequestConsumer) error {

	requestType := string(consumer.GetRequestType())

	messageHandler := func(msg *nats.Msg) {
		log.Debug(RECEIVER_LOG_PREFIX, fmt.Sprintf("Request '%s' received: %s", requestType, msg.Data))
		requestPtr := consumer.NewRequest()
		err := receiver.codec.Unpack(msg.Data, requestPtr)
		if err != nil {
			err = fmt.Errorf("Failed to unpack request '%s'. %s", requestType, err)
			log.Error(RECEIVER_LOG_PREFIX, err)
			return
		}

		response, err := consumer.Consume(requestPtr)
		if err != nil {
			err = fmt.Errorf("Failed to process request '%s'. %s", requestType, err)
			log.Error(RECEIVER_LOG_PREFIX, err)
			return
		}

		responseData, err := receiver.codec.Pack(response)
		if err != nil {
			err = fmt.Errorf("Failed to pack response '%s'. %s", requestType, err)
			log.Error(RECEIVER_LOG_PREFIX, err)
			return
		}

		err = receiver.connection.Publish(msg.Reply, responseData)
		if err != nil {
			err = fmt.Errorf("Failed to send response '%s'. %s", requestType, err)
			log.Error(RECEIVER_LOG_PREFIX, err)
			return
		}

		log.Debug(RECEIVER_LOG_PREFIX, fmt.Sprintf("Request '%s' response: %s", requestType, responseData))
	}

	_, err := receiver.connection.Subscribe(receiver.messageTopic+requestType, messageHandler)
	if err != nil {
		err = fmt.Errorf("Failed subscribe request '%s'. %s", requestType, err)
		return err
	}

	return nil
}

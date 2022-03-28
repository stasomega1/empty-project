package store

import (
	"encoding/json"
	"fmt"
	"project/inetrnal/app/model"

	"github.com/wagslane/go-rabbitmq"
)

type RabbitPublisher struct {
	*rabbitmq.Publisher
}

func (r *RabbitPublisher) SendSomeDbMessage(model model.DbModel) error {
	jsonByteData, err := json.Marshal(model)
	if err != nil {
		return fmt.Errorf("SendSomeDbMessage: %v", err)
	}

	err = r.Publish(jsonByteData, nil)
	if err != nil {
		return fmt.Errorf("SendSomeDbMessage: %v", err)
	}

	return nil
}

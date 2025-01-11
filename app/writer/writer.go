package writer

import (
	"encoding/json"
	"errors"
	"github.com/IBM/sarama"
	"golang-exercise/models"
)

type Writer struct {
	brokers          []string
	destinationTopic string
	producer         sarama.SyncProducer
}

var ErrBrokersNotProvided = errors.New("brokers not provided")

func NewWriter(opts ...func(*Writer)) (*Writer, error) {
	w := new(Writer)
	for _, o := range opts {
		o(w)
	}

	if len(w.brokers) == 0 {
		return nil, ErrBrokersNotProvided
	}

	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	var err error
	w.producer, err = sarama.NewSyncProducer(w.brokers, c)
	if err != nil {
		return nil, err
	}
	return w, nil
}

func (w *Writer) Write(event *models.Event) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return err
	}
	message := &sarama.ProducerMessage{
		Topic: w.destinationTopic,
		Key:   sarama.ByteEncoder(event.CompanyID),
		Value: sarama.ByteEncoder(eventBytes),
	}

	_, _, err = w.producer.SendMessage(message)
	if err != nil {
		return err
	}
	return nil
}

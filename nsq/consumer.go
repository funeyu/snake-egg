package nsq

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"snake/indexer"
)


type NewDocConsumer struct {
	c *nsq.Consumer
}

func CreateDocConsumer(cha string) (*NewDocConsumer, error) {
	config := nsq.NewConfig()
	q, _ := nsq.NewConsumer("NewDoc", cha, config)
	return &NewDocConsumer{c:q}, nil
}
type ListenNewDoc func(d indexer.Doc) error

func (nb *NewDocConsumer)Listen(l ListenNewDoc) {
	nb.c.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error{
		var d indexer.Doc
		err := json.Unmarshal(message.Body, &d)
		if err != nil {
			return err
		}
		e := l(d)
		if e != nil {
			return err
		}
		return nil
	}))
}

func (nb *NewDocConsumer) Init() error{
	err := nb.c.ConnectToNSQD("192.168.1.13:4150")
	if err != nil {
		return err
	}
	return nil
}

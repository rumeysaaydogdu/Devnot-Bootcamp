package example

import (
	"encoding/json"
	"fmt"
	rabbit "github.com/alperhankendi/golang-api/pkg/rabbitmq"

	"time"
)

type (
	Person struct {
		Name    string
		Surname string
		City    City
	}

	City struct {
		Name string
	}
)

func main() {

	var rabbitClient = rabbit.NewRabbitMqClient([]string{"127.0.0.1"}, "guest", "guest", "", rabbit.RetryCount(2), rabbit.PrefetchCount(3))

	onConsumed := func(message rabbit.Message) error {

		var consumeMessage Person
		var err = json.Unmarshal(message.Payload, &consumeMessage)
		if err != nil {
			return err
		}
		fmt.Println(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), " Message:", consumeMessage)
		return nil
	}

	onConsumed2 := func(message rabbit.Message) error {

		var consumeMessage Person
		var err = json.Unmarshal(message.Payload, &consumeMessage)
		if err != nil {
			return err
		}
		fmt.Println(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), " Message:", consumeMessage)
		time.Sleep(10000000)
		return nil
	}

	onConsumed3 := func(message rabbit.Message) error {

		var consumeMessage Person
		var err = json.Unmarshal(message.Payload, &consumeMessage)
		if err != nil {
			return err
		}
		fmt.Println(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), " Message:", consumeMessage)

		return nil
	}

	rabbitClient.AddConsumer("In.Person_Direct").SubscriberExchange("person", rabbit.Direct, "PersonV1_Direct").HandleConsumer(onConsumed)
	rabbitClient.AddConsumer("In.Person_Fanout").SubscriberExchange("", rabbit.Fanout, "PersonV3_Fanout").HandleConsumer(onConsumed2)
	rabbitClient.AddConsumer("In.Person_Topic").SubscriberExchange("person", rabbit.Topic, "PersonV4_Topic").HandleConsumer(onConsumed3)

	rabbitClient.RunConsumers()

}

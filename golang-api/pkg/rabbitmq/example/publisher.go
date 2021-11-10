package example

//
//import (
//	"fmt"
//	rabbit "github.com/alperhankendi/golang-api/pkg/rabbitmq"
//	"golang.org/x/net/context"
//)
//
//func main() {
//
//	var rabbitClient = rabbit.NewRabbitMqClient([]string{"127.0.0.1"}, "guest", "guest", "", rabbit.RetryCount(2))
//
//	rabbitClient.AddPublisher("PersonV1_Direct", rabbit.Direct, Person{})
//
//	var err = rabbitClient.Publish(context.TODO(), "123", Person{
//		Name:    "John",
//		Surname: "Jack",
//		City:    City{},
//	})
//	fmt.Print(err)
//}

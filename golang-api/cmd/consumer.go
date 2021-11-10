package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/alperhankendi/golang-api/internal/basket"
	"github.com/alperhankendi/golang-api/internal/config"
	rabbit "github.com/alperhankendi/golang-api/pkg/rabbitmq"
	"github.com/spf13/cobra"
	"runtime"
	"time"
)

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "A brief description of your command",
}

func init() {
	rootCmd.AddCommand(consumerCmd)
	var cfgFile string
	consumerCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.dev.yaml", "config file (default is $HOME/.golang-api.yaml)")
	ApiConfig, err := config.GetAllValues("./config/", cfgFile)
	if err != nil {
		panic(err)
	}
	consumerCmd.Run = func(cmd *cobra.Command, args []string) {

		var rabbitClient = rabbit.NewRabbitMqClient([]string{ApiConfig.RabbitMQSettings.Url}, ApiConfig.RabbitMQSettings.Username, ApiConfig.RabbitMQSettings.Password, "", rabbit.RetryCount(2))

		rabbitClient.AddConsumer("BasketCreated").
			SubscriberExchange("", rabbit.Direct, "Basket_Direct").
			HandleConsumer(ReadMessageFromRabbitMq)
		rabbitClient.RunConsumers()

		runtime.GC()

	}
}
func ReadMessageFromRabbitMq(message rabbit.Message) error {

	var consumeMessage basket.Basket
	var err = json.Unmarshal(message.Payload, &consumeMessage)
	if err != nil {
		return err
	}
	fmt.Println(time.Now().Format("Mon, 02 Jan 2006 15:04:05 "), " Message:", consumeMessage)
	//DO Something
	return nil
}

package cmd

import (
	"context"
	"fmt"
	"github.com/alperhankendi/golang-api/internal/basket"
	"github.com/alperhankendi/golang-api/internal/config"
	mongodata "github.com/alperhankendi/golang-api/internal/mongo"
	echoextentions "github.com/alperhankendi/golang-api/pkg/eechoextentions"
	"github.com/alperhankendi/golang-api/pkg/log"
	"github.com/alperhankendi/golang-api/pkg/mongoHelper"
	rabbit "github.com/alperhankendi/golang-api/pkg/rabbitmq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"time"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "A brief description of your application",
}

func init() {
	rootCmd.AddCommand(apiCmd)
	var cfgFile string
	var port string
	apiCmd.PersistentFlags().StringVarP(&port, "port", "p", "5000", "Restfull Service Port")
	apiCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.dev.yaml", "config file (default is $HOME/.golang-api.yaml)")
	ApiConfig, err := config.GetAllValues("./config/", cfgFile)
	if err != nil {
		panic(err)
	}
	apiCmd.Run = func(cmd *cobra.Command, args []string) {
		//application bootstrapper
		instance := echo.New()
		instance.Logger = log.SetupLogger()
		//custom middware , manipulate your response header
		instance.Use(ServerHeaderFunc)
		//Limit your payload size
		instance.Use(middleware.BodyLimit("1M"))
		//enable rate limiter
		instance.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(50)))
		//Gate Logger
		instance.Use(echoextentions.HookGateLoggerWithConfig(echoextentions.GateLoggerConfig{
			IncludeRequestBodies:  true,
			IncludeResponseBodies: true,
			Skipper:               echoextentions.Myskipper,
		}))

		db, err := mongoHelper.ConnectDb(ApiConfig.MongoSettings)
		if err != nil {
			log.Logger.Fatalf("Database connection problem,Error: %v", err)

		}
		var rabbitClient = rabbit.NewRabbitMqClient([]string{ApiConfig.RabbitMQSettings.Url}, ApiConfig.RabbitMQSettings.Username, ApiConfig.RabbitMQSettings.Password, "", rabbit.RetryCount(2))
		rabbitClient.AddPublisher("Basket_Direct", rabbit.Direct, basket.Basket{})
		rabbitClient.Publish(context.TODO(), "", basket.Basket{Id: "abc"})

		repository := mongodata.NewRepository(db)
		domainService := basket.NewService(repository, rabbitClient)
		resource := basket.NewResource(domainService)
		basket.RegisterHandlers(instance, resource)
		//Register Handlers ==> resource => domain service ==> reposiyory => mongo database
		log.Logger.Info("Service is starting, will serve on", port, " port")

		if err := instance.Start(fmt.Sprintf(":%s", port)); err != nil {
			log.Logger.Fatalf("shutting down the server: Error:", err.Error())
		}
		echoextentions.Shutdown(instance, time.Second*2)
	}

}

func ServerHeaderFunc(handlerFunc echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("ServerHeader", "TestServer/1.0")
		return handlerFunc(c)
	}
}

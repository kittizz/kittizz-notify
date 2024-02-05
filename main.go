package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/spf13/viper"
)

var notifys = make(map[string]*notify.Notify)

func init() {

	godotenv.Load()
	viper.AutomaticEnv()

	//env: NT_SYNOLOGY=token|chatid
	// fmt.Println(viper.GetString("NT_SYNOLOGY"))
	for _, v := range os.Environ() {
		variable := strings.Split(v, "=")
		envKey := variable[0]
		if strings.HasPrefix(strings.ToLower(envKey), "nt_") {
			key := strings.Replace(strings.ToLower(envKey), "nt_", "", -1)
			cfg := strings.Split(viper.GetString(envKey), "|")
			if len(cfg) < 2 {
				log.Fatalf("failed to parse %s=%s", envKey, v)
				continue
			}
			token := cfg[0]
			receiver, err := strconv.ParseInt(cfg[1], 10, 64)
			if err != nil {
				log.Fatalf("strconv error: %v", err)
				continue
			}

			telegramService, _ := telegram.New(token)
			telegramService.AddReceivers(receiver)
			n := notify.New()
			n.UseServices(telegramService)
			notifys[key] = n
			log.Printf("notify service %s is loaded", key)
		}
	}
	if len(notifys) == 0 {
		log.Panic("notify service is not configured")
	}
}

type Msg struct {
	Message string `json:"message" xml:"message" form:"message" query:"message"`
}

func main() {
	app := fiber.New()
	api := app.Group("/:domain", withNotify())
	{
		api.All("", func(c fiber.Ctx) error {
			noti := c.Locals("notify").(*notify.Notify)
			var msg Msg
			if err := c.Bind().Body(&msg); err != nil {
				if err := c.Bind().Query(&msg); err != nil {
					return err
				}
			}
			if err := noti.Send(c.Context(), c.Locals("domain").(string), msg.Message); err != nil {
				return err
			}

			return c.SendString("ok")
		})
	}

	log.Fatal(app.Listen(":" + os.Getenv("PORT")))
}

func withNotify() fiber.Handler {
	return func(c fiber.Ctx) error {
		domain := strings.ToLower(c.Params("domain"))
		c.Locals("notify", notifys[domain])
		c.Locals("domain", domain)
		return c.Next()
	}
}

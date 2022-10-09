package webapi

type (
	WebApi struct {
		TelegramWebApi *TelegramWebApi
		YandexWebApi   *YandexWebApi
	}

	Config struct {
		TG TG
		Ya Ya
	}
	TG struct {
		Webhook  string
		Endpoint string
		Token    string
	}

	Ya struct {
		ClientID     string
		ClientSecret string
		Host         string
	}
)

func New(conf Config) *WebApi {
	return &WebApi{
		TelegramWebApi: NewTelegramWebApi(conf.TG),
		YandexWebApi:   NewYandexWebApi(conf.Ya),
	}
}

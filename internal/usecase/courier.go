package usecase

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/OksidGen/market-couriers_salary/internal/entity"
)

type CourierUseCase struct {
	UserRepo   UserRepo
	TokenRepo  TokenRepo
	IncomeRepo IncomeRepo

	TelegramWebApi TelegramWebApi
	YandexWebApi   YandexWebApi
}

func NewCourierUseCase(userRepo UserRepo, tokenRepo TokenRepo, incomeRepo IncomeRepo, tgWebapi TelegramWebApi, yaWebapi YandexWebApi) *CourierUseCase {
	return &CourierUseCase{
		UserRepo:       userRepo,
		TokenRepo:      tokenRepo,
		IncomeRepo:     incomeRepo,
		TelegramWebApi: tgWebapi,
		YandexWebApi:   yaWebapi,
	}
}

func (uc *CourierUseCase) GetToken(code string, tgId int64) error {
	token := entity.Token{
		TGID: tgId,
	}
	if err := json.Unmarshal(uc.YandexWebApi.SendCodeForToken(code), &token); err != nil {
		return fmt.Errorf("error decode map.Token to entity.Token: %w", err)
	}
	if err := uc.TokenRepo.Create(&token); err != nil {
		return err
	}
	return nil
}

func (uc *CourierUseCase) TelegramParser(req *http.Request) map[string]interface{} {
	upd, err := uc.TelegramWebApi.ParseRequest(req)
	if err != nil {
		log.Printf("Error with parsing update request from TelegramWebApi: %v", err)
		return nil
	}
	return upd
}

func (uc *CourierUseCase) SendMessage(msg map[string]interface{}) {
	uc.TelegramWebApi.SendMessage(msg)
}

var (
	mainKeyboard = map[string]interface{}{
		"keyboard": [][]string{
			{"Расчет смены"},
			{"Прошлая неделя", "Текущая неделя"},
			{"Помощь"},
		},
	}
	startKeyboard = map[string]interface{}{
		"keyboard": [][]string{
			{"Авторизация"},
			{"Помощь"},
		},
	}
	authKeyboard = map[string]interface{}{
		"keyboard": [][]string{
			{"Проверка"},
			{"Помощь"},
		},
	}

	errorMessage = "Внутренняя ошибка, пожалуйста сообщите мне о ней, используя раздел `Помощь`"
	authMessage  = "Для авторизации в сервисе Вам необходимо перейти по [этой ссылке](%v) и войти, используя свой рабочий аккаунт Яндекс\\. Это необходимо, чтобы программа могла узнавать количество выполненных Вами доставок для корректного расчета дохода\\."
	waitMessage  = "Бот еще не получил информацию о Ваших доставках от Яндекс, возможно Вы забыли перейти по ссылке и авторизоваться или допустили ошибку при авторизации\\.\n\nПопробуйте еще раз \\- [ссылка](%v)"
	startMessage = "Добро пожаловать!\n\nЭтот сервис помогает курьерам Яндекс.Маркет самостоятельно узнавать свою заработную плату за день и прошлую/текущую недели.\n\nФормулы для расчета ориентированы на бо́льшую часть курьеров, то есть 75%, не вошедших в `топ25` своего склада. Так же в них не учитывается километраж (возможно функция будет добавлена позже), но пока что используются значения для минимального пробега(0-100).\n\nПрограмма может иметь некоторые ошибки (в том числе в формуле для расчета) если Вы столкнулись с ними или у Вас есть какие-либо предложения, пожалуйста, свяжитесь со мной, контакт Вы всегда можете найти в разделе `Помощь`, но если и этот раздел у Вас не работает пишите мне - @nightmarezero\n\nПриятного использования! Надеюсь этот сервис принесет Вам пользу!"
	okMessage    = "Авторизация успешно пройдена!"
	helpMessage  = "Если у Вас возникла какая-то проблема с сервисом или есть предложение по его улучшению, то напишите мне - @nightmarezero"
)

func (uc *CourierUseCase) CaseStart(msg map[string]interface{}) {
	token, err := uc.TokenRepo.Find(msg["chat_id"].(int64))
	if err != nil {
		msg["text"] = errorMessage
		log.Println(err)
		return
	}
	if token != "" {
		msg["text"] = "Главное меню"
		msg["reply_markup"] = mainKeyboard
		return
	}
	userCheck, err := uc.UserRepo.Check(msg["chat_id"].(int64))
	if err != nil {
		msg["text"] = errorMessage
		log.Println(err)
		return
	}
	if userCheck {
		msg["text"] = fmt.Sprintf(waitMessage, uc.YandexWebApi.GetURLForYandexAuth(msg["chat_id"].(int64)))
		msg["parse_mode"] = "MarkdownV2"
		msg["reply_markup"] = authKeyboard
		return
	}
	msg["text"] = startMessage
	msg["reply_markup"] = startKeyboard
}

func (uc *CourierUseCase) CaseAuth(msg map[string]interface{}) {
	user := entity.User{
		TGID:      msg["chat_id"].(int64),
		FirstName: msg["FirstName"].(string),
		LastName:  msg["LastName"].(string),
		UserName:  msg["UserName"].(string),
	}
	if err := uc.UserRepo.Create(&user); err != nil {
		msg["text"] = errorMessage
		log.Println(err)
		return
	}
	msg["text"] = fmt.Sprintf(authMessage, uc.YandexWebApi.GetURLForYandexAuth(msg["chat_id"].(int64)))
	msg["parse_mode"] = "MarkdownV2"
	msg["reply_markup"] = authKeyboard
}

func (uc *CourierUseCase) CaseCheck(msg map[string]interface{}) {
	res, err := uc.TokenRepo.Find(msg["chat_id"].(int64))
	if err != nil {
		log.Println(err)
		return
	}
	if res != "" {
		msg["text"] = okMessage
		msg["reply_markup"] = mainKeyboard
		return
	}
	msg["text"] = fmt.Sprintf(waitMessage, uc.YandexWebApi.GetURLForYandexAuth(msg["chat_id"].(int64)))
	msg["parse_mode"] = "MarkdownV2"
}

func (uc *CourierUseCase) CaseCalculate(msg map[string]interface{}) {
	token, err := uc.TokenRepo.Find(msg["chat_id"].(int64))
	if err != nil {
		log.Println(err)
		return
	}
	if token == "" {
		msg["text"] = errorMessage
		return
	}
	switch uc.YandexWebApi.CheckCompletedTask(token) {
	case "more":
		msg["text"] = "Еще выполнены не все доставки! Попробуйте еще раз в конце рабочего дня"
		return
	case "empty":
		msg["text"] = "Нет доставок на сегодня"
		return
	}
	res := uc.YandexWebApi.CalculateIncome(token)
	income := entity.Income{
		TGID:         msg["chat_id"].(int64),
		Count:        res["count"],
		Client:       res["client"],
		PvzPstLocker: res["pvz_pst_locker"],
		Multi5:       res["multi5"],
		Income:       res["income"],
	}
	if err := uc.IncomeRepo.Create(&income); err != nil {
		log.Println(err)
		msg["text"] = errorMessage
		return
	}
	msg["text"] = income.ToTableStirng()
	msg["parse_mode"] = "MarkdownV2"
}

func (uc *CourierUseCase) CaseHelp(msg map[string]interface{}) {
	msg["text"] = helpMessage
}

func (uc *CourierUseCase) CaseWeek(msg map[string]interface{}, week string) {
	incomes := uc.IncomeRepo.Select(msg["chat_id"].(int64), week)
	switch week {
	case "last":
		week = "прошлую"
	case "current":
		week = "текущую"
	}
	out := fmt.Sprintf("Статистика дохода за %v неделю:", week)
	totalInc := 0
	for _, inc := range incomes {
		out += fmt.Sprintf("\n\n```%v```\n", inc.ToTableStirng())
		totalInc += inc.Income
	}
	out += fmt.Sprintf("\nИтого общий доход : %v ₽", totalInc)
	msg["text"] = out
	msg["parse_mode"] = "MarkdownV2"
}

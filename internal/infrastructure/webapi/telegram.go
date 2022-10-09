package webapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	Update struct {
		UpdateId int `json:"update_id"`
		Message  struct {
			Text string `json:"text"`
			From struct {
				Id        int64  `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				UserName  string `json:"username"`
			} `json:"from"`
		} `json:"message"`
	}

	Message struct {
		Chat_id     int64      `json:"chat_id"`
		Text        string     `json:"text"`
		ParseMode   string     `json:"parse_mode,omitempty"`
		ReplyMarkup [][]string `json:"reply_markup,omitempty"`
	}

	APIResource struct {
		OK          bool            `json:"ok"`
		Result      json.RawMessage `json:"result"`
		ErrorCode   int             `json:"error_code"`
		Description string          `json:"description"`
	}

	TelegramWebApi struct {
		token   string
		client  http.Client
		host    string
		webhook string
	}
)

func NewTelegramWebApi(conf TG) *TelegramWebApi {
	TelegramWebApi := &TelegramWebApi{
		token:   conf.Token,
		client:  http.Client{},
		host:    fmt.Sprintf("%v/bot%v", conf.Endpoint, conf.Token),
		webhook: conf.Webhook,
	}

	if err := TelegramWebApi.getMe(); err != nil {
		log.Fatalf("Telegram bot - GetMe - Error: %v", err)
	}
	if err := TelegramWebApi.setWebhook(); err != nil {
		log.Fatalf("Telegram bot - setWebhook - Error: %v", err)
	}

	return TelegramWebApi
}

func (api *TelegramWebApi) ParseRequest(req *http.Request) (map[string]interface{}, error) {
	var upd Update
	if err := decodeBody(req.Body, &upd); err != nil {
		return nil, fmt.Errorf("error decoding body - %w", err)
	}
	out := map[string]interface{}{
		"FromId":    upd.Message.From.Id,
		"FirstName": upd.Message.From.FirstName,
		"LastName":  upd.Message.From.LastName,
		"UserName":  upd.Message.From.UserName,
		"Text":      upd.Message.Text,
	}
	return out, nil
}

func (api *TelegramWebApi) SendMessage(input map[string]interface{}) {
	if err := api.sendRequest("sendMessage", input); err != nil {
		log.Printf("Error sending message: %v\n", err)
	}
}

func (api *TelegramWebApi) getMe() error {
	if err := api.sendRequest("getMe", nil); err != nil {
		return err
	}
	return nil
}

func (api *TelegramWebApi) setWebhook() error {
	param := map[string]interface{}{
		"url": "",
	}
	if err := api.sendRequest("setWebhook", param); err != nil {
		return err
	}
	param["url"] = api.webhook
	if err := api.sendRequest("setWebhook", param); err != nil {
		return err
	}
	return nil
}

func (api *TelegramWebApi) sendRequest(endpoint string, params map[string]interface{}) error {
	url := fmt.Sprintf("%v/%v", api.host, endpoint)
	bodyJSON, err := buildJSONBody(params)
	if err != nil {
		return fmt.Errorf("error building JSON body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyJSON))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := api.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var apiResp APIResource
	if err := decodeBody(resp.Body, &apiResp); err != nil {
		return err
	}
	if !apiResp.OK {
		return fmt.Errorf("error: %d - desc: %s", apiResp.ErrorCode, apiResp.Description)
	}
	return nil
}

func buildJSONBody(params map[string]interface{}) ([]byte, error) {
	if params == nil {
		return nil, nil
	}
	out, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("error marshal params: %w", err)
	}
	return out, nil
}

func decodeBody(body io.Reader, out any) error {
	data, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &out)
	if err != nil {
		return err
	}
	return nil
}

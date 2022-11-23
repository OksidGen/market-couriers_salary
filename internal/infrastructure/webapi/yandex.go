package webapi

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"strconv"

	"golang.org/x/exp/slices"
)

type (
	task struct {
		RoutePointID int `json:"routePointId"`
		// TaskID int `json:"taskId"`
		// Type          string `json:"type"`
		OrderType  string `json:"orderType"`
		TaskStatus string `json:"taskStatus"`
		PlaceCount int    `json:"placeCount"`
		// OrdinalNumber int    `json:"ordinalNumber"`
	}

	taskMassive struct {
		Tasks []task `json:"tasks"`
	}

	summary struct {
		UnfinishedOrderCount int `json:"unfinishedOrderCount"`
		TotalOrderCount      int `json:"totalOrderCount"`
	}
	countOfOrder struct {
		Summary summary `json:"summary"`
	}

	YandexWebApi struct {
		clientID     string
		clientSecret string
		client       *http.Client
		host         string
	}
)

const (
	priceTask     = 90
	priceMulti    = 60
	pricePackage  = 5
	minimumIncome = 3500

	yandexAuth = "https://oauth.yandex.ru"
)

func NewYandexWebApi(conf Ya) *YandexWebApi {
	return &YandexWebApi{
		clientID:     conf.ClientID,
		clientSecret: conf.ClientSecret,
		client:       &http.Client{},
		host:         conf.Host,
	}
}

func (api *YandexWebApi) GetURLForYandexAuth(id int64) string {
	host := "https://oauth.yandex.ru"
	end_point := "/authorize"
	resp_type := "response_type=code"
	client_id := "client_id=" + api.clientID
	tgid := "state=" + strconv.FormatInt(id, 10)
	url := host + end_point + "?" + resp_type + "&" + client_id + "&" + tgid
	return url
}

func (api *YandexWebApi) CheckCompletedTask(token string) string {
	coo := countOfOrder{}
	getFromBody(api.sendGETRequest("/api/tasks/order-delivery", token), &coo)
	log.Println(coo)
	if coo.Summary.TotalOrderCount == 0 {
		return "empty"
	}
	if coo.Summary.UnfinishedOrderCount == 0 {
		return "end"
	}
	return "more"
}

func (api *YandexWebApi) CalculateIncome(token string) map[string]int {
	tm := taskMassive{}
	getFromBody(api.sendGETRequest("/api/tasks", token), &tm)
	tm.brushMassive()

	out := map[string]int{
		"pvz_pst_locker": 0,
		"count":          0,
		"client":         0,
		"multi5":         0,
		"income":         0,
	}

	for _, task := range tm.Tasks {
		switch task.OrderType {
		case "CLIENT":
			if task.PlaceCount > 4 {
				out["multi5"] += task.PlaceCount
				out["income"] += task.PlaceCount * priceMulti
			}
			out["client"]++
		default:
			out["pvz_pst_locker"] += task.PlaceCount
			out["income"] += task.PlaceCount * pricePackage
			out["count"]++
		}
		out["income"] += priceTask
	}
	if out["income"] < minimumIncome {
		out["income"] = minimumIncome
	}
	return out
}

func (api *YandexWebApi) SendCodeForToken(code string) []byte {
	return api.sendPOSTRequest("/token", code)
}

func (api *YandexWebApi) sendPOSTRequest(endpoint, code string) []byte {
	URL := fmt.Sprint(yandexAuth, endpoint)
	v := url.Values{
		"grant_type": {"authorization_code"},
		"code":       {code},
	}
	req, err := http.NewRequest(http.MethodPost, URL, strings.NewReader(v.Encode()))
	if err != nil {
		log.Printf("error create post request: %v\n", err)
		return nil
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+api.basicAuth())
	resp, err := api.client.Do(req)
	if err != nil {
		log.Printf("error send post request: %v\n", err)
		return nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error read response body: %v\n", err)
		return nil
	}
	return body
}

func (api *YandexWebApi) sendGETRequest(endpoint, token string) []byte {
	url := fmt.Sprint(api.host, endpoint)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("error create get request: %v\n", err)
		return nil
	}
	req.Header.Add("Authorization", "OAuth "+token)
	resp, err := api.client.Do(req)
	if err != nil {
		log.Printf("error send get request: %v\n", err)
		return nil
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error read response body: %v\n", err)
		return nil
	}
	return body
}

func (api *YandexWebApi) basicAuth() string {
	auth := fmt.Sprintf("%v:%v", api.clientID, api.clientSecret)
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func getFromBody(body []byte, res any) {
	err := json.Unmarshal(body, &res)
	if err != nil {
		log.Printf("error unmarshalling body: %v\n", err)
	}
}

func (tm *taskMassive) brushMassive() {
	var temp []task
	for _, oneTask := range tm.Tasks {
		if oneTask.TaskStatus == "DELIVERED" {
			i := slices.IndexFunc(temp, func(t task) bool { return t.RoutePointID == oneTask.RoutePointID })
			if i != -1 {
				temp[i].PlaceCount += oneTask.PlaceCount
			} else {
				temp = append(temp, oneTask)
			}
		}
	}
	tm.Tasks = temp
}

package main

import (
	"encoding/json"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
)

var (
	r         = gin.Default()
	securekey = "/key-1212"
	mcache    = cache.New(cache.NoExpiration, cache.NoExpiration)
)

func main() {

	getRoutes() //издеся маршруты роутим

	r.Run()
}

func getRoutes() {

	r.GET(securekey+"/debug/:debug", debugGet)
	r.POST(securekey+"/debug", debugPost)

	r.GET(securekey+"/gamestart/:id", gamestart)
	r.GET(securekey+"/holeCards/:id/:cards", holeCards)
	r.POST(securekey+"/dealHoleCardsEvent", dealHoleCardsEventPost)
	r.GET(securekey+"/stageEvent/:id/:stage/:board", stageEvent)
	r.POST(securekey+"/getAction", getActionPost)
	r.POST(securekey+"/actionEvent", actionEventPost)
	r.POST(securekey+"/showdownEvent", showdownEventPost)
	r.POST(securekey+"/winEvent", winEventPost)
	r.GET(securekey+"/gameOverEvent/:id", gameOverEvent)

}

func debugGet(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	debug := c.Param("debug")
	c.Status(200)
	println(botname + " : " + debug)
}
func debugPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	println(botname + " : " + c.PostForm("DEBUG"))

}

//началась новая раздача, все боты присылают это сообщение. Но кто пришлет первым - будет главным и в текущей раздаче будет присылать основную инфу
func gamestart(c *gin.Context) {
	//botname := c.Request.Header.Get("User-Agent")//botname=undefined, потому что в евенте gamestart боту не известно пока свое место
	id := c.Param("id")

	// Get the string associated with the key from the cache
	idVal, found := mcache.Get(id)
	if !found {
		//это первый постучавшийся бот, сохраним id раздачи в базе и отправим боту, что он главный
		mcache.Set(id, 1, cache.DefaultExpiration)

		c.String(200, "you are main")
	} else {
		//это не первый бот - увеличим значение для поля на 1 и отправим боту Удачи!
		//в значении хранится количество ботов в итоге
		i := idVal.(int)
		mcache.Set(id, i+1, cache.DefaultExpiration)
		c.String(200, "Good luck!")
	}
}

//бот из academy говорит, что ему раздали карты
func holeCards(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	println(botname + " : " + c.Param("id") + " | " + c.Param("cards"))

}

//главный бот из academy говорит, что всем раздали карты
func dealHoleCardsEventPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)

	// defining a struct instance
	var tablestate TableState

	// JSON array to be decoded
	// to an array in golang
	Data := []byte(c.PostForm("JSON"))

	// decoding JSON array to
	// the country array
	err := json.Unmarshal(Data, &tablestate)

	if err != nil {

		// if error is not nil
		// print error
		println(err)
	}

	// printing decoded array
	// values one by one

	print(botname + " | ")
	fmt.Println(tablestate)

}

type TableState struct {
	HandID  int64    `json:"handID"`
	State   string   `json:"state"`
	BB      float64  `json:"BB"`
	Ante    float64  `json:"ante"`
	Players []Player `json:"players"`
}

type Player struct {
	Name  string  `json:"name"`
	Stack float64 `json:"stack"`
}

//главный бот из academy говорит, что поменялся уровень раздачи preflop flop turn river showdown
func stageEvent(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	println(botname + " : " + c.Param("id") + " | " + c.Param("stage") + " | " + c.Param("board"))

}

//бот из academy спрашивает свое действие
func getActionPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")

	var pg PlayerGetAction

	Data := []byte(c.PostForm("JSON"))
	err := json.Unmarshal(Data, &pg)

	if err != nil {
		println(err)
	}

	print(botname + " | ")
	fmt.Println(pg)

	//ответ надо сгенерить в одном из вариантов:
	//FOLD
	//CHECK
	//CALL
	//BET:40.0
	//RAISE:120.0

	//Заглушим пока ответы простыми штуками
	if pg.ToCall == 0 {
		c.String(200, "CHECK")
	} else if pg.ToCall == pg.BB {
		c.String(200, "CALL")
	}

}

type PlayerGetAction struct {
	HandID     int64   `json:"handID"`
	BB         float64 `json:"BB"`
	ToCall     float64 `json:"tocall"`
	MinRaise   float64 `json:"minraise"`
	TotalPot   float64 `json:"totalpot"`
	MainPot    float64 `json:"mainpot"`
	Stack      float64 `json:"stack"`
	NumRaise   int64   `json:"numtraise"`
	NumPlayers int64   `json:"numplayers"`
}

//главный бот из academy говорит, что кто-то действие призвел
func actionEventPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)

	var pa PlayerAction

	Data := []byte(c.PostForm("JSON"))
	err := json.Unmarshal(Data, &pa)

	if err != nil {
		println(err)
	}

	print(botname + " | ")
	fmt.Println(pa)
}

//{"handID":271,"player":"Grant A.","action":"BIG_BLIND","amount":20.0}
type PlayerAction struct {
	HandID     int64   `json:"handID"`
	PlayerName string  `json:"player"`
	Action     string  `json:"action"`
	Amount     float64 `json:"amount"`
	TotalPot   float64 `json:"totalpot"`
	MainPot    float64 `json:"mainpot"`
}

//главный бот из academy говорит, что кто-то показал карты
func showdownEventPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)

	var ps PlayerShowCard

	Data := []byte(c.PostForm("JSON"))
	err := json.Unmarshal(Data, &ps)

	if err != nil {
		println(err)
	}

	print(botname + " | ")
	fmt.Println(ps)

}

//{"handID":270,"player":"Grant A.","cards":"KsJh"}
type PlayerShowCard struct {
	HandID     int64  `json:"handID"`
	PlayerName string `json:"player"`
	Cards      string `json:"cards"`
}

//главный бот из academy говорит, что кто-то выиграл раздачу
func winEventPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)

	var pw PlayerWin

	Data := []byte(c.PostForm("JSON"))
	err := json.Unmarshal(Data, &pw)

	if err != nil {
		println(err)
	}

	print(botname + " | ")
	fmt.Println(pw)

}

//{"handID":270,"player":"Grant A.","amount":340.0,"handname":"Two Pair, Aces and Kings, Jack kicker"}
type PlayerWin struct {
	HandID     int64   `json:"handID"`
	PlayerName string  `json:"player"`
	Amount     float64 `json:"amount"`
	Handname   string  `json:"handname"`
}

//главный бот из academy говорит, что раздача окончена
func gameOverEvent(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	id := c.Param("id")

	//удалим из кеша эту раздачу
	mcache.Delete(id)

	println(botname + " : " + id + " | hand finished")

}

package main

import (
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
	r.GET(securekey+"/stageEvent/:id/:stage", stageEvent)

	r.GET(securekey+"/showdownEvent/:id/:playername/:cards", showdownEvent)
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
	println(botname + " : " + c.PostForm("JSON"))

}

//главный бот из academy говорит, что поменялся уровень раздачи preflop flop turn river showdown
func stageEvent(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	println(botname + " : " + c.Param("id") + " | " + c.Param("stage"))

}

//главный бот из academy говорит, что кто-то показал карты
func showdownEvent(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	println(botname + " : " + c.Param("id") + " | " + c.Param("playername") + " | " + c.Param("cards"))

}

//главный бот из academy говорит, что кто-то выиграл раздачу
func winEventPost(c *gin.Context) {
	botname := c.Request.Header.Get("User-Agent")
	c.Status(200)
	println(botname + " : " + c.PostForm("JSON"))

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

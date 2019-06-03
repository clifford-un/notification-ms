package controllers

import (
	"notification-ms/app"
	"strconv"

	"github.com/revel/revel"
)

type Test struct {
	*revel.Controller
}

type testModel struct {
	id   int64
	test string
}

func getTestRedis(id int64) testModel {
	idString := strconv.FormatInt(id, 10)
	test, err := app.RedisClient.Get(idString).Result()
	if err != nil {
		test = "None"
	}
	return testModel{id, test}
}

func (test Test) Index() revel.Result {
	data := make(map[string]interface{})
	data["id"] = "1"
	data["test"] = getTestRedis(1).test
	return test.RenderJSON(data)
}

func (test Test) Show() revel.Result {
	data := make(map[string]interface{})
	data["id"] = test.Params.Route.Get("id")
	data["test"] = getTestRedis(1).test
	return test.RenderJSON(data)
}

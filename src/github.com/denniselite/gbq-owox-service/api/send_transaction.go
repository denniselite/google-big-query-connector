package api

import (
	"log"
	"encoding/json"
	"gopkg.in/validator.v2"
	"github.com/kataras/iris"
	. "github.com/denniselite/toolkit/errors"
	"github.com/denniselite/gbq-owox-service/structs"
	. "github.com/denniselite/gbq-owox-service/manager"
)

func (c *Context) Data(ctx *iris.Context) {

	// Создаем уникальный ID запроса для отслеживания в менеджере,
	// разбираем запрос, валидируем
	uuid := ctx.GetString("uuid")
	var signature structs.SendRequest
	if err := ctx.ReadJSON(&signature); err != nil {
		ctx.JSON(HttpApiError(NewError(err, BadRequest)))
		return
	}

	log.Printf("%s Signature: %#v\n", uuid, signature)
	if err := validator.Validate(signature); err != nil {
		ctx.JSON(HttpApiError(err))
		return
	}

	// Вызыываем через RPC rabbitmq нужный метод менеджера и отдаем результат
	body, _ := json.Marshal(signature)
	res, err := c.Rmq.Rpc(QueueSend, uuid, body)
	if err != nil {
		ctx.JSON(HttpApiError(err))
		return
	}

	ctx.JSON(iris.StatusNoContent, res)
}

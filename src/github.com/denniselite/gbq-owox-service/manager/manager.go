package manager

import (
	"encoding/json"
	"gopkg.in/validator.v2"
	"github.com/jinzhu/gorm"
	"github.com/denniselite/toolkit/conn"
	. "github.com/toolkit/errors"
	"github.com/denniselite/gbq-owox-service/structs"
	"golang.org/x/net/context"
	"cloud.google.com/go/bigquery"
	"strconv"
	"time"
)


const (
	QueueSend = "gbqm.send"
)

type GoogleBiqQueryManager struct {
	Db              *gorm.DB
	Rmq             conn.RmqInt
	BigQueryCtx     context.Context
	Client     		*bigquery.Client
	DataSet         *bigquery.Dataset
	ProjectID       string
	DataSetName     string
	TableName		string
}

// Инициализация менеджера
func (m *GoogleBiqQueryManager) Run(rmq conn.RmqInt) {
	m.Rmq = rmq

	// Создаем BQ-контекст и клиент для требуемого dataSet-а
	var err error
	m.BigQueryCtx = context.Background()
	m.Client, err = bigquery.NewClient(m.BigQueryCtx, m.ProjectID)
	Oops(err)
	m.DataSet = m.Client.DatasetInProject(m.ProjectID, m.DataSetName)

	go func() {
		Oops(m.Rmq.ConsumeRpc(QueueSendTransaction, m.SendTransaction, -1))
	}()
}

// Отправка входящей транзакции в Google Biq Query
func (m *GoogleBiqQueryManager) SendTransaction(data []byte) (res conn.RmqMessage, err error)  {

	// Разбираем запрос
	rq := new(structs.SendRequest)
	err = json.Unmarshal(data, rq)
	if err != nil {
		return
	}

	err = validator.Validate(rq)
	if err != nil {
		return
	}

	// Создаем коннект к таблице и загрузчик
	table := m.DataSet.Table(m.TableName)
	u := table.Uploader()

	// Заполняем структуру транзакции
	item := new(interface{})

	var items []*interface{}
	items = append(items, item)

	// Загружаем, если все ок - отдаем пустой ответ и статус - 200
	if err = u.Put(m.BigQueryCtx, items); err != nil {
		return
	}

	response := new(structs.EmptyResponse)
	res = response
	return
}

// Получение переменной типа time.Time по количеству миллисекунд
// Вычисляем целое количество секунд и наносекунд,
// через time.Unix(...) получаем необходимую переменную типа time.Time;
func (m *GoogleBiqQueryManager) getTimeFromMilliseconds(ms int64) (dateTime time.Time) {
	sec := ms / 1000
	nsec := (ms - (sec) * 1000) * 1000000
	dateTime = time.Unix(sec, nsec)
	return
}
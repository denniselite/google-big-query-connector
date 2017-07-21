package main

import (
	"log"
	"fmt"
	"flag"
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"github.com/denniselite/toolkit/conn"
	. "github.com/toolkit/errors"
	"github.com/denniselite/gbq-owox-service/manager"
	"github.com/denniselite/gbq-owox-service/structs"
	"github.com/denniselite/gbq-owox-service/api"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "config", "", "config filename")
	flag.StringVar(&configFile, "c", "", "config filename (shorthand)")
	flag.Parse()

	cfg := new(structs.Config)
	data, err := ioutil.ReadFile(configFile)
	Oops(err)

	Oops(yaml.Unmarshal(data, &cfg))

	// get rabbit connection
	var rmq *conn.Rmq

	// get rabbit connection
	rabbitConnectionString := fmt.Sprintf("amqp://%s:%s@%s:%d",
		cfg.Rabbit.Username,
		cfg.Rabbit.Password,
		cfg.Rabbit.Host,
		cfg.Rabbit.Port,
	)

	man := new(manager.GoogleBiqQueryManager)
	man.ProjectID = cfg.BigQuery.ProjectID
	man.DataSetName = cfg.BigQuery.DataSet
	man.TableName = cfg.BigQuery.TableName
	fmt.Printf("Google Big Query connection params:\n===================================\nProject: %s\nDataset: %s\nTable: %s\n===================================\n", man.ProjectID, man.DataSetName, man.TableName)

	run := func() {
		man.Run(rmq)
	}

	rmq, err = conn.NewRmq(rabbitConnectionString, run)
	if err != nil {
		log.Println("Failed connect to rabbit")
		Oops(err)
	}

	// run manager
	run()

	// setup api server
	ctx := &api.Context{
		Rmq: rmq,
		RmqString: rabbitConnectionString,
	}

	r := ctx.NewRouter()

	//run API application
	log.Printf("Listen HTTP port: %d", cfg.Listen)
	r.Listen(fmt.Sprintf(":%d", cfg.Listen))
}

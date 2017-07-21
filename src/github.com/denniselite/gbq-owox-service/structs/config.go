package structs

type Config struct {
	Listen          int
	Rabbit          RabbitConfig
	BigQuery        BigQueryConfig `yaml:"bigQuery"`
}

type RabbitConfig struct {
	Host        string
	Port        int
	Username    string
	Password    string
}

type BigQueryConfig struct {
	DataSet   string `yaml:"dataSet"`
	ProjectID   string `yaml:"projectID"`
	TableName string `yaml:"tableName"`
}
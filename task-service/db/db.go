package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jaustinmiles/ostodo/task-service/common"
	"github.com/jaustinmiles/ostodo/task-service/tasks"
	"io"
	"net/http"
	"os"
	"strconv"
)

type TaskDatabase struct {
	User   string
	Addr   string
	Port   int
	Region string
	Name   string
	db     *sql.DB
}

func NewTaskDatabase() *TaskDatabase {
	port, _ := strconv.Atoi(os.Getenv("TASK_DB_PORT"))
	return &TaskDatabase{
		User:   os.Getenv("TASK_DB_USER"),
		Addr:   os.Getenv("TASK_DB_ADDR"),
		Name:   os.Getenv("TASK_DB_NAME"),
		Port:   port,
		Region: os.Getenv("AWS_REGION"),
	}
}

func (taskDb *TaskDatabase) Connect() error {
	l := common.GetLogger()
	l.Info("attempting to register rds certs with mysql")
	err := RegisterRDSMysqlCerts(http.DefaultClient)
	if err != nil {
		l.Error("error registering certs to rds")
		return err
	}
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		l.Error("Error loading default configuration", err)
		return err
	}
	dbEndpoint := fmt.Sprintf("%s:%d", taskDb.Addr, taskDb.Port)
	authToken, err := auth.BuildAuthToken(
		context.Background(),
		dbEndpoint,
		taskDb.Region,
		taskDb.User,
		cfg.Credentials,
	)
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?tls=rds&allowCleartextPasswords=true",
		taskDb.User, authToken, dbEndpoint, taskDb.Name,
	)
	l.Info("opening database with access string: ", dsn)
	db, err := sql.Open("mysql", dsn)
	taskDb.db = db
	return err
}

func (taskDb *TaskDatabase) InsertTask(task tasks.Task) error {
	l := common.GetLogger()
	q := "INSERT INTO `tasks` (user, uuid, name, completiontime, repetitions) VALUES (?, ?, ?, ?, ?)"
	insert, err := taskDb.db.Prepare(q)
	if err != nil {
		l.Error("failed to prepare db query for inserting task")
		return err
	}
	defer insert.Close()

	_, err = insert.Exec(task.User, task.UUID, task.Name, task.CompletionTime, task.Repetitions)
	return err
}

func (taskDb *TaskDatabase) Close() {
	defer taskDb.db.Close()
}

func RegisterRDSMysqlCerts(c *http.Client) error {
	resp, err := c.Get("https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem")
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	pem, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	rootCertPool := x509.NewCertPool()
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		return fmt.Errorf("couldn't append certs from pem")
	}

	err = mysql.RegisterTLSConfig("rds", &tls.Config{RootCAs: rootCertPool, InsecureSkipVerify: true})
	if err != nil {
		return err
	}
	return nil
}

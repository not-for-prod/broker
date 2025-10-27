package main

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
	"github.com/not-for-prod/broker/consumer"
	"github.com/not-for-prod/broker/producer"
	"github.com/stretchr/testify/suite"
)

const (
	natsURL = "nats://localhost:4222"
	pgDSN   = "postgres://postgres:postgres@localhost/postgres?sslmode=disabled"
)

type TestSuite struct {
	suite.Suite
	nats      *nats.Conn
	jetStream nats.JetStream
	db        *pgxpool.Pool
	producer  *producer.Producer
	consumer  *consumer.Consumer
}

func (suite *TestSuite) SetupSuite() {
	var err error

	suite.nats, err = nats.Connect(natsURL)
	suite.Require().NoError(err)

	suite.jetStream, err = suite.nats.JetStream()
	suite.Require().NoError(err)

	suite.db, err = pgxpool.New(context.Background(), pgDSN)
	suite.Require().NoError(err)
}

func (suite *TestSuite) TearDownSuite() {
	suite.nats.Close()
	suite.db.Close()
}

func (suite *TestSuite) TestDelivery() {}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

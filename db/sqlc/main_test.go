package db

import (
	"Ice/internal/config"
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

var testStore Store

func TestMain(m *testing.M) {
	var cm config.ConfigManager
	conf := cm.LoadConfig("../../internal/config/.")

	connPool, err := pgxpool.New(context.Background(), conf.Postgres.DatabaseSource())
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}

func TestCreateUser(t *testing.T) {
	user, err := testStore.CreateUser(context.Background(), &CreateUserParams{
		Username: "zhansgan",
		Nickname: "zhansgan",
		Password: "zhansgan",
		Email:    "zhansgan@example.com",
		Gender:   1,
	})
	require.NoError(t, err)
	require.NotEmpty(t, user)
}

package user

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go-app/config"
	"go-app/database"
	"go-app/domain"
	"log"
	"testing"
	"time"
)

var pgConStr string

func init() {
	ctx := context.Background()

	c, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		panic(err)
	}

	connStr, err := c.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		panic(err)
	}

	// check the connection to the database
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)

	pgConStr = connStr
	log.Println(fmt.Sprintf("Users database started successfully. Connection String: %s", pgConStr))
}

func Test_Should_Create_User(t *testing.T) {

	gormDb := config.ConnectTestPostgres(pgConStr)
	database.Migrate(gormDb)

	userRepo := NewUserRepository(gormDb)

	user := domain.User{Name: "mert", Age: 26}
	savedUser, err := userRepo.CreateUser(user)
	if err != nil {
		log.Fatal(err.Message)
	}

	assert.Equal(t, user.Name, savedUser.Name)
}

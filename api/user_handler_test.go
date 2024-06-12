package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/db"
	"github.com/nickemma/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type testDB struct {
	db.UserStore
}

func (tdb *testDB) teardown(t *testing.T) {
	if err := tdb.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func setUptest(t *testing.T) *testDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoUri))
	if err != nil {
		t.Fatal(err)
	}

	return &testDB{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestCreateUser(t *testing.T) {
	tdb := setUptest(t)
	defer tdb.teardown(t)

	app := fiber.New()
	UserHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", UserHandler.HandleCreateUser)

	params := types.CreateUserParams{
		FirstName: "test",
		LastName:  "tester",
		Email:     "test@gmail.com",
		Password:  "Test1234#",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	response, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user types.User
	json.NewDecoder(response.Body).Decode(&user)
	if user.FirstName != params.FirstName {
		t.Errorf("expected %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected %s but got %s", params.Email, user.Email)
	}
}

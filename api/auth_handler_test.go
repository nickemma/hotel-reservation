package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/nickemma/hotel-reservation/db"
	"github.com/nickemma/hotel-reservation/types"
)

func testInsertUser(t *testing.T, userStore db.UserStore) *types.User {
	user, err := types.CreateNewUserFromParams(types.CreateUserParams{
		Email:     "tester@gmail.com",
		FirstName: "tester",
		LastName:  "test",
		Password:  "Test1234#",
	})

	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user
}

func TestAuthenticate(t *testing.T) {
	tdb := setUptest(t)
	defer tdb.teardown(t)
	insertedUser := testInsertUser(t, tdb.UserStore)

	app := fiber.New()
	AuthHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/auth", AuthHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "tester@gmail.com",
		Password: "Test1234#",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status code %d, got %d", http.StatusOK, res.StatusCode)
	}

	var authRes AuthResponse
	if err := json.NewDecoder(res.Body).Decode(&authRes); err != nil {
		t.Fatal(err)
	}
	if authRes.Token == "" {
		t.Fatal("expected token, got empty string")
	}
	insertedUser.EncryptedPassword = ""

	if !reflect.DeepEqual(authRes.User, insertedUser) {
		t.Fatal("expected user to be the same as the inserted user")
	}
}

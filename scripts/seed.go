package main

import (
	"context"
	"log"

	"github.com/nickemma/hotel-reservation/db"
	"github.com/nickemma/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(fname, lName, email string) {
	user, err := types.CreateNewUserFromParams(types.CreateUserParams{
		FirstName: fname,
		LastName:  lName,
		Email:     email,
		Password:  "Tech1234#",
	})

	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.CreateUser(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {
	hotel := types.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Size:  "Single Room",
			Price: 99.9,
		},
		{
			Size:  "Double Room",
			Price: 499.9,
		},
		{
			Size:  "Deluxe Room",
			Price: 999.9,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID

		_, err := roomStore.CreateRoom(ctx, &room)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Paris", "France", 3)
	seedHotel("The cozy hotel", "The Nederlands", 4)
	seedHotel("Don't Give up", "Rwanda", 1)
	seedUser("Techie", "Emma", "Techieemma@gmail.com")
}

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoUri))
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
	userStore = db.NewMongoUserStore(client)
}

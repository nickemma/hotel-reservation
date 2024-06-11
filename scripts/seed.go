package main

import (
	"context"
	"fmt"
	"log"

	"github.com/nickemma/hotel-reservation/db"
	"github.com/nickemma/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()
	// connection to mongodb
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.MongoUri))

	// check if there is an error while connecting to mongodb
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Database(db.DBNAME).Drop(ctx); err != nil {
		log.Fatal(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	hotel := types.Hotel{
		Name:     "Paris",
		Location: "France",
		Rooms:    []primitive.ObjectID{},
	}
	insertedHotel, err := hotelStore.CreateHotel(ctx, &hotel)

	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.Room{
		{
			Type:      types.SingleRoomType,
			BasePrice: 99.9,
		},
		{
			Type:      types.SeaSideRoomType,
			BasePrice: 499.9,
		},
		{
			Type:      types.DeluxeRoomType,
			BasePrice: 999.9,
		},
	}

	for _, room := range rooms {
		room.HotelID = insertedHotel.ID

		insertedRoom, err := roomStore.CreateRoom(ctx, &room)

		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("seeding the database", insertedRoom)
	}
}

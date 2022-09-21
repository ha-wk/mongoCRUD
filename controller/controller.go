package controller

import (
	"context"
	"encoding/json"
	"example/mongodemo/model"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connection = "mongodb+srv://tuna:tuna123@cluster0.m1u5qqj.mongodb.net/?retryWrites=true&w=majority"

const dbname = "NETFLIX"
const colname = "watchlist"

var collection *mongo.Collection

func init() {

	clientop := options.Client().ApplyURI(connection)

	client, err := mongo.Connect(context.TODO(), clientop)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connection Established")

	collection = client.Database(dbname).Collection(colname)

	fmt.Println("Collection instance is ready")
}

//MONGODB HELPERS

func inserOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)

	if err != nil {
		log.Fatal(err)

	}
	fmt.Println("Inserted 1 movie in db with id:", inserted.InsertedID)
}

//UPDATE RECORD

func updateOneMovie(movieId string) {

	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("modified count:", result.ModifiedCount)
}

//delete one movie

func deleteOneMovie(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	delteCount, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("movie got deleted with dletecount as", delteCount)
}

///delete many

func deleteAllMovie() int64 {
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("number of movies delete:", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

//extracting all movies

func getAllMovies() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}

	var movies []primitive.M

	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies
}

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

//CREATE AND UPDATE IN DATABASE THROUGH GOLANG

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	_ = json.NewDecoder(r.Body).Decode(&movie)
	inserOneMovie(movie)
	json.NewEncoder(w).Encode(movie)
}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-control-Allow-Methods", "POST")

	params := mux.Vars(r)

	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-control-Allow-Methods", "DELETE")

	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}

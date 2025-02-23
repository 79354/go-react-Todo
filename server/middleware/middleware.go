package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-react-Todo/server/models"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

/*Go automatically detects and executes init() before anything else in that package.
	If multiple init() functions exist in different files of the same package, they run in an unspecified order but before main()*/

func init(){
	loadTheEnv()
	createDBInstance()
}

func loadTheEnv(){
	err:= godotenv.Load(".env")
	if err != nil{
		log.Fatal("Error loading the .env file")
	}
}

func createDBInstance(){
	// Retrieve connection details from .env
	connectionString:= os.Getenv("DB_URI")
	dbName:= os.Getenv("DB_NAME")
	collName:= os.Getenv("DB_COLLECTION_NAME")

	// Configure mongodb client and establish connection, context.TODO() place holder	
	clientOptions:= options.Client().ApplyURI(connectionString)
	client, err:= mongo.Connect(context.TODO(), clientOptions)
	if err != nil{
		log.Fatal(err)
	}

	// Ping Client is used to test the reachability of a network device on an IP network
	client.Ping(context.TODO(), nil)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("connected to mongodb!")

	/*client is the entire system (MongoDB connection).
	client.Database("LibraryDB") selects a specific library (database).
	client.Database("LibraryDB").Collection("Fiction") picks a bookshelf (collection).
	Now you can read, add, remove books (documents in the collection).*/

	collection = client.Database(dbName).Collection(collName)
	fmt.Println("collection instance created")
}

func GetAllTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www.form urlencoded")		// multipart/form-data : for significant size payload
	w.Header().Set("Acess-Control-Allow-Origin", "*")
	payload:= getAllTasks()
	json.NewEncoder(w).Encode(payload)
}

func CreateTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")	// allow request from anywhere use *
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")	// CORS header key

	var task models.ToDoList

	// Converts incoming JSON to struct
	json.NewDecoder(r.Body).Decode(&task)

	//Inserts into MongoDB (Go struct → BSON)
	insertOneTask(task)

	// Sends back the created task as JSON
	json.NewEncoder(w).Encode(task)
}

func TaskComplete(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params:= mux.Vars(r)
	TaskComplete(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func UndoTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods","PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params:= mux.Vars(r)
	UndoTask(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteTask(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "applicaiton/x-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	params:= mux.Vars(r)
	deleteOneTask(params["id"])
}

func DeleteAllTasks(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "applicaiton/x-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	count := deleteAllTasks()
	json.NewEncoder(w).Encode(count)
}


// Getting Tasks (Database → Server → Client):
// returns a slice(map[string]interface{})
func getAllTasks()[]primitive.M{
	// Fetches BSON documents from MongoDB
	// passing a empty query .D{{}} -> to fetch everything in the collection
	cur, err:= collection.Find(context.Background(), bson.D{{}})
	if err != nil{
		log.Fatal(err)
	}

	var results []primitive.M
	// Converts BSON to Go structures
	for cur.Next(context.Background()){
		var result bson.M
		e:= cur.Decode(&result)
		if e != nil{
			log.Fatal(e)
		}
		results = append(results, result)
	}

	if err:= cur.Err(); err != nil{
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return results
}

func taskComplete(task string){
	id, _ := primitive.ObjectIDFromHex(task)
	filter:= bson.M{"id":id}
	update:= bson.M{"$set":bson.M{"status":true}}
	result, err:= collection.UpdateOne(context.Background())
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("modified count:", result.ModifiedCount) 

}

func insertOneTask(task models.ToDoList){
	InsertResult, err:= collection.InsertOne(context.Background(), task)
	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("Inserted a single result: ", InsertResult.InsertedID)
}

func undoTask(task string){
	id := primitive.ObjectIDFromHex(task)
	filter:= bson.M{"id": id}
	update:= bson.M{"$set":bson.M{}}
	result, err:= collection.UpdateOne(context.Background(), filter, update)
	fmt.Println("modified count:", result)
}

func deleteOneTask(task string){
	id, _ := primitive.ObjectIDFromHex(task)
	filter:= bson.M{"id": id}	// M is an unordered representation of a BSON document.
	d, err:= collection.DeleteOne(context.Background(), filter)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("Deleted task: ", d.DeletedCount)
}

func deleteAllTasks() int64{
	d, err:= collection.DeleteMany(context.Background(), bson.D{}, nil)
	if err != nil{
		log.Fatal(err)
	}
	fmt.Println("deleted document", )
	return d.DeletedCount
}

/*
	createDBInstance()
	1. Load environment variables (credentials like DB_URI, DB_NAME).
	2. Establish a connection to MongoDB using the mongo.Connect() function.
	3. Verify the connection using client.Ping().
	4. Access a specific database and collection where your data will be stored.
	5. Store the collection reference in a global variable (collection), so all functions can interact with it.


	CORS: need to be set for each method type(GET, POST, PUT, etc)
		w.Header().Set("Access-Control-Allow-Origin", "*")
	SOP: same origin policy, to overcome this security measure
*/
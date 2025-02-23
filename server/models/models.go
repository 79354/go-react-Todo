package model

import(
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ToDoList struct{
	ID		primitive.ObjectID	`json:"_id,omitempty" bson:"_id,omitempty"`
	Task	string				`json:"task,omitempty"`
	Status	bool				`json:"status,omitempty"`	// omits zero value such as Status: false, for compactness
}

/*
	The data transformations are:

	Client → Server: JSON → Go struct
	Server → Database: Go struct → BSON
	Database → Server: BSON → Go struct
	Server → Client: Go struct → JSON

The encoding/json package and MongoDB driver handle these conversions automatically based on the struct tags,
*/

/*
	what does the struct do?
	1. ID: A unique MongoDB identifier
	2. Task: The actual task text
	3. Status: A completion status flag


	primitive.ObjectID - MongoDB's unique identifier type

	The struct tags (json and bson) are metadata that tell the JSON encoder/decoder 
	and MongoDB driver how to convert the struct to and from their respective formats.

	BSON (Binary JSON) is MongoDB's binary-encoded serialization format
*/
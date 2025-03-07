package helpers

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ObjectIDToString(id primitive.ObjectID) string {
	return id.Hex()
}

func StringToObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

package handlers

import (
	"context"

	"github.com/RIpallol541/PeogectX/db"
	"github.com/RIpallol541/PeogectX/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// HandleAuth handles user authentication
func HandleAuth(user models.User) map[string]string {
	collection := db.Client.Database("myapp").Collection("users")
	var dbUser models.User
	filter := bson.M{"username": user.Username}
	err := collection.FindOne(context.TODO(), filter).Decode(&dbUser)
	if err != nil || user.Password != dbUser.Password {
		return map[string]string{"status": "Authentication failed"}
	}

	return map[string]string{"status": "Authentication successful"}
}
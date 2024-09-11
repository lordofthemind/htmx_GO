// In your `initializers` package
package initializers

import (
	"go.mongodb.org/mongo-driver/mongo"
)

// GetDatabase returns the MongoDB database instance for the given database name.
func GetDatabase(client *mongo.Client, dbName string) *mongo.Database {
	return client.Database(dbName)
}

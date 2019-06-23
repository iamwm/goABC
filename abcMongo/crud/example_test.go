package crud

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

func TestDocumentationExamples(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cs := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cs))

	require.NoError(t, err)
	defer client.Disconnect(ctx)

	db := client.Database("documentation_examples")

	InsertExamples(t, db)
	QueryToplevelFieldsExamples(t, db)
	QueryEmbeddedDocumentsExamples(t, db)
	QueryArrayExamples(t, db)
	QueryArrayEmbeddedDocumentsExamples(t, db)
}

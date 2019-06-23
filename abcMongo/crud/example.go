package crud

import (
	"context"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)
import "go.mongodb.org/mongo-driver/mongo"

func requireCursorLength(t *testing.T, cursor *mongo.Cursor, length int) {
	i := 0
	for cursor.Next(context.Background()) {
		i++
	}

	require.NoError(t, cursor.Err())
	require.Equal(t, i, length)
}

func containsKey(doc bson.Raw, key ...string) bool {
	_, err := doc.LookupErr(key...)
	if err != nil {
		return false
	}
	return true
}

func InsertExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_insert")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		result, err := coll.InsertOne(
			context.Background(),
			bson.D{
				{"item", "canvas"},
				{"qty", 100},
				{"tags", bson.A{"cotton"}},
				{"size", bson.D{
					{"h", 28},
					{"w", 35.5},
					{"uom", "cm"},
				}},
			},
		)

		require.NoError(t, err)
		require.NotNil(t, result.InsertedID)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"item", "canvas"},
			},
		)

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		result, err := coll.InsertMany(
			context.Background(),
			[]interface{}{
				bson.D{
					{"item", "journal"},
					{"qty", int32(25)},
					{"tags", bson.A{"blank", "red"}},
					{"size", bson.D{
						{"h", 14},
						{"w", 21},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mat"},
					{"qty", int32(25)},
					{"tags", bson.A{"gray"}},
					{"size", bson.D{
						{"h", 27.9},
						{"w", 35.5},
						{"uom", "cm"},
					}},
				},
				bson.D{
					{"item", "mousepad"},
					{"qty", 25},
					{"tags", bson.A{"gel", "blue"}},
					{"size", bson.D{
						{"h", 19},
						{"w", 22.85},
						{"uom", "cm"},
					}},
				},
			})
		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 3)
	}
}

func QueryToplevelFieldsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_top")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{},
		)

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", "D"}},
		)

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"status", bson.D{{"$in", bson.A{"A", "D"}}}}})
		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
				{"qty", bson.D{{"$lt", 30}}},
			})
		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{
					"$or",
					bson.A{
						bson.D{{"status", "A"}},
						bson.D{{"qty", bson.D{{"$lt", 30}}}},
					},
				},
			})
		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"status", "A"},
				{"$or", bson.A{
					bson.D{{"qty", bson.D{{"$lt", 30}}}},
					bson.D{{"item", primitive.Regex{Pattern: "^p", Options: ""}}},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}
}

func QueryEmbeddedDocumentsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_embedded")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "A"},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"size", bson.D{
					{"h", 8.5},
					{"w", 11},
					{"uom", "in"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"size", bson.D{
					{"h", 22.85},
					{"w", 30},
					{"uom", "cm"},
				}},
				{"status", "D"},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"size", bson.D{
					{"h", 10},
					{"w", 15.25},
					{"uom", "cm"},
				}},
				{"status", "A"},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size", bson.D{
					{"h", 14},
					{"w", 21},
					{"uom", "cm"},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size.uom", "in"},
			}, )

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"size.h", bson.D{
					{"$lt", 15},
				}},
				{"size.uom", "in"},
				{"status", "D"},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}
}

func QueryArrayExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_array")

	err := coll.Drop(context.Background())
	require.NoError(t, err)

	{
		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"qty", 25},
				{"tags", bson.A{"black", "red"}},
				{"dim_cm", bson.A{14, 21}},
			},
			bson.D{
				{"item", "notebook"},
				{"qty", 50},
				{"tags", bson.A{"red", "black"}},
				{"dim_cm", bson.A{14, 21}},
			},
			bson.D{
				{"item", "paper"},
				{"qty", 100},
				{"tags", bson.A{"red", "black", "plain"}},
				{"dim_cm", bson.A{14, 21}},
			},
			bson.D{
				{"item", "planner"},
				{"qty", 75},
				{"tags", bson.A{"black", "red"}},
				{"dim_cm", bson.A{22.85, 30}},
			},
			bson.D{
				{"item", "postcard"},
				{"qty", 45},
				{"tags", bson.A{"blue"}},
				{"dim_cm", bson.A{10, 15.25}},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"tags", bson.A{"red", "black"}}}, )

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{{"tags", bson.D{{"$all", bson.A{"red", "black"}}}}})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"tags", "red"},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"dim_cm", bson.D{{"$gt", 25},}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"dim_cm", bson.D{
					{"$gt", 15},
					{"$lt", 20},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"tags", bson.D{
					{"$size", 3},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}
}

func QueryArrayEmbeddedDocumentsExamples(t *testing.T, db *mongo.Database) {
	coll := db.Collection("inventory_query_array_embedded")

	err := db.Drop(context.Background())
	require.NoError(t, err)

	{
		docs := []interface{}{
			bson.D{
				{"item", "journal"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 5},
					},
					bson.D{
						{"warehouse", "C"},
						{"qty", 15},
					},
				}},
			},
			bson.D{
				{"item", "notebook"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "C"},
						{"qty", 5},
					},
				}},
			},
			bson.D{
				{"item", "paper"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 60},
					},
					bson.D{
						{"warehouse", "B"},
						{"qty", 15},
					},
				}},
			},
			bson.D{
				{"item", "planner"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "A"},
						{"qty", 40},
					},
					bson.D{
						{"warehouse", "B"},
						{"qty", 5},
					},
				}},
			},
			bson.D{
				{"item", "postcard"},
				{"instock", bson.A{
					bson.D{
						{"warehouse", "B"},
						{"qty", 15},
					},
					bson.D{
						{"warehouse", "C"},
						{"qty", 35},
					},
				}},
			},
		}

		result, err := coll.InsertMany(context.Background(), docs)

		require.NoError(t, err)
		require.Len(t, result.InsertedIDs, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{
					"instock", bson.D{
					{"warehouse", "A"},
					{"qty", 5},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 1)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock", bson.D{
					{"qty", 5},
					{"warehouse", "A"},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 0)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.0.qty", bson.D{
					{"$lte", 20},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.qty", bson.D{
					{"$lte", 20},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 5)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock", bson.D{
					{"$elemMatch", bson.D{
						{"qty", bson.D{
							{"$gt", 10},
							{"$lt", 20},
						}},
					}},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 3)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.qty", bson.D{
					{"$gt", 10},
					{"$lt", 20},
				}},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 4)
	}

	{
		cursor, err := coll.Find(
			context.Background(),
			bson.D{
				{"instock.qty", 5},
				{"instock.warehouse", "A"},
			})

		require.NoError(t, err)
		requireCursorLength(t, cursor, 2)
	}
}

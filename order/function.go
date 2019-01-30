package order

import (
	"context"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/jsonapi"
	"github.com/jonaskay/cola-inventory-functions/config"
)

// Create saves a new Order entity to Datastore and prints it as a JSON object.
func Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", jsonapi.MediaType)

	ctx := context.Background()

	dsClient, err := datastore.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}

	orders := make([]*config.Order, 0)
	q := datastore.NewQuery("Order").Order("-CreatedAt").Limit(1)

	if _, err := dsClient.GetAll(ctx, q, &orders); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}

	if len(orders) > 0 && orders[0].DeliveredAt.IsZero() {
		w.WriteHeader(http.StatusForbidden)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "An existing order is waiting for delivery",
			Status: "403",
		}})
		return
	}

	orderKey := datastore.IncompleteKey("Order", nil)
	order := &config.Order{CreatedAt: time.Now()}

	key, err := dsClient.Put(ctx, orderKey, order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}
	order.ID = key.ID

	if err := jsonapi.MarshalPayload(w, order); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}
}

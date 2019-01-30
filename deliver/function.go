package deliver

import (
	"context"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/datastore"
	"github.com/google/jsonapi"
	"github.com/jonaskay/cola-inventory-functions/config"
)

// Latest marks the latest Order entity as delivered and prints it as a JSON
// object.
func Latest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "PATCH")
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

	keys, err := dsClient.GetAll(ctx, q, &orders)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}

	if len(keys) == 0 {
		w.WriteHeader(http.StatusForbidden)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Forbidden",
			Detail: "There are no orders waiting for delivery",
			Status: "403",
		}})
		return
	}

	k := keys[0]
	o := orders[0]
	o.ID = k.ID
	o.DeliveredAt = time.Now()

	if _, err := dsClient.Put(ctx, k, o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}

	if err := jsonapi.MarshalPayload(w, o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}
}

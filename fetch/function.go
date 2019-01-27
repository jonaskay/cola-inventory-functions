package fetch

import (
	"context"
	"net/http"
	"os"

	"cloud.google.com/go/datastore"
	"github.com/google/jsonapi"
	"github.com/jonaskay/cola-inventory-functions/config"
)

// Latest returns the latest Order entity and prints it as a JSON object.
func Latest(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
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
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Not Found",
			Status: "404",
		}})
		return
	}

	k := keys[0]
	o := orders[0]
	o.ID = k.ID

	if err := jsonapi.MarshalPayload(w, o); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
			Title:  "Internal Server Error",
			Status: "500",
		}})
		return
	}
}

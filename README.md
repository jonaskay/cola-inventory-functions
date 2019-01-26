# Cola Inventory Functions

This is a collection of **Google Cloud Functions for Go** that provide the backend for [Cola Inventory PWA](https://github.com/jonaskay/cola-inventory-pwa).

## Folder structure

* `config/` provides shared functionality for the project.
* `deliver/`, `fetch/`, `order/` contain the deployable Google Cloud Functions.

## Deploying

You can deploy the functions using the [`gcloud` command-line tool](https://cloud.google.com/sdk/gcloud/reference/functions/deploy):

    $ gcloud functions deploy NAME --entry-point ENTRY_POINT --set-env-vars PROJECT_ID=YOUR_PROJECT_ID --runtime go111 --trigger-http

## Development

Each function directory has a `cmd/` directory used for local testing. By running `cmd/main.go` you can start a local HTTP server which registers the parent function as an HTTP handler.

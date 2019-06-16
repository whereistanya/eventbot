// Package p contains an HTTP Cloud Function to update a database record
// to show when a slack user called this last.
package p

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"cloud.google.com/go/firestore"
)

// HelloWorld prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func Set(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Not a post: %v", r.Method)
		return
	}

	err := r.ParseForm()
	if err != nil {
		fmt.Fprintf(w, "Couldn't parse form: %s\n", err)
		return
	}

	username := r.Form["user_name"][0]
	text := r.Form["text"][0]
	command := r.Form["command"][0]

	projectID := "cloud-cloud-cloud-1359"

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		fmt.Fprintf(w, "Failed to create client: %v\n", err)
		return
	}

	defer client.Close()

	collection := client.Collection("afternoontea")
	doc := collection.Doc(username)
	current, err := doc.Get(ctx)
	if err != nil {
		fmt.Fprintf(w, "No existing data for %s: %v\n", username, err)
	} else {
		fmt.Fprintf(w, "Last one you told me about was:%+v\n", current.Data())
	}
	_, err = collection.Doc(username).Set(ctx, map[string]interface{}{
		"date":  time.Now(),
		"emoji": ":grimacing:",
	})
	if err != nil {
		fmt.Fprintf(w, "Aw, failed to write a document: %v\n", err)
		return
	}

	fmt.Fprintf(w, "Updated, %s [%s, %s]!\n", username, command, text)
}

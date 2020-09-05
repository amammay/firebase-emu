package main

import (
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/functions/metadata"
	"context"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"github.com/amammay/firebase-emu/fsemu"
	"log"
	"time"
)

func main() {

	event := fsemu.EmuResource{ProjectId: "mammay-play-test", Address: "http://localhost:8080"}
	emuRegisters := []fsemu.EmuRegister{{
		TriggerFn:    WriteSkills,
		TriggerType:  fsemu.FirestoreOnWrite,
		ResourcePath: "skills/{id}",
	}}

	if err := event.RegisterToEmu(emuRegisters); err != nil {
		panic(err)
	}

	if err := funcframework.Start("6000"); err != nil {
		panic(err)
	}

	ctx := context.Background()

	firestoreClient, err := firestore.NewClient(ctx, "mammay-play-test")
	if err != nil {
		log.Fatalf("firestore.NewClient: %v", err)
	}

	if _, err := firestoreClient.Collection("skills").NewDoc().Create(ctx, map[string]interface{}{
		"skill":     "Running",
		"timestamp": firestore.ServerTimestamp,
	}); err != nil {
		log.Fatalf("firestoreClient.Collection.NewDoc().Set: %v", err)
	}

	doc := firestoreClient.Doc("jimmy/stins")
	doc.Set(ctx, map[string]interface{}{
		"skill":     "Running",
		"timestamp": firestore.ServerTimestamp,
	})
	val, err := doc.Get(ctx)
	if err != nil {
		log.Fatalf("firestoreClient.Get: %v", err)
	}
	log.Println(val.Data())

}

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

// FirestoreValue holds Firestore fields.
type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log the interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     interface{} `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

func WriteSkills(ctx context.Context, e FirestoreEvent) error {

	meta, err := metadata.FromContext(ctx)
	if err != nil {
		return fmt.Errorf("metadata.FromContext: %v", err)
	}
	log.Printf("Function triggered by change to: %v", meta.Resource)
	log.Printf("Old value: %+v", e.OldValue)
	log.Printf("New value: %+v", e.Value)
	return nil
}

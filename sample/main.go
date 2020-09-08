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

	go func() {
		ctx := context.Background()

		firestoreClient, err := firestore.NewClient(ctx, "dummy")
		if err != nil {
			log.Fatalf("firestore.NewClient: %v", err)
		}
		defer firestoreClient.Close()

		if _, err := firestoreClient.Collection("skills").NewDoc().Create(ctx, map[string]interface{}{
			"skill":     "Running",
			"timestamp": firestore.ServerTimestamp,
		}); err != nil {
			log.Fatalf("firestoreClient.Collection.NewDoc().Set: %v", err)
		}
	}()
	event := fsemu.EmuResource{ProjectId: "dummy", Address: "http://localhost:8080"}
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
}

type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

type FirestoreValue struct {
	CreateTime time.Time   `json:"createTime"`
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

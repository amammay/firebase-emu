package fsemu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/GoogleCloudPlatform/functions-framework-go/funcframework"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

type EventType string

const (
	FirestoreOnCreate EventType = "providers/cloud.firestore/eventTypes/document.create"
	FirestoreOnUpdate EventType = "providers/cloud.firestore/eventTypes/document.update"
	FirestoreOnDelete EventType = "providers/cloud.firestore/eventTypes/document.delete"
	FirestoreOnWrite  EventType = "providers/cloud.firestore/eventTypes/document.write"
)

type ServiceType string

const (
	FirestoreService ServiceType = "firestore.googleapis.com"
)

type EmuTrigger struct {
	EventTrigger EventTrigger `json:"eventTrigger"`
}
type EventTrigger struct {
	Resource  string      `json:"resource"`
	EventType EventType   `json:"eventType"`
	Service   ServiceType `json:"service"`
}

type EmuRegister struct {
	TriggerFn    interface{}
	TriggerType  EventType
	ResourcePath string
}

type EmuResource struct {
	ProjectId string
	Address   string
}

func (events EmuResource) RegisterToEmu(registers []EmuRegister) error {
	httpClient := http.DefaultClient
	for _, element := range registers {

		emuTrigger := *events.createFsEmuTrigger(element)

		funcName := getFunctionName(element.TriggerFn)
		reqUrl := fmt.Sprintf("%s/emulator/v1/projects/%s/triggers/%s", events.Address,
			events.ProjectId, funcName)
		marshal, err := json.Marshal(emuTrigger)
		if err != nil {
			return fmt.Errorf("json.Marshal(): %v", err)
		}
		reader := bytes.NewReader(marshal)

		requester, err := http.NewRequest(http.MethodPut, reqUrl, reader)
		if err != nil {
			return fmt.Errorf("http.NewRequest(): %v", err)
		}

		emuResponse, err := httpClient.Do(requester)
		if err != nil {
			return fmt.Errorf("httpClient.Do(): %v", err)
		}
		if emuResponse.StatusCode != 200 {
			return fmt.Errorf("emuResponse status code: %v status: %s", emuResponse.StatusCode, emuResponse.Status)
		}

		funcPath := fmt.Sprintf("/functions/projects/%s/triggers/%s", events.ProjectId, funcName)
		funcframework.RegisterEventFunction(funcPath, element.TriggerFn)
	}

	return nil

}

func (events EmuResource) createFsEmuTrigger(element EmuRegister) *EmuTrigger {
	emuTrigger := &EmuTrigger{EventTrigger: EventTrigger{
		//fully qualified name for the emulator resource we will be registering against
		Resource: fmt.Sprintf("projects/%s/databases/(default)/documents/%s", events.ProjectId, element.ResourcePath),
		//what type emulator event we are respecting
		EventType: element.TriggerType,
		//what type of service (firestore/rtdb)
		Service: FirestoreService,
	}}
	return emuTrigger
}

func getFunctionName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	// firebase emu doesnt like path escapes so we replace / with dots
	name = strings.ReplaceAll(name, "/", ".")
	return name
}

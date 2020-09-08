# Introduction
The firebase team has done a really awesome job at providing a rich and feature complete environment for running and testing out various products all locally (firestore emulator, real time database emulator, pub/sub emulator).

# Purpose?
Most of the support/docs/guides out there mainly cover using the firebase-tools cli and using Node.Js.
 
This post is mainly to show that it's possible to recreate a nice local development setup with the other languages that GCP Cloud Functions supports, such as Java, Python, and Go.

> for the purpose of this write up, we will be using golang


# How to use

1. Make sure you have the firestore emulator installed locally on your machine using the `firebase-tools` cli

2. Fire up the firestore emulator `java -jar ~/.cache/firebase/emulators/cloud-firestore-emulator-v[version number installed].jar --functions_emulator localhost:6000`

3. With using the `github.com/GoogleCloudPlatform/functions-framework-go/funcframework` package replace your calls to `RegisterEventFunction` to something like 

```go
	event := fsemu.EmuResource{ProjectId: "dummyid", Address: "http://localhost:8080"}
	emuRegisters := []fsemu.EmuRegister{{
        // pass your function in that needs called upon the event firing
		TriggerFn:    WriteSkills,
        // type of firestore trigger
		TriggerType:  fsemu.FirestoreOnWrite,
        // path the firestore collection/document
		ResourcePath: "skills/{id}",
	}}

	if err := event.RegisterToEmu(emuRegisters); err != nil {
		panic(err)
	}

	if err := funcframework.Start("6000"); err != nil {
		panic(err)
	}

	
```

4. Fire up you go program with setting the `FIRESTORE_EMULATOR_HOST` env variable from the firestore emulator output. example `export FIRESTORE_EMULATOR_HOST=localhost:8080`




 

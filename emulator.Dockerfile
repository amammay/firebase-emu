FROM alpine
RUN apk --no-cache add curl openjdk11
RUN curl https://storage.googleapis.com/firebase-preview-drop/emulator/cloud-firestore-emulator-v1.11.5.jar --output firebaseemu.jar

CMD ["java", "-jar", "firebaseemu.jar"]


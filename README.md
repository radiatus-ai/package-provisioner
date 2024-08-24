+ service that consumes messages from pubsub
+ grab secrets from the db / secret store
+ grab package values from db
+ provision infra / app / agent
+ update db

+ optionally the db adds this info to the msg and the provisioner can't talk to the db.
+ diagraming state is important, so we can have a robust UX
+ websocket support in the api will be required for shared user sessions


```
gcloud components install pubsub-emulator
gcloud beta emulators pubsub start --project=your-project-id
$(gcloud beta emulators pubsub env-init)
```

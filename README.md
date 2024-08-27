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


```yaml
# this should be json schema as well, allowing each package to have a dag of composeable steps as well
plan:
    - source: pkg 
      executor: terraform 
    # - source: pkg-new 
    #   executor: opentofu 
    #   generated: true  
```
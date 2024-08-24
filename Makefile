test:
	go test ./...

start:
	go run cmd/deployer/main.go

push-message:
	curl -X POST -d '{"message": {"data": "ewogICJwcm9qZWN0X2lkIjogInByb2otMTIzIiwKICAicGFja2FnZV9pZCI6ICJwa2ctNDU2IiwKICAicGFja2FnZSI6IHsKICAgICJ0eXBlIjogImRlcGxveW1lbnQiLAogICAgInBhcmFtZXRlcl9kYXRhIjogewogICAgICAidmVyc2lvbiI6ICIxLjAuMCIsCiAgICAgICJlbnZpcm9ubWVudCI6ICJwcm9kdWN0aW9uIgogICAgfSwKICAgICJvdXRwdXRzIjogewogICAgICAidXJsIjogImh0dHBzOi8vZXhhbXBsZS5jb20vYXBwIiwKICAgICAgInN0YXR1cyI6ICJzdWNjZXNzIgogICAgfQogIH0sCiAgImNvbm5lY3RlZF9pbnB1dF9kYXRhIjogewogICAgImRhdGFiYXNlX3VybCI6ICJwb3N0Z3JlczovL3VzZXI6cGFzc3dvcmRAaG9zdDo1NDMyL2RibmFtZSIsCiAgICAiYXBpX2tleSI6ICJhYmNkZWYxMjM0NTYiCiAgfQp9", "id": "123"}}'  -H 'Content-Type: application/json' localhost:8080/push

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o dist/main ./cmd/deployer/main.go

upload:
	docker compose build provisioner-deploy && \
        docker tag cloud-canvas-provisioner-deploy:latest us-central1-docker.pkg.dev/rad-containers-hmed/cloud-canvas/provisioner:latest && \
				docker push us-central1-docker.pkg.dev/rad-containers-hmed/cloud-canvas/provisioner:latest && \
				gcloud run deploy provisioner \
					--image=us-central1-docker.pkg.dev/rad-containers-hmed/cloud-canvas/provisioner:latest \
					--execution-environment=gen2 \
					--region=us-central1 \
					--project=rad-dev-canvas-kwm6 \
					&& gcloud run services update-traffic provisioner --to-latest --region us-central1 --project=rad-dev-canvas-kwm6

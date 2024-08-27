package pubsub

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	// Added import for io
	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

type Executor interface {
	PostOutputToAPI(projectID string, packageID string, outputData map[string]interface{}, action models.DeployStatus) error
}

type Subscriber struct {
	cfg      *config.Config
	deployFn func(models.DeploymentMessage) error
	executor Executor // Changed from *Executor to Executor
}

func NewSubscriber(cfg *config.Config, deployFn func(models.DeploymentMessage) error, executor Executor) *Subscriber {
	return &Subscriber{
		cfg:      cfg,
		deployFn: deployFn,
		executor: executor,
	}
}

func (s *Subscriber) HandlePush(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received push request: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Method not allowed: %s", r.Method)
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		http.Error(w, "Error reading request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var pushRequest struct {
		Message struct {
			Data []byte `json:"data,omitempty"`
			ID   string `json:"id"`
			Ack  string `json:"ack_id,omitempty"`
		} `json:"message"`
		Subscription string `json:"subscription"`
	}

	if err := json.Unmarshal(body, &pushRequest); err != nil {
		log.Printf("Error unmarshaling push request: %v", err)
		http.Error(w, "Error processing message", http.StatusBadRequest)
		return
	}

	// Acknowledge the message immediately
	w.WriteHeader(http.StatusOK)

	// Process the message asynchronously
	go func() {
		log.Printf("Processing message ID: %s", pushRequest.Message.ID)
		log.Printf("Received message data: %.100s", string(pushRequest.Message.Data))

		var deploymentMsg models.DeploymentMessage
		if err := json.Unmarshal(pushRequest.Message.Data, &deploymentMsg); err != nil {
			log.Printf("Error unmarshaling deployment message: %v", err)
			return
		}

		// don't print the deploymentMsg, it has secrets
		// log.Printf("%s package: %+v", deploymentMsg.Action, deploymentMsg)
		log.Printf("%s package: %s", deploymentMsg.Action, deploymentMsg.PackageID)
		if err := s.deployFn(deploymentMsg); err != nil {
			log.Printf("Error deploying package: %v", err)
			errorDeployData := map[string]interface{}{
				"error": err.Error(),
			}
			if postErr := s.executor.PostOutputToAPI(deploymentMsg.ProjectID, deploymentMsg.PackageID, errorDeployData, models.Failed); postErr != nil {
				log.Printf("Failed to post error to API: %v", postErr)
			}
		}
	}()
}

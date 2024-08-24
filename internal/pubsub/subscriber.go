package pubsub

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/radiatus-ai/package-provisioner/internal/config"
	"github.com/radiatus-ai/package-provisioner/pkg/models"
)

type Subscriber struct {
	cfg      *config.Config
	deployFn func(models.DeploymentMessage) error
}

func NewSubscriber(cfg *config.Config, deployFn func(models.DeploymentMessage) error) *Subscriber {
	return &Subscriber{
		cfg:      cfg,
		deployFn: deployFn,
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

	log.Printf("Received raw body: %s", string(body))

	var pushRequest struct {
		Message struct {
			Data []byte `json:"data,omitempty"`
			ID   string `json:"id"`
		} `json:"message"`
	}

	if err := json.Unmarshal(body, &pushRequest); err != nil {
		log.Printf("Error unmarshaling push request: %v", err)
		http.Error(w, "Error processing message", http.StatusBadRequest)
		return
	}

	log.Printf("Received message ID: %s", pushRequest.Message.ID)
	log.Printf("Received message data: %s", string(pushRequest.Message.Data))

	var deploymentMsg models.DeploymentMessage
	if err := json.Unmarshal(pushRequest.Message.Data, &deploymentMsg); err != nil {
		log.Printf("Error unmarshaling deployment message: %v", err)
		http.Error(w, "Error processing message", http.StatusBadRequest)
		return
	}

	log.Printf("Deploying package: %+v", deploymentMsg)
	if err := s.deployFn(deploymentMsg); err != nil {
		log.Printf("Error deploying package: %v", err)
		http.Error(w, "Error processing message", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deployed package: %s", deploymentMsg.Package.Type)
	w.WriteHeader(http.StatusOK)
}

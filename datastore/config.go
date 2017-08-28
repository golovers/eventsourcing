package datastore

type Config struct {
	// Google project ID
	ProjectID string `envconfig:"GOOGLE_PROJECT_ID"`
	// Data connector namespace on datastore
	Namespace string `envconfig:"DATASTORE_NAMESPACE"`
}

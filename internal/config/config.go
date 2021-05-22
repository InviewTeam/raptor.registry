package config

type Settings struct {
	RegistryAddress     string `env:"REGISTRY_ADDRESS"`
	DatabaseAddress     string `env:"MONGO_ADDRESS"`
	DatabaseName        string `env:"MONGO_DB_NAME"`
	WorkersCollection   string `env:"WORKERS_COLLECTION"`
	AnalyzersCollection string `env:"ANALYZERS_COLLECTION"`
	ReportsCollection   string `env:"REPORTS_COLLECTION"`
	RMQAddress          string `env:"RMQ_ADDRESS"`
	WorkerQueue         string `env:"WORKER_QUEUE"`
	AnalyzerQueue       string `env:"ANALYZER_QUEUE"`
}

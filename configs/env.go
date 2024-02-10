package configs

func getEnv() PostgresConfig {
	envVars := PostgresConfig{
		Host:     "localhost",
		Port:     "5432",
		User:     "pilolo",
		Password: "sredev",
		DbName:   "task-management-cli",
		SSLMode:  "disable",
	}

	return envVars
}

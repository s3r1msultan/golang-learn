package initializers

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func InitLogger() {
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Warn("Failed to log to file, using default stderr")
	} else {
		log.SetOutput(file)
	}

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.DebugLevel)
}

func LogError(action string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["error"] = err.Error()
	fields["action"] = action
	log.WithFields(fields).Error("An error occurred")
}

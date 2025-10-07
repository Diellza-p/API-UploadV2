package configs

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file once during package initialization
	err := godotenv.Load()
	if err != nil {
		// Log warning instead of fatal - allow app to continue with defaults
		if Logger != nil {
			Logger.Warn("Could not load .env file, using environment variables or defaults", "error", err)
		}
	}
}

func EnvMongoURI() string {
	if uri := os.Getenv("MONGODB_URI"); uri != "" {
		return uri
	}
	return "mongodb://localhost:27017/synapp" // default fallback
}

func EnvNotificationServiceURL() string {
	if url := os.Getenv("NOTIFICATION_SERVICE_URL"); url != "" {
		return url
	}
	return "http://localhost:3007" // default fallback
}

func EnvMediaDir() string {
	if dir := os.Getenv("MEDIADIR"); dir != "" {
		return dir
	}
	return "/tmp/media" // default fallback
}

func EnvStreamDir() string {
	if dir := os.Getenv("STREAMDIR"); dir != "" {
		return dir
	}
	return "/tmp/stream" // default fallback
}

func RedisURL() string {
	if url := os.Getenv("REDIS_URL"); url != "" {
		// Remove redis:// prefix if present, keep only host:port
		if strings.HasPrefix(url, "redis://") {
			return strings.TrimPrefix(url, "redis://")
		}
		return url
	}
	return "localhost:6379" // default fallback
}

func NOTIFICATIONCHANNEL() string {
	if channel := os.Getenv("NOTIFICATION_CHANNEL"); channel != "" {
		return channel
	}
	return "notifications" // default fallback
}

func INITMEDIADIR() string {
	if dir := os.Getenv("INIT_MEDIA_DIR"); dir != "" {
		return dir
	}
	return "/tmp/media" // default fallback
}

func EnvDBHost() string {
	if host := os.Getenv("DB_HOST"); host != "" {
		return host
	}
	return "localhost" // default fallback
}

func EnvDBUser() string {
	if user := os.Getenv("DB_USER"); user != "" {
		return user
	}
	return "postgres" // default fallback
}

func EnvDBPassword() string {
	if password := os.Getenv("DB_PASSWORD"); password != "" {
		return password
	}
	return "password123" // default fallback
}

func EnvDBName() string {
	if name := os.Getenv("DB_NAME"); name != "" {
		return name
	}
	return "synapp" // default fallback
}

func EnvDBPort() string {
	if port := os.Getenv("DB_PORT"); port != "" {
		return port
	}
	return "5432" // default fallback
}

func FEMALEAVATAR() string {
	if avatar := os.Getenv("FEMALEAVATAR"); avatar != "" {
		return avatar
	}
	return "default_female.jpg" // default fallback
}

func MALEAVATAR() string {
	if avatar := os.Getenv("MALEAVATAR"); avatar != "" {
		return avatar
	}
	return "default_male.jpg" // default fallback
}

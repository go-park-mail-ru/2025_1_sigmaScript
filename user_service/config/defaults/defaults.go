package defaults

const (
	ListenerPort = ":8082"

	DBHost            = "127.0.0.1"
	DBPort            = 5433
	DBUser            = "filmlk_user"
	DBPassword        = "filmlk_password"
	DBName            = "filmlk"
	DBMaxOpenConns    = 10
	DBMaxIdleConns    = 5
	DBConnMaxLifetime = "300s"
	DBConnMaxIdleTime = "60s"

	StorageRelativePath = "/static/avatars/"
)

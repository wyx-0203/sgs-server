package global

const (
	PORT = "8080"
	// MYSQL_DSN = "sgs:123456@tcp(host.docker.internal:3306)/sgs?charset=utf8&parseTime=True&loc=Local"
	MYSQL_DSN = "sgs:123456@tcp(123.56.19.80:3306)/sgs?charset=utf8&parseTime=True&loc=Local"

	SSL_IS_ON   = true
	SSL_CRT     = "/app/acapp.pem"
	SSL_CRT_KEY = "/app/acapp.key"

	JWT_SECRET = "sgssgs"
)

module github.com/thanhpp/prom

go 1.16

replace (
	github.com/coreos/go-systemd => github.com/coreos/go-systemd/v22 v22.0.0
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	firebase.google.com/go v3.13.0+incompatible
	github.com/appleboy/go-fcm v0.1.5
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.7.1
	github.com/go-redis/redis/v8 v8.8.2
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.2.0
	github.com/spf13/viper v1.7.1
	github.com/streadway/amqp v1.0.0
	go.etcd.io/etcd v3.3.25+incompatible
	go.uber.org/zap v1.16.0
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/api v0.13.0
	google.golang.org/grpc v1.32.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.21.8
	sigs.k8s.io/yaml v1.2.0 // indirect
)

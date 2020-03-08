package mock

//go:generate mockgen -package=mock -destination=./kinesisapi.go github.com/aws/aws-sdk-go/service/kinesis/kinesisiface KinesisAPI

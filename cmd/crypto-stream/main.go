package main

import (
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/rs/zerolog/log"
)

const streamName = "xxx"

func main() {
	creds := credentials.NewStaticCredentials("xxx", "xxxx", "")
	sess := session.Must(session.NewSession(&aws.Config{
		CredentialsChainVerboseErrors: nil,
		Credentials:                   creds,
		Region:                        aws.String("eu-central-1"),
	}))
	kinesisClient := kinesis.New(sess)
	interruptChan := make(chan os.Signal, 1)
	stream := NewStream()
	publisher := NewPublisher(kinesisClient, streamName)

	err := stream.Init()
	if err != nil {
		log.Panic().Err(err).Msg("")
		panic(err.Error())
	}

	listen, err := stream.GetListener()
	if err != nil {
		log.Panic().Err(err).Msg("")
		panic(err.Error())
	}

	go func() {
		for {
			tr, err := listen()
			if err != nil {
				log.Error().Err(err).Msg("")
			}

			log.Info().Fields(map[string]interface{}{
				"value":    tr.GetUSDValue(),
				"happened": tr.GetDate(),
			}).Msg("Transaction")

			b, err := json.Marshal(tr)
			if err != nil {
				log.Error().Err(err).Msg("problem while marshaling")
				continue
			}

			err = publisher.Publish(b)
			if err != nil {
				log.Error().Err(err).Msg("problem sending to kinesis")
			}
		}
	}()

	select {
	case <-interruptChan:
		err := stream.Stop()
		if err != nil {
			log.Error().Err(err).Msg("")
		}
	}

}

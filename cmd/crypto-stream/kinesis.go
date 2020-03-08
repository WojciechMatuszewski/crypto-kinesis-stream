package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/aws/aws-sdk-go/service/kinesis/kinesisiface"
)

type Publisher struct {
	streamName string
	client     kinesisiface.KinesisAPI
}

func NewPublisher(client kinesisiface.KinesisAPI, streamName string) Publisher {
	return Publisher{
		streamName: streamName,
		client:     client,
	}
}

func (p Publisher) Publish(data []byte) error {
	_, err := p.client.PutRecord(&kinesis.PutRecordInput{
		Data:         data,
		PartitionKey: aws.String("1"),
		StreamName:   aws.String(p.streamName),
	})

	return err
}

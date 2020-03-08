package main

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kinesis"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"tweets-kinesis/cmd/crypto-stream/mock"
)

const mockStreamName = "STREAM"

func TestPublisher_Publish(t *testing.T) {

	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewMockKinesisAPI(ctrl)
		publisher := NewPublisher(client, mockStreamName)

		in := []byte("test")
		client.EXPECT().PutRecord(&kinesis.PutRecordInput{
			Data:         in,
			PartitionKey: aws.String("1"),
			StreamName:   aws.String(mockStreamName),
		}).Return(nil, nil)

		err := publisher.Publish(in)
		assert.NoError(t, err)
	})

	t.Run("failure", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		client := mock.NewMockKinesisAPI(ctrl)
		publisher := NewPublisher(client, mockStreamName)

		in := []byte("test")
		client.EXPECT().PutRecord(&kinesis.PutRecordInput{
			Data:         in,
			PartitionKey: aws.String("1"),
			StreamName:   aws.String(mockStreamName),
		}).Return(nil, errors.New("boom"))

		err := publisher.Publish(in)
		assert.Error(t, err)
	})
}

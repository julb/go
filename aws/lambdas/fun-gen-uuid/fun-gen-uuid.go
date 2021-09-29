package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/julb/go/pkg/dto"
	"github.com/julb/go/pkg/logging"
	"github.com/julb/go/pkg/util/identifier"
)

func HandleRequest(ctx context.Context) (*dto.ValueDTO, error) {
	uuid := identifier.Generate()

	logging.Debugf("uuid generated: %s", uuid)

	return &dto.ValueDTO{
		Value: uuid,
	}, nil
}

func main() {
	lambda.Start(HandleRequest)
}

package main

import (
	"bytes"
	"context"
	"fmt"
	htmlTemplate "html/template"
	"io/ioutil"
	"os"
	"strings"
	textTemplate "text/template"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-xray-sdk-go/xray"
	"github.com/julb/go/pkg/dto"
	log "github.com/julb/go/pkg/logging"
	"github.com/julb/go/pkg/xaws"
)

type GenerateNotificationContentEvent dto.GenerateNotificationContentDTO
type LambdaSettings struct {
	LogLevel     string
	BucketName   string `yaml:"bucketName" json:"bucketName"`
	BucketRegion string `yaml:"bucketRegion" json:"bucketRegion"`
}

func HandleRequest(ctx context.Context, event *GenerateNotificationContentEvent) (*dto.LargeContentWithSubjectDTO, error) {
	// get the lambda settings
	lambdaSettings := &LambdaSettings{
		LogLevel:     os.Getenv("J3_LOGGING_LEVEL"),
		BucketName:   os.Getenv("J3_LAMBDA_BUCKET_NAME"),
		BucketRegion: os.Getenv("J3_LAMBDA_BUCKET_REGION"),
	}

	// configure log level
	if lambdaSettings.LogLevel != "" {
		log.SetLevel(lambdaSettings.LogLevel)
	}

	// configure x-ray
	err := xray.Configure(xray.Config{})
	if err != nil {
		return nil, err
	}
	xray.SetLogger(xaws.NewXRayLoggerProxy())

	// Trace input
	log.Debugf("s3 bucket holding template is: %s (region: %s)", lambdaSettings.BucketName, lambdaSettings.BucketRegion)

	// Get the resolvable paths.
	templatePath := getTemplateResolvablePath(ctx, lambdaSettings, event)
	log.Debugf("path resolved for the template: %s", templatePath)

	// Finds the template in the repository
	templateContent, err := getTemplateContentWithPath(ctx, lambdaSettings, templatePath)
	if err != nil {
		return nil, err
	}
	log.Debugf("template content resolved - rendering.")

	// Render the template content
	return renderTemplate(ctx, lambdaSettings, event, string(templateContent))
}

// Get the template resolve path
func getTemplateResolvablePath(ctx context.Context, lambdaSettings *LambdaSettings, event *GenerateNotificationContentEvent) string {
	fileName := fmt.Sprintf("%s.%s", event.Name, event.Type.TemplatingMode.FileExtension)
	return strings.ToLower(strings.Join([]string{event.Trademark, event.Locale.String(), event.Type.Alias, fileName}, "/"))
}

// Get the resolvable paths valid for the given event
func getTemplateContentWithPath(ctx context.Context, lambdaSettings *LambdaSettings, path string) ([]byte, error) {
	// get a aws session
	sess := session.Must(session.NewSession())

	// get a s3 client
	s3Client := s3.New(sess, aws.NewConfig().WithRegion(lambdaSettings.BucketRegion))
	xray.AWS(s3Client.Client)

	// read the file from s3
	s3GetObjectResult, err := s3Client.GetObjectWithContext(ctx, &s3.GetObjectInput{
		Bucket: aws.String(lambdaSettings.BucketName),
		Key:    aws.String(path),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				log.Errorf("object with key %s does not exist in the bucket %s", path, lambdaSettings.BucketName)
			default:
				log.Errorf("error when getting content of %s in the bucket %s: %s", path, lambdaSettings.BucketName, err.Error())
			}
		}
		return nil, err
	}

	defer s3GetObjectResult.Body.Close()

	return ioutil.ReadAll(s3GetObjectResult.Body)
}

func renderTemplate(ctx context.Context, lambdaSettings *LambdaSettings, event *GenerateNotificationContentEvent, templateContent string) (*dto.LargeContentWithSubjectDTO, error) {
	// Parse the template and render
	var buf bytes.Buffer
	if event.Type.TemplatingMode == *dto.HtmlTemplatingMode {
		tpl, err := htmlTemplate.New("html").Parse(templateContent)
		if err != nil {
			return nil, err
		}

		// Execute the template
		err = tpl.Execute(&buf, event.Parameters)
		if err != nil {
			return nil, err
		}
	} else if event.Type.TemplatingMode == *dto.TextTemplatingMode {
		tpl, err := textTemplate.New("txt").Parse(templateContent)
		if err != nil {
			return nil, err
		}

		// Execute the template
		err = tpl.Execute(&buf, event.Parameters)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("unsupported templating mode: %s", event.Type.TemplatingMode.Alias)
	}

	if event.Type.HasSubject {
		// subject on first line
		subject, err := buf.ReadString('\n')
		if err != nil {
			return nil, err
		}
		subject = strings.TrimSpace(subject)

		// empty line
		_, err = buf.ReadString('\n')
		if err != nil {
			return nil, err
		}

		// content is the remaining
		content := buf.String()

		return &dto.LargeContentWithSubjectDTO{
			MimeType: event.Type.TemplatingMode.MimeType,
			Subject:  subject,
			Content:  content,
		}, nil
	} else {
		return &dto.LargeContentWithSubjectDTO{
			MimeType: event.Type.TemplatingMode.MimeType,
			Content:  buf.String(),
		}, nil
	}
}

func main() {
	lambda.Start(HandleRequest)
}

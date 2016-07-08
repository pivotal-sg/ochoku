package imageservice

import (
	"bytes"
	"fmt"
	"image"
	"log"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/micro/go-platform/config"
	"github.com/satori/go.uuid"
)

// S3Store will persist the image to S3, and return the URI
type S3Store struct {
	Config config.Config
}

// generateS3URI will return you a s3 link to the uploaded resource
// based on the bucket/region/filename
func generateS3URI(filename, region, bucket string) string {
	uri := url.URL{
		Scheme: "https",
		Host:   fmt.Sprintf("%s.s3-%s.amazonaws.com", bucket, region),
		Path:   filename,
	}
	return uri.String()
}

// SaveImage to a S3 bucket, as jpeg with default compression.  This
// will return the URI of said image.
func (s3store S3Store) SaveImage(img image.Image) (string, error) {
	var aws_id, aws_secret, aws_token, aws_region, s3_bucket string

	aws_id = s3store.Config.Get("aws", "id").String("")
	aws_secret = s3store.Config.Get("aws", "secret").String("")
	aws_token = s3store.Config.Get("aws", "token").String("")
	aws_region = s3store.Config.Get("aws", "region").String("ap-southeast-1")
	s3_bucket = s3store.Config.Get("aws", "bucket").String("ochoku")

	// id, secret, token
	creds := credentials.NewStaticCredentials(aws_id, aws_secret, aws_token)

	conf := &aws.Config{Region: aws.String(aws_region)}
	conf = conf.WithCredentials(creds)

	svc := s3.New(session.New(conf))

	// create file like for writing to s3
	imgBuf := &bytes.Buffer{}
	if err := encodeToJpeg(img, imgBuf); err != nil {
		return "", err
	}

	// setup the upload to s3
	filename := uuid.NewV4().String() + ".jpg"
	params := &s3.PutObjectInput{
		Bucket:      aws.String(s3_bucket),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(imgBuf.Bytes()),
		ContentType: aws.String("image/jpeg"),
	}

	// put alls the object!
	if resp, err := svc.PutObject(params); err != nil {
		log.Printf("Failed to put the object to s3: '%v'\n\n resp: '%s'", err, resp.String())
		return "", err
	}

	uri := generateS3URI(filename, aws_region, s3_bucket)

	return uri, nil
}

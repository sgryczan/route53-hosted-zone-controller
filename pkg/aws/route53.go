package aws

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/google/uuid"
	"k8s.io/utils/pointer"
)

var cfg aws.Config

func init() {
	var err error
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err = config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
}

func GetZoneDetailByName(zone string) (*route53.ListHostedZonesByNameOutput, error) {
	// Using the Config value, create the Route53 client
	svc := route53.NewFromConfig(cfg)

	ctx := context.Background()

	output, err := svc.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{
		DNSName: &zone,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	js, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, err
	}
	log.Printf("%s\n", js)

	return output, nil
}

func HostedZoneExists(zone string) (*bool, error) {

	output, err := GetZoneDetailByName(zone)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if len(output.HostedZones) == 0 {
		return pointer.Bool(false), nil
	}
	return pointer.Bool(true), err
}

func CreateHostedZone(zone string) error {
	svc := route53.NewFromConfig(cfg)
	ctx := context.Background()

	output, err := svc.CreateHostedZone(ctx, &route53.CreateHostedZoneInput{
		CallerReference: genUUID(),
		Name:            pointer.String(zone),
	})

	if err != nil {
		return err
	}

	js, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("%s\n", js)

	return nil
}

func DeleteHostedZone(id string) error {
	svc := route53.NewFromConfig(cfg)
	ctx := context.Background()

	output, err := svc.DeleteHostedZone(ctx, &route53.DeleteHostedZoneInput{
		Id: pointer.String(id),
	})

	if err != nil {
		return err
	}

	js, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return err
	}
	log.Printf("%s\n", js)

	return nil
}

func genUUID() *string {
	uuidRaw := uuid.New()
	uuid := strings.Replace(uuidRaw.String(), "-", "", -1)
	return &uuid
}

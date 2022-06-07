package aws

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"

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

func GetZoneByName(zone string) (*route53.ListHostedZonesByNameOutput, error) {
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

	printJSON(output)

	return output, nil
}

func GetZoneDetailByName(zone string) (*route53.GetHostedZoneOutput, error) {
	// Using the Config value, create the Route53 client
	svc := route53.NewFromConfig(cfg)
	ctx := context.Background()

	z, err := GetZoneByName(zone)
	if err != nil {
		return nil, err
	}
	zoneID := strings.Replace(*z.HostedZones[0].Id, "/hostedzone/", "", -1)

	output, err := svc.GetHostedZone(ctx, &route53.GetHostedZoneInput{
		Id: &zoneID,
	})

	if err != nil {
		log.Println(err)
		return nil, err
	}

	printJSON(output)

	return output, nil
}

func HostedZoneExists(zone string) (*bool, error) {

	output, err := GetZoneByName(zone)

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

	printJSON(output)

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

	printJSON(output)

	return nil
}

func DeleteZoneDelegation(delegationName string, nameservers []string, zoneID string, roleArnToAssume string) error {
	config, err := AssumeRoleArn(roleArnToAssume)
	if err != nil {
		return err
	}

	r53svc := route53.NewFromConfig(*config)

	nsList := []types.ResourceRecord{}
	for _, ns := range nameservers {
		nsList = append(nsList, types.ResourceRecord{Value: pointer.String(ns + ".")})
	}

	batch := &types.ChangeBatch{
		Changes: []types.Change{
			{
				Action: types.ChangeActionDelete,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name:            pointer.String(delegationName),
					Type:            types.RRTypeNs,
					TTL:             pointer.Int64(300),
					ResourceRecords: nsList,
				},
			},
		},
	}

	printJSON(batch)

	output, err := r53svc.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  batch,
		HostedZoneId: &zoneID,
	})

	printJSON(output)

	if err != nil {
		return err
	}

	return nil
}

func CreateZoneDelegation(delegationName string, nameservers []string, zoneID string, roleArnToAssume string) error {
	config, err := AssumeRoleArn(roleArnToAssume)
	if err != nil {
		return err
	}

	r53svc := route53.NewFromConfig(*config)

	nsList := []types.ResourceRecord{}
	for _, ns := range nameservers {
		nsList = append(nsList, types.ResourceRecord{Value: pointer.String(ns + ".")})
	}

	batch := &types.ChangeBatch{
		Changes: []types.Change{
			{
				Action: types.ChangeActionUpsert,
				ResourceRecordSet: &types.ResourceRecordSet{
					Name:            pointer.String(delegationName),
					Type:            types.RRTypeNs,
					TTL:             pointer.Int64(300),
					ResourceRecords: nsList,
				},
			},
		},
	}

	printJSON(batch)

	output, err := r53svc.ChangeResourceRecordSets(context.Background(), &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  batch,
		HostedZoneId: &zoneID,
	})

	printJSON(output)

	if err != nil {
		return err
	}

	return nil
}

func AssumeRoleArn(roleArnToAssume string) (*aws.Config, error) {
	// get credentials for role in target account
	stssvc := sts.NewFromConfig(cfg)
	ctx := context.Background()

	result, err := stssvc.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         &roleArnToAssume,
		RoleSessionName: pointer.String("route53-controller"),
	})

	if err != nil {
		return nil, err
	}

	config := &aws.Config{
		Region: "us-east-1",
		Credentials: credentials.NewStaticCredentialsProvider(
			*result.Credentials.AccessKeyId,
			*result.Credentials.SecretAccessKey,
			*result.Credentials.SessionToken,
		),
	}

	return config, nil
}

func GetNameServers(zone string) ([]string, error) {
	zoneDetail, err := GetZoneDetailByName(zone)
	if err != nil {
		return nil, err
	}

	return zoneDetail.DelegationSet.NameServers, nil
}

func genUUID() *string {
	uuidRaw := uuid.New()
	uuid := strings.Replace(uuidRaw.String(), "-", "", -1)
	return &uuid
}

func printJSON(v interface{}) {
	js, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Printf("error marshalling output: %s\n", err.Error())
		return
	}
	log.Printf("%s\n", js)
}

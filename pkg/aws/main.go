package aws

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func TxtFuncMap() template.FuncMap {
	f := sprig.TxtFuncMap()
	delete(f, "env")
	delete(f, "expandenv")

	// Add some extra functionality
	extra := template.FuncMap{
		"awsLookupvpc":     AwsLookupVpc,
		"awsLookupSubnets": AwsLookupSubnets,
		"awsRegion":        AwsRegion,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func AwsRegion() string {
	return "us-east-1"
}

func AwsLookupSubnets(tagsJoined string) string {
	return ""
}

func AwsLookupVpc(tagsJoined string) (string, error) {
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"),
	})
	tags := strings.Split(tagsJoined, ",")
	svc := ec2.New(sess)
	req := &ec2.DescribeVpcsInput{}
	req.Filters = BuildTagFilter(tags)
	resp, err := svc.DescribeVpcs(req)
	if err != nil {
		if ec2err, ok := err.(awserr.Error); ok && ec2err.Code() == "InvalidVpcID.NotFound" {
			resp = nil
		} else {
			log.Printf("Error on AwsLookupVpc: %s", err)
			return "", err
		}
	}

	if resp == nil {
		// Sometimes AWS just has consistency issues and doesn't see
		// our instance yet. Return an empty state.
		return "", nil
	}

	if len(resp.Vpcs) == 0 {
		return "", nil
	}

	vpc := resp.Vpcs[0]
	if vpc != nil {
		return *vpc.VpcId, nil
	}
	return "", nil
}

func BuildTagFilter(tags []string) []*ec2.Filter {
	filters := []*ec2.Filter{}

	for _, tag := range tags {
		parts := strings.SplitN(tag, "/", 2)
		if len(parts) != 2 {
			panic(fmt.Errorf("expected TAG/VALUE got %s", tag))
		}
		tagName := fmt.Sprintf("tag:%s", *aws.String(parts[0]))
		filters = append(filters, &ec2.Filter{
			Name:   &tagName,
			Values: []*string{aws.String(parts[1])},
		})
	}
	return filters

}

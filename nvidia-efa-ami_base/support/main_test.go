package main

import (
    "testing"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
)

func TestGenCaseBody(t *testing.T) {
	// Works with empty iid
	iid := ec2metadata.EC2InstanceIdentityDocument{}
	body, err := genCaseBody(iid)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Fatal("Body is empty")
	}

	// Works with populated iid
	iid.InstanceID = "i-234234"
	iid.Region = "us-west-2"
	iid.AvailabilityZone = "a"
	iid.InstanceType = "p4d.24xl"
	iid.AccountID = "2222"
	iid.ImageID = "1111"
	iid.KernelID = "3333"

	body, err = genCaseBody(iid)
	if err != nil {
		t.Fatal(err)
	}
	if len(body) == 0 {
		t.Fatal("Body is empty")
	}
	t.Log(body)
}
 

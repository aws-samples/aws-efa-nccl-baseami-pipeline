/**
 * Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
 * the Software, and to permit persons to whom the Software is furnished to do so.
 * 
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
 * FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
 * COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
 * IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
 * CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package main

/**
 * Support SDK Godocs: https://docs.aws.amazon.com/sdk-for-go/api/service/support
 * Support Go Src code: https://github.com/aws/aws-sdk-go/blob/37a82efacad413c32032d9e120bc84ae54162164/service/support/api.go#L1514
 */

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/support"
	"os"
	"path/filepath"
	"text/template"
)

var (
	CC_EMAIL_1         = "foo@foo.com"
	CC_EMAIL_2         = "baz@baz.com"
	CC_EMAILS          = []*string{&CC_EMAIL_1, &CC_EMAIL_2}
	CASE_SUBJECT = "GPU Faults/Errors Encountered | Hardware Cordon Requested"
	CASE_BODY_TEMPLATE = `
Please cordon off the following ec2 instance:
  * Instance Id: {{.InstanceID}}
  * Region: {{.Region}} 
  * AZ: {{.AvailabilityZone}}
  * Instance Type: {{.InstanceType}}
  * Account Id: {{.AccountID}}
  * Image Id: {{.ImageID}}
  * Kernel Id: {{.KernelID}}
   
Attached are Nvidia Bug report logs.
   
Refer to Failure Modes - Require Degrade of the underlying EC2 Host
in the support playbook for more info`
)

/**
 * Uploads NVidia logs to AWS Support.
 * Returns:
 *   error: Err in process of uploading the logs
 *   string: AttachmentSetId
 */
func uploadLogs(client *support.Support, nvidialogs []string) (string, error) {
	// Read in attachments into memory
	attchs := make([]*support.Attachment, len(nvidialogs))
	for i, fp := range nvidialogs {
		attch := new(support.Attachment)
		data, err := os.ReadFile(fp)
		if err != nil {
			return "", err
		}
		attch.SetData(data)
		attch.SetFileName(filepath.Base(fp))
		attchs[i] = attch

	}

	// Upload attachments to support
	ats := new(support.AddAttachmentsToSetInput)
	ats.SetAttachments(attchs)
	resp, err := client.AddAttachmentsToSet(ats)
	if err != nil {
		return "", err
	}
	return *resp.AttachmentSetId, nil
}

/**
 * Uses Ec2 Metadata API to enrich a support case body with instance details.
 * Returns:
 *   string: Case body text
 *   error: Errors generated in process of calling Ec2 API or text/template
 */
func genCaseBody(iid ec2metadata.EC2InstanceIdentityDocument) (string, error) {
	// Populate case body with metadata on Ec2 instance
	tmpl, err := template.New("body").Parse(CASE_BODY_TEMPLATE)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, iid); err != nil {
		return "", err
	}
	return body.String(), nil
}

/**
 * Creates a support case to request cordoning an ec2 instance
 * Returns:
 *   string: Created case's CaseId
 *   error: errors in process of creating the case
 */
func RequestNodeCordon(nvidialogs []string) (string, error) {
	mySession := session.Must(session.NewSession())

	/**
	 * We can get credentials 3 ways:
	 * https://github.com/aws/aws-sdk-go#configuring-credentials
	 * Initial credentials loaded from SDK's default credential chain. Such as
     * the environment, shared credentials (~/.aws/credentials), or EC2 Instance
     * Role. These credentials will be used to to make the STS Assume Role API.
	 */

	// Support's API is only available in us-east-1
	// See: https://docs.aws.amazon.com/general/latest/gr/awssupport.html
	client := support.New(mySession, aws.NewConfig().WithRegion("us-east-1"))

	atsId, err := uploadLogs(client, nvidialogs)
	if err != nil {
		return "", err
	}

	emdClient := ec2metadata.New(mySession)
	if !emdClient.Available() {
		return "", errors.New("Cannot connect to ec2 metadata service")
	}

	iid, err := emdClient.GetInstanceIdentityDocument()
	if err != nil {
		return "", err
	}

	body, err := genCaseBody(iid)
	if err != nil {
		return "", err
	}
	fmt.Println(body)

	// Populate case fields
	supportCase := new(support.CreateCaseInput)
	supportCase.SetCcEmailAddresses(CC_EMAILS)
	supportCase.SetSubject(CASE_SUBJECT)
	supportCase.SetCommunicationBody(body)
	supportCase.SetIssueType("technical")
	supportCase.SetLanguage("en")

	// A list of support case service codes & category codes can be found using the CLI:
	// $ aws support describe-services --region=us-east-1
	supportCase.SetServiceCode("amazon-elastic-compute-cloud-linux")
	supportCase.SetCategoryCode("instance-issue")

	// A list of suport case severity levels and associated code can be found using the CLI:
	// $ aws support describe-severity-levels
	// critical = 15 minute target (SLO) for first response by AWS Enterprise Support
	// urgent = 1 hour target (SLO) for first response by AWS Enterprise Support
	supportCase.SetSeverityCode("critical")
	supportCase.SetAttachmentSetId(atsId)

	if err := supportCase.Validate(); err != nil {
		return "", err
	}

	res, err := client.CreateCase(supportCase)
	if err != nil {
		return "", err
	}
	return *res.CaseId, nil

}

func main() {
	caseID, err := RequestNodeCordon([]string{"./nvidia-bug-report.log.gz"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Case Id: %s\n", caseID)
}

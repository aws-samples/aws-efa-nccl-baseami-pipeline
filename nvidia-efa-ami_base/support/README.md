To use the support API the EC2 instance will call on your behalf the aws support API.
Minimum permissions to create a case with the nvidia-bug-report attachement is below.
You can create a policy and add the policy to an existing IAM profile attached to the EC2 instance

````json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "support:AddAttachmentsToSet",
                "support:AddCommunicationToCase",
                "support:CreateCase",
                "support:PutCaseAttributes"
            ],
            "Resource": "*"
        }
    ]
}
````

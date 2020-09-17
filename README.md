## AWS EFA and NCCL Base AMI Build Pipeline
The base EFA/NCCL Base AMI can help you quickly get started with running distributed training workloads on AWS with our EFA enabled instances (p3dn, etc).
These scripts can be used as examples for both AL2 and Ubuntu 18.04 the following stack is installed 

NVIDIA Driver 450.xx
CUDA 11
cuDNN 8
NCCL 2.7.8
EFA latest driver
AWS-OFI-NCCL 
FSx kernel and client driver and utilities
Intel OneDNN

The docker build files can show an example implentation of the requirements to setup EFA/NCCL in container context.

CUDA 11
cuDNN 8
NCCL 2.7.8
EFA
AWS-OFI-NCCL


## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.


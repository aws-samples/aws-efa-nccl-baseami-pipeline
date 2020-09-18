## AWS EFA and NCCL Base AMI Build Pipeline
The base EFA/NCCL Base AMI can help you quickly get started with running distributed training workloads on AWS with our EFA enabled instances (p3dn, etc).
Included are sample buildspecs which you integrate with a CodeBuild/CodePipeline for automatic builds.
These scripts can be used as examples for both AL2 and Ubuntu 18.04 the following stack is installed 

- NVIDIA Driver 450.xx
- CUDA 11
- cuDNN 8
- NCCL 2.7.8
- EFA latest driver
- AWS-OFI-NCCL 
- FSx kernel and client driver and utilities
- Intel OneDNN
- Docker

The docker build files can show an example implentation of the requirements to setup EFA/NCCL in container context.

- CUDA 11
- cuDNN 8
- NCCL 2.7.8
- EFA
- AWS-OFI-NCCL

## Packer Instructions
In the `nvidia-efa-ami_base` dir you will find packer scripts for Amazon Linux 2 and Ubuntu 18.04. Generally you just need to modify the `variables:{}` json and execute the packer build
````json
"variables": {
    "region": "<region>",
    "flag": "al2-base",
    "subnet_id": "<subnetid>",
    "security_groupids": "<security_group_ids_list>",
    "build_ami": "`<lastest_base_ami_>",
    "efa_pkg": "aws-efa-installer-latest.tar.gz",
    "intel_mkl_version": "intel-mkl-2020.0-088",
    "cuda_version": "cuda-11-0",
    "cudnn_version": "libcudnn8"
  },
````  

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.


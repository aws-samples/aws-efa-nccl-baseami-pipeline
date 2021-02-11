## AWS EFA and NCCL Base AMI/Docker Build Pipeline
The base EFA/NCCL Base AMI can help you quickly get started with running distributed training workloads on AWS with our EFA enabled instances (p3dn, g4dn, and p4d)
Included are sample buildspecs which you integrate with a CodeBuild/CodePipeline for automatic builds.
These scripts can be used as examples for both AL2 and Ubuntu 18.04 the following stack is installed. The docker build file is an example implentation of the requirements to setup EFA/NCCL in a container context for ECS/Batch/EKS.

- NVIDIA Driver 460.xx
- CUDA 11.2
- NVIDIA Fabric Manager 460.xx (version locked to the nvidia driver)
- cuDNN 8
- NCCL 2.8.3
- EFA latest driver
- AWS-OFI-NCCL 
- FSx kernel and client driver and utilities
- Intel OneDNN
- NVIDIA runtime Docker

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
    "cuda_version": "cuda-drivers-455 cuda-toolkit-11-1",
    "cudnn_version": "libcudnn8"
  },
````  
After filling in the `variables` check that the packer script is validated.
````
packer validate nvidia-efa-ml-al2.yml
packer build nvidia-efa-ml-al2.yml
````

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.


## AWS EFA and NCCL Base AMI/Docker Build Pipeline
### Amazon Linux 2022 Status
included in nvidia-efa-ami_base/nvidia-efa-ml-al2022.yml is a beta for NVIDIA GPU support in AL2022
- [X] NVIDIA Driver install path
- [X] CUDA install path
- [X] nvidia-docker install path (working nvidia-container-toolkit)
- [ ] GDRCopy
- [ ] EFA/GDR support
- [ ] FSx support

The base EFA/NCCL Base AMI can help you quickly get started with running distributed training workloads on AWS with our EFA enabled instances (p3dn, g4dn, and p4d)
Included are sample buildspecs which you integrate with a CodeBuild/CodePipeline for automatic builds.
These scripts can be used as examples for both AL2 and Ubuntu 18.04 the following stack is installed. The docker build file is an example implentation of the requirements to setup EFA/NCCL in a container context for ECS/Batch/EKS.

- NVIDIA Driver 510.xx*
- CUDA 11.6
- NVIDIA Fabric Manager (version locked to the nvidia driver)
- cuDNN 8
- NCCL 2.12.7
- EFA latest driver
- AWS-OFI-NCCL 
- FSx kernel and client driver and utilities
- Intel OneDNN
- NVIDIA runtime Docker

*Note: With the 510+ driver - G2/P2 instances are not supported. If you need to work with p2/g2 instances the last supported driver is on the 470 branch.

## Packer Instructions
In the `nvidia-efa-ami_base` dir you will find packer scripts for Amazon Linux 2 and Ubuntu 18.04. Generally you just need to modify the `variables:{}` json and execute the packer build
````json
"variables": {
    "region": "us-east-1",
    "flag": "<flag>",
    "subnet_id": "<subnetid>",
    "security_groupids": "<security_group_id,security_group_id",
    "build_ami": "<buildami>",
    "efa_pkg": "aws-efa-installer-latest.tar.gz",
    "intel_mkl_version": "intel-mkl-2020.0-088",
    "nvidia_version": "510.47.03",
    "cuda_version": "cuda-toolkit-11-6 nvidia-gds-11-6",
    "cudnn_version": "libcudnn8",
    "nccl_version": "v2.12.7-1"
  },
````  
After filling in the `variables` check that the packer script is validated.
````
packer validate nvidia-efa-ml-al2.yml
packer build nvidia-efa-ml-al2.yml
````
## Accelerator Metrics/Error Handling in Cloudwatch
In this repo we also have an accelerator metrics and error handling (XIDs) custom metric which will push key metrics into cloudwatch. This is particularly useful in situations where you have an abstracted view of the underlying accelerator and unable to monitor metrics directly. 
For NVIDIA GPUS the following metrics are captured:
![dashboard](imgs/accelerator_dashboard.png?raw=true "Dashboard")
Accelerator kernel utilization
Memory utilization
Memory free
Memory used
SM clocks
Memory clocks
Total uncorrectable ECC Errors

The metric code is natively added to all AMIs built from this repo but you can use it directly in your AMIs as well. If interested you can extend this code to use your own metrics montitor as long as you follow this JSON schema:
````json
{
  "Id": 1,
  "AcceleratorName": "NVIDIA A100-SXM4-40GB",
  "AcceleratorDriver": "470.42.01",
  "Metrics": {
    "<Metric_Name>": <Metric_Value>,
    "<Metric_Name>": <Metric_Value>
  }
}
````
Furthermore we have added error handling specifically for NVIDIA GPUs in Cloudwatch Logs. A logstream is created which will lift ```NVRM: ...``` related messages in the syslog of the instance and push them to Cloudwatch.
![error log](imgs/example_error.png?raw=true "Example CW logs")

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.


#!/bin/bash

export USE_SYSTEM_NCCL=1
export USE_NCCL=1
export USE_DISTRIBUTED=1                # skip setting this if you want to enable OpenMPI backend
export USE_QNNPACK=1
export USE_PYTORCH_QNNPACK=1
export TORCH_CUDA_ARCH_LIST="5.3;6.2;7.2;7.5;8.0;8.6"

export PYTORCH_BUILD_VERSION=1.10.0  # without the leading 'v', e.g. 1.3.0 for PyTorch v1.3.0
export PYTORCH_BUILD_NUMBER=1
export NCCL_INCLUDE_DIRS="/opt/nccl/build/include"
export NCCL_LIBRARIES="/opt/nccl/build/lib"
export NCCL_ROOT_DIR="/opt/nccl/build"

git clone --recursive --branch v1.10.0 http://github.com/pytorch/pytorch
cd pytorch

apt-get install python3-pip cmake libopenblas-dev -y
pip3 install -r requirements.txt
pip3 install scikit-build ninja

python3 setup.py bdist_wheel

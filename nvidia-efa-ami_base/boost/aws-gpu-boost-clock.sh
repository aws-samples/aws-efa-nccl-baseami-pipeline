#!/bin/bash

if [ ! nvidia-smi ]; then
   echo "nvidia-smi not found"
   exit 0
fi

sudo nvidia-smi -pm 1
sudo nvidia-smi --auto-boost-default=0

GPUNAME=$(nvidia-smi -L | head -n1)
echo $GPUNAME

if [[ $GPUNAME == *"A100-SXM4-40GB"* ]]; then
   nvidia-smi -ac 1215,1410
elif [[ $GPUNAME == *"A100-SXM4-80GB"* ]]; then
   nvidia-smi -ac 1593,1410
elif [[ $GPUNAME == *"V100"* ]]; then
   nvidia-smi -ac 877,1530
elif [[ $GPUNAME == *"K80"* ]]; then
   nvidia-smi -ac 2505,875
elif [[ $GPUNAME == *"T4"* ]]; then
   nvidia-smi -ac 5001,1590
elif [[ $GPUNAME == *"M60"* ]]; then
   nvidia-smi -ac 2505,1177
else
   echo "unsupported gpu"
fi

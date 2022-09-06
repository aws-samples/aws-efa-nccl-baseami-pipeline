import boto3
import os
from datetime import datetime
from time import sleep
import sys
import subprocess
import json
import gzip
import urllib


BINARY_PATH=sys.argv[1]

#endpoint_url = "http://localhost:8000/"
### CHOOSE REGION ####
#EC2_REGION = 'us-east-1'

###CHOOSE NAMESPACE PARMETERS HERE###
my_NameSpace = 'AcceleratorMetrics' 

### CHOOSE PUSH INTERVAL ####
sleep_interval = 10

### CHOOSE STORAGE RESOLUTION (BETWEEN 1-60) ####
store_reso = 1

### INSTANCE INFO
#### Support IMDS-v2
def send_request(url, header, method):
    req = urllib.request.Request(url, headers=header, method=method)
    with urllib.request.urlopen(req) as f:
        response = f.read().decode('utf8')
    return response

HEADER = {'X-aws-ec2-metadata-token-ttl-seconds' : '21600'}
TOKEN = send_request('http://169.254.169.254/latest/api/token', HEADER, 'PUT')

HEADER = {'X-aws-ec2-metadata-token' : TOKEN}

BASE_URL = 'http://169.254.169.254/latest/meta-data/'
INSTANCE_ID = send_request(BASE_URL + 'instance-id', HEADER, 'GET')
INSTANCE_TYPE = send_request(BASE_URL + 'instance-type', HEADER, 'GET')
INSTANCE_AZ = send_request(BASE_URL + 'placement/availability-zone', HEADER, 'GET')
print(INSTANCE_AZ)
EC2_REGION = INSTANCE_AZ[:-1]

cloudwatch = boto3.client('cloudwatch', region_name=EC2_REGION)
#event_system = cloudwatch.meta.events

PUSH_TO_CW = True

def create_metric_shard(i,d,n,m):
    metric_shard=[]
    MY_DIMENSIONS=[
                    {
                        'Name': 'Id',
                        'Value': INSTANCE_ID
                    },
                    {
                        'Name': 'InstanceType',
                        'Value': INSTANCE_TYPE
                    },
                    {
                        'Name': 'AcceleratorIndex',
                        'Value': str(i)
                    },
                    {
                        'Name': 'AcceleratorName',
                        'Value': str(n)
                    },
                    {
                        'Name': 'AcceleratorDriver',
                        'Value': str(d)
                    }
                ]
    AGR_DIMENSIONS=[
                    {
                        'Name': 'Id',
                        'Value': INSTANCE_ID
                    },
                    {
                        'Name': 'InstanceType',
                        'Value': INSTANCE_TYPE
                    },
                    {
                        'Name': 'AcceleratorName',
                        'Value': str(n)
                    },
                    {
                        'Name': 'AcceleratorDriver',
                        'Value': str(d)
                    }
                ]
    for key, value in m.items():
        a={'MetricName':key,'Dimensions':MY_DIMENSIONS,'Unit':'None','StorageResolution': store_reso,'Value':int(value)}
        metric_shard.append(a)
    for key, value in m.items():
        a={'MetricName':key,'Dimensions':AGR_DIMENSIONS,'Unit':'None','StorageResolution': store_reso,'Value':int(value)}
        metric_shard.append(a)
    return metric_shard

def gzip_request_body(request, **kwargs):
    gzipped_body = gzip.compress(request.body)
    request.headers.add_header('Content-Encoding', 'gzip')
    request.data = gzipped_body

def logResults(metric_shard):
    if (PUSH_TO_CW):
#        event_system.register('before-sign.cloudwatch.PutMetricData', gzip_request_body)
        cloudwatch.put_metric_data(
            Namespace=my_NameSpace,
            MetricData=metric_shard
        )


def main():
        while True:
            PUSH_TO_CW = True
            accel_metric_list=[]
            accel_metric_shard=subprocess.check_output(BINARY_PATH, universal_newlines=True)
            accel_metric_list=accel_metric_shard.splitlines()
            for accel in range(len(accel_metric_list)):
                d=json.loads(accel_metric_list[accel])
                ametric_shard=create_metric_shard(d['Gpu_index'],d['Driver'],d['Gpu_name'],d['Metrics'])
                print(ametric_shard)
                logResults(ametric_shard)
            sleep(sleep_interval)

if __name__=='__main__':
    main()

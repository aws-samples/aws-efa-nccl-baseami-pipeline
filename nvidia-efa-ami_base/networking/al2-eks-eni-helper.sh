#!/usr/bin/env bash

get_metadata()
{
     TOKEN=$(curl -X PUT "http://169.254.169.254/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 60")
     attempts=60
     false
     while [ "${?}" -gt 0 ]; do
         if [ "${attempts}" -eq 0 ]; then
         echo "Failed to get metdata"
         exit 1
         fi
         meta=$(curl -H "X-aws-ec2-metadata-token: $TOKEN" http://169.254.169.254/latest/meta-data/${1})
         if [ "${?}" -gt 0 ]; then
             let attempts--
             sleep 0.5
             false
         fi
     done
     echo "$meta"
}

get_gateway() 
{
    gateway=$(sipcalc $1 | grep HostMin | awk '{print $2}')
    echo $gateway
}

PRIMARY_MAC=$(get_metadata 'mac')
PRIMARY_IF=$(ip -o link show | grep -F "link/ether $PRIMARY_MAC" | awk -F'[ :]+' '{print $2}')
ALL_MACS=$(get_metadata 'network/interfaces/macs')

MAC_ARRAY=($ALL_MACS)
TABLE_ID=1001
PREF_ID=32767
for MAC in "${MAC_ARRAY[@]}"; do
    TRIMMED_MAC=$(echo $MAC | sed 's:/*$::')
    IF_NAME=$(ip -o link show | grep -F "link/ether $TRIMMED_MAC" | awk -F'[ :]+' '{print $2}')

    if [ "$IF_NAME" = "$PRIMARY_IF" ]; then
    echo "Primary Interface so no need to setup"
    else
        IF_IP=$(ip --family inet address show $IF_NAME | grep -o 'inet [^/ ]*' | cut -f2 -d' ')
    echo $IF_IP
        CIDR=$(get_metadata "network/interfaces/macs/$TRIMMED_MAC/vpc-ipv4-cidr-block")
        echo $CIDR
        echo $TABLE_ID $PREF_ID
        GATEWAY_IP=$(get_gateway $CIDR)
        ip route add $CIDR dev $IF_NAME table $TABLE_ID proto kernel scope link src $IF_IP
        ip route add default via $GATEWAY_IP dev $IF_NAME table $TABLE_ID
        ip rule add from $IF_IP lookup $TABLE_ID pref $PREF_ID
        ((TABLE_ID=TABLE_ID+1))
        ((PREF_ID=PREF_ID+1))   
    fi
done

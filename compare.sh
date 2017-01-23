#!/bin/bash

# Want to compare speed with and without reggie

# Assume reggie is running on 5000 and normal registry on 6000
# normal registry should be exposed on mesh network

docker pull amouat/large-image-arm:100
docker tag amouat/large-image-arm:100 localhost:5000/large-image-arm:100
docker tag amouat/large-image-arm:100 localhost:6000/large-image-arm:100

start_reggie=$(date +%s%N)
docker service create --mode global localhost:5000/large_image-arm:100 compare_test
docker push localhost:5000/large-image-arm:100

# Want unique code to avoid picking up old versions
# Could use squash and different file to stop caching
docker service ls -f name=reggie | awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ exit 1} }'
stop_reggie=$(date +%s%N)

echo "Reggie took: " $(($stop_reggie-$start_reggie)) "ns"
docker service rm compare_test

start_classic=$(date +%s%N)
docker service create --mode global localhost:6000/large_image:unique_code compare_test
docker push localhost:6000/large-image-arm:100
docker service ls -f name=reggie | awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ exit 1} }'
stop_classic=$(date +%s%N)


echo "Classic took: " $(($stop_reggie-$start_reggie)) "ns"


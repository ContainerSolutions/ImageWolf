#!/bin/bash
set -e

# Want to compare speed with and without reggie

# Assume reggie is running on 5000 and normal registry on 6000
# normal registry should be exposed on mesh network

docker service rm compare_test_1 compare_test_2 || true
docker pull amouat/large-image-arm:100
docker tag amouat/large-image-arm:100 localhost:5000/large-image-arm:100
docker tag amouat/large-image-arm:100 localhost:6000/large-image-arm:100
#To avoid pinning to RepoID, need to use SHA
image_sha=$(docker inspect --format {{.Id}} localhost:5000/large-image-arm:100)

start_reggie=$(date +%s%N)
docker service create --name compare_test_1 --mode global $image_sha
docker push localhost:5000/large-image-arm:100

ready=$(docker service ls -f name=compare_test_1 | \
  awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')

while [[ $ready != "0" ]]; do
  ready=$(docker service ls -f name=compare_test_1 | \
    awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')
  # don't like this sleep in a timed loop....
  sleep 0.1
done

stop_reggie=$(date +%s%N)

echo "Reggie took: " $(($stop_reggie-$start_reggie)) "ns"
docker service rm compare_test_1

start_classic=$(date +%s%N)
docker service create --mode global --name compare_test_2 127.0.0.1:6000/large-image-arm:100 
docker push 127.0.0.1:6000/large-image-arm:100
ready=$(docker service ls -f name=compare_test_2 | \
  awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')

while [[ $ready != "0" ]]; do
  ready=$(docker service ls -f name=compare_test_2 | \
    awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')
  # don't like this sleep in a timed loop....
  sleep 0.1
done
stop_classic=$(date +%s%N)


echo "Classic took: " $(($stop_classic-$start_classic)) "ns"


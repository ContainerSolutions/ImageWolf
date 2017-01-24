#!/bin/bash
set -e

# Want to compare speed with and without reggie

# Assume reggie is running on 5000 and normal registry on 6000
# normal registry should be exposed on mesh network

image_name="large-image-arm:108"

docker service rm compare_test_1 compare_test_2 || true

#The big problem with the test is that this really needs to be a new
#image, not cached on hosts
#Try using a short lived service to get rid of all images
docker service rm image_rm || true
docker service create --name image_rm --mode global \
  --mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
  amouat/docker-arm \
  bash -c "docker rmi localhost:5000/$image_name 127.0.0.1:6000/$image_name $image_name || true; sleep infinity"

ready=$(docker service ls -f name=image_rm | \
  awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')

while [[ $ready != "0" ]]; do
  ready=$(docker service ls -f name=image_rm | \
    awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')
  # don't like this sleep in a timed loop....
  sleep 0.1
done

docker service rm image_rm

docker pull amouat/$image_name
#There's some weirdness about which interface to use here
#Something to do with networking/swarm/registry...
docker tag amouat/$image_name localhost:5000/$image_name
docker tag amouat/$image_name 127.0.0.1:6000/$image_name

#To avoid pinning to RepoID, need to use SHA
#Also ensures using latest version
image_sha=$(docker inspect --format {{.Id}} localhost:5000/$image_name)

start_reggie=$(date +%s%N)
docker service create --name compare_test_1 --mode global $image_sha
docker push localhost:5000/$image_name

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
docker service rm compare_test_1 || true
docker service rm image_rm || true
docker service create --name image_rm --mode global \
  --mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
  amouat/docker-arm \
  bash -c "docker rmi localhost:5000/$image_name 127.0.0.1:6000/$image_name $image_name || true; sleep infinity"

ready=$(docker service ls -f name=image_rm | \
  awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')

while [[ $ready != "0" ]]; do
  ready=$(docker service ls -f name=image_rm | \
    awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')
  # don't like this sleep in a timed loop....
  sleep 0.1
done

docker tag amouat/$image_name 127.0.0.1:6000/$image_name

start_classic=$(date +%s%N)
docker service create --mode global --name compare_test_2 127.0.0.1:6000/$image_name 
docker push 127.0.0.1:6000/$image_name
ready=$(docker service ls -f name=compare_test_2 | \
  awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')

while [[ $ready != "0" ]]; do
  ready=$(docker service ls -f name=compare_test_2 | \
    awk 'NR==2 {split($4,a,"/"); if (a[1] != a[2]){ print "1" } else { print "0" }}')
  # don't like this sleep in a timed loop....
  sleep 0.1
done
stop_classic=$(date +%s%N)
docker service rm compare_test_2 || true


echo "Classic took: " $(($stop_classic-$start_classic)) "ns"


[![Build Status](https://travis-ci.org/ArangoGutierrez/ImageWolf.svg?branch=master)](https://travis-ci.org/ArangoGutierrez/ImageWolf)
[![Go Report Card](https://goreportcard.com/badge/github.com/ArangoGutierrez/ImageWolf)](https://goreportcard.com/report/github.com/ArangoGutierrez/ImageWolf)

ImageWolf - Fast Distribution of Docker Images on Clusters
==========================================================

ImageWolf is a PoC that provides a blazingly fast way to get Docker images
loaded onto your cluster, allowing updates to be pushed out quicker.

ImageWolf works alongside existing registries such as the Docker Hub, Quay.io
as well as self-hosted registries.

The PoC for ImageWolf uses the BitTorrent protocol spread images around the
cluster as they are pushed.

## Video

### Docker Swarm

[![asciicast](https://asciinema.org/a/DowEjf7Inqhtu4ZQsvZfA2b0j.png)](https://asciinema.org/a/DowEjf7Inqhtu4ZQsvZfA2b0j)

### Kubernetes

[![asciicast](https://asciinema.org/a/01rQtDxr67y4Gtu85KpBJ9cz2.png)](https://asciinema.org/a/01rQtDxr67y4Gtu85KpBJ9cz2)


## Getting Started

ImageWolf is currently alpha software and intended as a PoC - please don't run it in
production!

### Docker Swarm Mode

To start ImageWolf, run the following on your Swarm master:

```
docker network create -d overlay --attachable wolf
docker service create --network wolf --name imagewolf --mode global \
       --mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
       containersol/imagewolf
```

The ImageWolf service is now running on all nodes in our cluster.

The next step is to link ImageWolf with a registry. Whenever an image is pushed to
the registry, ImageWolf will immediately pull it and distribute across all the
nodes. To set up a private registry linked to ImageWolf:


```
# First find the id of the ImageWolf task running on this node
# This should work, but is a bit of a hack
TASK=$(docker ps -f name="imagewolf." --format {{.ID}})

# Configuration for the notification endpoint

export REGISTRY_NOTIFICATIONS_ENDPOINTS=$(cat <<EOF
    - name: imagewolf
      disabled: false
      url: http://${TASK}:8000/registryNotifications
EOF
)

# Start up a single instance of the registry
docker run -d --name registry-wolf --network wolf -p 5000:5000 -p 5001:5001 \
           -e REGISTRY_NOTIFICATIONS_ENDPOINTS \
           containersol/registry-wolf
```


You can then push an image to the registry running on the local node:

```
docker tag redis localhost:5000/myimage
docker push localhost:5000/myimage
```

ImageWolf should immediately see the push and distribute the image to the other
nodes. You can see what's going on by running `docker service logs imagewolf`.

We can now start another global service using this image:

```
# Use the digest of the image to avoid problems with repo lookups
IMAGE_HASH=$(docker inspect --format {{.Id}} localhost:5000/myimage)
docker service create --name test-service --mode global $IMAGE_HASH
```

In order to monitor progress, you can either pass `-d=false` when starting the
service or run `docker service ps test-service`. Note that nodes will reject
jobs until ImageWolf completes loading the image onto the node.

### Kubernetes

In Kubernetes, there is the [DaemonSet](https://kubernetes.io/docs/concepts/workloads/controllers/daemonset/)
concept which ensures that all (or some) nodes run a copy of a pod. As nodes
are added to the cluster, pods are added to them. As nodes are removed from the
cluster, those pods are garbage collected.

Then a [headless Service](https://kubernetes.io/docs/concepts/services-networking/service/#headless-services)
is created allowing ImageWolf to discover all the peers.


To start ImageWolf, deploy the DaemonSet and its associated headless Service using
the following command:

```
kubectl apply -f kubernetes.yaml
```

Testing:

```
# You may need to wait few minutes before Kubernetes display the service public IP
export IMAGEWOLF_IP=$(kubectl get svc imagewolf --no-headers | awk '{print $3}')

# Simulate a Docker Hub webhook
curl "${IMAGEWOLF_IP}/hubNotifications" \
   -H 'Content-Type: application/json' \
   -d '{
      "push_data": {
         "tag": "latest"
      },
      "repository": {
         "repo_name": "redis"
      }
   }'

# Inspect the logs
kubectl logs -l app=imagewolf

# Check the stats
curl "${IMAGEWOLF_IP}/stats"
```

## Integration with Docker Hub

The Docker Hub has a web hooks feature which can be used to call a remote
service when an image is pushed. When ImageWolf receives the callback, it will
pull the image and distribute to the cluster, which is *significantly*
faster than all nodes pulling individually from the Docker Hub.

To use this feature, you will need to expose the ImageWolf service so that it is
accessible to the Hub. This can be done by adding the flag `-p 8000:8000` to the
`service create` command. You can then add the URL or IP address of your server
as a webhook, specifying hubNotifications as the path e.g:
http://mycluster.com/hubNotifications. If your cluster runs on a internal
network you can use a service such as ngrok to forward calls.

## Stats

There are no hard numbers yet.

The real improvements are expected on large clusters, where multiple Docker
engines pull images simultaneously. Also whilst a ramped deployment may avoid
the "stampeding herd" problem swamping the registry, deployment times will still
be longer as whenever a container is deployed to a node without the image a new
pull will take place - with ImageWolf the image will already be on the node and
the container will start immediately.

## Other Approaches

Using a global or distributed file system to back a Docker registry can also
achieve many of the benefits of ImageWolf.

## Multiarch

ImageWolf was tested on a Raspberry PI cluster as well as in the Google cloud. You
should find that the above instructions work identically on 32-bit ARM (armv7l)
as well as x86_64 through the magic of multi-arch images.

## Feedback

This is a PoC. If it is useful or interesting to you, please get in touch and
let us know.

 - adrian.mouat@container-solutions.com

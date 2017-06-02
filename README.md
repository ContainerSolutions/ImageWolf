Reggie - The Cluster-first Image Registry
=========================================

Reggie is a PoC into building a "cluster-first registry" for Docker images. The
purpose of reggie is to provide a blazingly fast way to get new images loaded
onto your cluster, allowing updates to be pushed out quicker.

Reggie is not intended to replace existing registry services such as the Docker
Hub or Quay.io. Instead, it works alongside such services. The centralised
service continues to provide a stable and reliable store for images over time,
whereas Reggie provides a cluster-local cache of the images.

The PoC for Reggie uses the BitTorrent protocol and the existing Docker registry
to spread images around the cluster as they are pushed.

## Video

## Getting Started

The PoC was developed for Docker Swarm Mode. If there is sufficient interest,
versions for Kubernetes and other cluster managers will follow. Reggie is
currently alpha software and intended as a PoC - please don't run it in
production!

To start Reggie, run:

```
docker network create -d overlay --attachable reggie
docker service create --network reggie --name reggie --mode global \
       --mount type=bind,src=/var/run/docker.sock,dst=/var/run/docker.sock \
       amouat/reggie
```

Now we have a Reggie instance running across all nodes in our cluster.

The next step is to link Reggie with a registry. Whenever an image is pushed to
the registry, Reggie will immediately pull it and distribute across all the
nodes. At the moment Reggie only works with the private Docker registry and will
distribute all pushed images. To set up a private registry linked to Reggie:


```
# First find the id of the Reggie task running on this node
# This should work, but is a bit of a hack
TASK=$(docker ps -f name="reggie." --format {{.ID}})

# Configuration for the notification endpoint

export REGISTRY_NOTIFICATIONS_ENDPOINTS=$(cat <<EOF
    - name: reggie
      disabled: false
      url: http://${TASK}:8000/registryNotifications
EOF
)

# Start up a single instance of the registry
docker run -d --name registry-reggie --network reggie -p 5000:5000 -p 5001:5001 \
           -e REGISTRY_NOTIFICATIONS_ENDPOINTS \
           amouat/registry-reggie
```


You can then push an image to the registry running on the local node:

```
docker tag alpine localhost:5000/myimage
docker push localhost:5000/myimage
```

Reggie should immediately see the push and distribute the image to the other
nodes. You can see what's going on by running `docker service logs reggie`.

We can now start another global service using this image:

```
TK
```

TK say something aobut speed, large images

## Integration with Docker Hub

The Docker Hub has a web hooks feature which can be used to call a remote
service when an image is pushed. When Reggie recieves the callback, it can then
pull the image and distribute to the cluster, which will be *significantly*
faster than all nodes pulling individually from the Docker Hub.

This isn't implemented yet, but it should be straightforward. 

## Stats

There are no hard numbers yet.

The real improvements are expected on large clusters, where multiple Docker
engines pull images simultaneuously. Also whilst ramped deployments may avoid
the "stampeding herd" problem swamping the reigstry, they also hugely extend the
time taken to deploy new versions as pulls are performed in serial (in reggie
the startup time of containers will be much faster as the image is already on
the node).

## Other Approaches

Using a global or distributed file system to back a Docker registry will also
achieve many of the benefits of Reggie. 

## Multiarch

Reggie was tested on a Raspberry PI cluster as well as in the Google cloud. You
should find that the above instructions work identically on 32-bit ARM (armv7l)
as well as x86_64 through the magic of multi-arch images.

## Bugs & Improvements

Reggie is a PoC currently and there are a lot of rough edges:

 - Services have to be started using the Image ID to avoid repo pinning problems
 - No optimisations have been carried out
 - The internal use of the Docker CLI and sock is a bit hacky
 - Swarm will attempt to pull images on nodes at the same time as Reggie is
   distributing the image

Assuming there is interest in Reggie, the next step will be to change the hacked
together code into a coherent solution.

## Feedback

This is a PoC. If it is useful or interesting to you, please get in touch and
let us know.

 - adrian.mouat@container-solutions.com

# TODO
## Bugs & Improvements

ImageWolf is a PoC currently and there are a lot of rough edges:

 - [ ] Services have to be started using the Image ID to avoid repo pinning problems
 - [ ] No optimisations have been carried out
 - [ ] The internal use of the Docker CLI and sock is a bit hacky
 - [ ] If ImageWolf is still distributing the image when a service is created, nodes
   will attempt to pull from the registry simultaneous with ImageWolf pushing
   the image
 - [ ] Allow Google Cloud Platform container registry webhook using [Pub/Sub](https://cloud.google.com/container-registry/docs/configuring-notifications)
 - [ ] Build a go-test pkg for CI
 - [x] Replace `exec.Command` with the docker API "github.com/docker/docker/client"
      `client.NewEnvClient()` + client.ImagePull

Assuming there is interest in ImageWolf, the next step will be to change the hacked
together code into a coherent solution.

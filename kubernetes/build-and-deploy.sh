set DOCKER_ACCOUNT=<your-docker-account>
make
make manifests
make docker-buildx PLATFORMS=linux/amd64 IMG=<your-docker-account>/container-as-a-service:<version>
make docker-push IMG=<your-docker-account>/container-as-a-service:<version>
make install
make deploy IMG=<your-docker-account>/container-as-a-service:<version>
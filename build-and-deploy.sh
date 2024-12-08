set DOCKER_ACCOUNT=nightdeath
make
make manifests
make docker-buildx PLATFORMS=linux/amd64 IMG=nightdeath/container-as-a-service:1.7
make docker-push IMG=nightdeath/container-as-a-service:1.7
make install
make deploy IMG=nightdeath/container-as-a-service:1.7
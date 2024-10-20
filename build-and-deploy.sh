set DOCKER_ACCOUNT=nightdeath
make
make manifests
make docker-build IMG=nightdeath/container-as-a-service:latest
make docker-push IMG=nightdeath/container-as-a-service:latest
make install
make deploy IMG=nightdeath/container-as-a-service:latest
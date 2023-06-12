# telepresence-compose-demo

## To build & push the image

```cli
VERSION=0.1.3
docker build -f docker/userapi/Dockerfile --platform=linux/amd64 -t knlambert/telepresence-compose-demo-userapi:${VERSION} .
docker push knlambert/telepresence-compose-demo-userapi:${VERSION}
```

```cli
VERSION=0.1.5
docker build -f docker/contactapi/Dockerfile --platform=linux/amd64 -t knlambert/telepresence-compose-demo-contactapi:${VERSION} .
docker push knlambert/telepresence-compose-demo-contactapi:${VERSION}
```

## Deploy in cluster

```cli
kubectl apply -f infra[userapi.go](cmd%2Fuserapi.go)
```

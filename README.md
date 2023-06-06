# telepresence-compose-demo

## To build & push the image

```cli
VERSION=0.1.3
docker build . --platform=linux/amd64 -t knlambert/telepresence-compose-demo-userapi:${VERSION}
docker push knlambert/telepresence-compose-demo-userapi:${VERSION}
```

## Deploy in cluster

```cli
kubectl apply -f infra
```

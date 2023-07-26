
## DEPLOY TO DOCKERHUB

```bash
podman build -t amixp/nginx-log-generator-v2:latest .
podman login -u "amixp" -p "**********" docker.io
podman push amixp/nginx-log-generator-v2:latest
```

On destination:

`docker pull docker.io/amixp/nginx-log-generator-v2:latest`

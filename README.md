# Kitewave

![Tests](https://github.com/jh-chee/kitewave/actions/workflows/test.yml/badge.svg)

Sample Description

----

## Install
### Docker
Builds both HTTP and RPC images, and runs `docker compose up`.
```
make compose
```

### Helm
Build both images first. Charts are inside the [deployment](deployment/kitewave) folder
```
make docker-build
cd deployment
helm install kitewave kitewave
```
---
## Todo
- Add gh action to push images to docker hub
- Add HPA

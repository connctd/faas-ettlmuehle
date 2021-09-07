# FAAS OpenFaas

This function generates dummy waterlevel values depending on day of year and daytime.

## Requirements

You need to download the golang http template first

`faas template pull https://github.com/openfaas-incubator/golang-http-template`

## How to build and deploy

`make build` generates build directory and build docker image

`make push` pushes docker image to registry

`make install` installs necessary secret, cronjob and instructs openfaas gateway to deploy service. Please note: right now we are NOT exposing the openfaas gw to the outside world. To let the faas-cli communicate with it we need to do a port forwarding via kubectl like so: `kubectl port-forward GATEWAY-POD 8080 -n openfaas`

`make uninstall` removes secrets and deployment


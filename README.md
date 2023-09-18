# Manifests for two different environment test and prod

Manifests are built using Kustomize. Changes between environments are configured using patches.

The small test application exposes an endpoint that returns a value saved in a ConfigMap. When the value changes in the ConfigMap, the endpoint should return the new value.

## Installation and run

`.env` files are not commited to VCS so they need to be created in the local environment:

```shell
cat > ./k8s/overlays/prod/.env <<EOF
MSG=Hello testers of the PROD environment!
EOF

cat > ./k8s/overlays/test/.env <<EOF
MSG=Hello testers of the TEST environment!
EOF
```

Deploy the 2 environments:

```shell
kubectl apply -k ./k8s/overlays/test
kubectl apply -k ./k8s/overlays/prod
```

Test locally using port forward:

```shell
kubectl port-forward service/test-helloworld 6666:80
```

and in a separate shell:

```shell
curl -s http://127.0.0.1:6666/
```

## The application

If a basic Go application that exposes an HTTP endpoint. It has a test so the github action can validate the code. I chose to read the value to be exposed from the environment since this is what most applications do.

## ConfigMap as environment variables

I chose to map the ConfigMap as environment variables since this is the most common pattern. Because of this, changing values in the ConfigMap requires a restart of the application in order for it to even notice the new values.

Sice Kustomize deploys configMaps with different names each time (and automatically adjusts the deployment to match the new name), then simply applying a new changed configMap will trigger a rollout. One problem with this approach is that old configmaps are not cleaned up.

## ConfigMap as volume

If the ConfigMap is mapped as a volume, then the application wouldn't even need a restart. The volume is automatically updated by Kubernetes and our application would only need to read the value from a file instead of environment. Of course the value in the file would be cached so a monitor would have to run in order to notice the change and invalidate the cache. If that is not desired, still a restart is required.

## Monitor the configmap

Using Operator-sdk we can define our own custom resource and write a kubernetes operator that would react to any changes to the configmap. It could then identify the applications that use the custom resource (new configmap) and trigger a rollout.
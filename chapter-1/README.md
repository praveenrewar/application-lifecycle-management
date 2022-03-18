# Creating templated manifests and deploying them safely

The `/app` directory defines a simple web-app which has a response which is configured by the environment variable `USER_NAME`. We have built a docker image for the same so that we can easily deploy it on top of Kubernetes.
We can find the image [here](https://hub.docker.com/r/prewar/simple-server/).

To start off, lets have a look at our starter config in `/starter-config`.
```bash
$ cd starter-config
$ ls
config.yaml
```

`config.yaml` defines a manifest which creates a Deployment using the image of our `simple-server` and a Service which allows us to send requests to it.

```yaml
---
apiVersion: v1
kind: Service
metadata:
  namespace: default
  name: simple-server-app
spec:
  ports:
  - port: 8081
    targetPort: 8081
  selector:
    simple-app: ""
---
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: simple-server-app
spec:
  replicas: 3
  selector:
    matchLabels:
      simple-app: ""
  template:
    metadata:
      labels:
        simple-app: ""
    spec:
      containers:
      - name: simple-server-app
        image: docker.io/prewar/simple-server:latest
        env:
        - name: USER_NAME
          value: "John Doe"

```
While defining the container which is running the image we set the value of environment variable `USER_NAME` to "John Doe".

## Deploying resources using `kapp`
To deploy our resources onto the cluster, we will use `kapp`
```bash
kapp deploy -a simple-server -f config.yaml
```
After confirming the changes, `kapp` waits for our deployment to reach it's desired state. We group our resources into an abstract app called 'simple-server'.

We should be able to see the resources created by running,
```bash
kapp inspect -a simple-server -t
```

Let's have a look at how our application looks like. Let's forward the port exposed by the service to our system.
```bash
kubectl port-forward svc/simple-server-app 8081:8081
```

We should be able to see "Hello, John Doe", if we open up `localhost:8081` in our browsers.

Let's see how `kapp` deals with updates to the cluster. You can use a text editor of your choice to change the value of `spec.replicas` to 2.

Now we can update the app by deploying it again,
```bash
kapp deploy -a simple-server -f config.yaml -c
```
`kapp` show the changes to the user, asks for confirmation and then waits for the resource to reach it's desired state.

## Templating our manifest
We will now work towards two goals:
- Allow some values to be configurable via templating and 
- Reduce redundancy in our manifest

Import the "data" module at the top of the config file.
```yaml
#@ load("@ytt:data", "data")
```

We can now template parts of the manifest, like ports exposed and targeted and the name of the user. We can do something like this to achieve this:
```yaml
#...
spec:
  ports:
  - port: #@ data.values.svc_port
    targetPort: #@ data.values.app_port
#...
#...
        - name: USER_NAME
          value: #@ data.values.user_name
```

To define descriptions and default values for our supplied values we can create a file `values-schema.yaml` with the following contents.

```yaml
#@data/values-schema
---
#@schema/desc "Port number for the service."
svc_port: 8081
#@schema/desc "Target port for the application."
app_port: 8081
#@schema/desc "Name used in hello message from app when app is pinged."
user_name: stranger
```
To see how the the YAML is templated using the default values we can run,
```bash
ytt template -f .
```

We are adding labels at a number of places in the manifest, we can define a function which adds these to reduce redundancy.
```yaml
#@ def labels():
simple-app: ""
#@ end
```
We can now call this function to add the map of values wherever we require labels.
```yaml
#...
  selector:
    simple-app: ""
#...
  selector:
    matchLabels: #@ labels()
#...
  template:
    metadata:
      labels: #@ labels()
#...
```

Running,
```bash
ytt template -f .
```
will yield the templated result.

Now if we were to provide custom values to override the defaults. We would add another file `values.yml`
```
#@data/values
---
user_name: soumik
```

We can now pipe the templated manifest from `ytt` to `kapp` and deploy the templated resources.

```bash
ytt template -f . | kapp deploy -a simple-server -f - --yes
```

We can now use `kubectl port-forward` to see that the changes have taken affect✨

At the end of this chapter, the files in `starter-config` should look something like `/config` (excluding `values.yml` that we used to provid custom input)

In chapter-2, we will replicate this inner loop workflow declaratively.

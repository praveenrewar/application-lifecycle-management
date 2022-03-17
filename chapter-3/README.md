## Ensuring immutability of artifacts and manifests

Before we package and version our workflows so that they can be consumed. We need to ensure that our manifests and images or artifacts consumed by the manifests are immutable.

This is to prevent inconsistencies in multiple installations of the same package. If 2 users were to install version `1.0.0` of `hello.corp.com` their clusters should have the exact same changes.

We generally refer to images using tags, something like:
```yaml
image: index.docker.io/prewar/simple-server:latest
```
However, a different image can be pushed with the same tag and the reference is does not immutable.

We could alternatively use a digest reference,
```yaml
image: index.docker.io/prewar/simple-server@sha256:452f0d74bd18c7110022815e6c8eeafeaada6f93eeb030d3efcd5c0df3eadcbd
```
This would refer to one and only one image. And is thus immutable. We would want our manifests to use immutable references to ensure consistency.

In this chapter we will:
- See how kbld resolves images to their digests
- Bundle our manifests into an OCI image and push it to a registry so that it can be used by a package

## Resolving images to their digest form using `kbld`
To start off lets have a look at the `/starter-config`. It consists of what we had at the end of chapter-1.
```bash
$ cd starter-config
$ ls
config.yml        values-schema.yml
```
We would want to first template using `ytt` and then convert the image references to their digest form.
```bash
$ ytt template -f . | kbld -f .
```
We can see that the output has all the images resolved to their digest form. We could pipe this output to `kapp` in order to deploy manifests with resolved images.

However, since we intend to bundle the manifest into a bundle using `imgpkg`. We will be generate a lock file which points `imgpkg` to the resolved references.

Lets create a folder for our lockfile.
```bash
$ mkdir .imgpkg
```
We can now generate the lockfile by running
```bash
ytt template -f . | kbld -f - --imgpkg-lock-output .imgpkg/images.yml 
```
If we take a look at the generated file it looks something like this,
```yaml
---
apiVersion: imgpkg.carvel.dev/v1alpha1
images:
- annotations:
    kbld.carvel.dev/id: docker.io/prewar/simple-server:latest
    kbld.carvel.dev/origins: |
      - resolved:
          tag: latest
          url: docker.io/prewar/simple-server:latest
      - preresolved:
          url: index.docker.io/prewar/simple-server@sha256:452f0d74bd18c7110022815e6c8eeafeaada6f93eeb030d3efcd5c0df3eadcbd
  image: index.docker.io/prewar/simple-server@sha256:452f0d74bd18c7110022815e6c8eeafeaada6f93eeb030d3efcd5c0df3eadcbd
kind: ImagesLock
```

## Bundling our configuration using `imgpkg`
Now that `imgpkg` knows where it can find the resolved images, we can bundle our manifests and push it to a registry.

Make sure that you run `docker login` to authenticate into Deocker Hub before continuing with the next steps.

The command to push the image that has our manmifests bundled will look something like:
```bash
$ imgpkg push -b <docker-hub-username>/hello-app -f .
```
The end result is something like this.
```bash
imgpkg push -b 100mik/hello-app -f .
dir: .
dir: .imgpkg
file: .imgpkg/images.yml
file: config.yml
file: values-schema.yml
Pushed 'index.docker.io/100mik/hello-app@sha256:03f4c4ae65a33a15db62036936af0bf8fa74fe465bdeb1aa1f1ffafdd508dc0e'
Succeeded
```

Now we can use the immutable image which has our configuration with references to resolved image as a source for out Package.

We define the workflow a Package installation creates on the server similar to how we defined a workflow in our App CR. An `imgpkg` bundle is one of the sources we can "fetch" our resources from. So our Package will be pointing towards the bundled image with our manifests which is immutable instead of a mutable git repository.

Let's dive into how we can author, version and distribute packages in chapter-4!
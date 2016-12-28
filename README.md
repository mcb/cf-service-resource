# Cloud Foundry Service Resource

An output only resource that will create/bind service to a
Cloud Foundry Application.  
based on  
* [concourse/cf-resource](https://github.com/concourse/cf-resource)  
* [mcb/cf-docker-resource](https://github.com/mcb/cf-docker-resource)

## Source Configuration

* `api`: *Required.* The address of the Cloud Controller in the Cloud Foundry
  deployment.
* `username`: *Required.* The username used to authenticate.
* `password`: *Required.* The password used to authenticate.
* `organization`: *Required.* The organization to push the application to.
* `space`: *Required.* The space to push the application to.
* `skip_cert_check`: *Optional.* Check the validity of the CF SSL cert.
  Defaults to `false`.

## Behaviour

### `out`: Create and Bind Service to a Cloud Foundry Application

1. Create and Bind service to the Cloud Foundry deployed application.
2. Restage Application

#### Parameters

* `service`: *Required.* Service name.
* `plan`: *Required.* Plan name of the service.
* `instance_name`: *Required.* Service instance name.
* `manifest`: *(Either) Required.* Path to an application manifest file.
* `current_app_name`: *(Either) Required.* The name of the application to bind service.  
When both are listed, `manifest` is used.
* `delete`: *Optional.* Default `false`. (not yet implemented)

## Pipeline example

```yaml
---
resource_types:
  - name: cf-service-resource
    type: docker-image
    source:
      repository: mcb/cf-service-resource

resources:
  - name: resource-web-app
    type: git
    source:
      uri: https://github.com/idahobean/cf-service-resource-test.git

  - name: foobar-cf
    type: cf
    source:
      api: https://api.foo.bar.cfapps.io
      username: USERNAME
      password: PASSWORD
      organization: ORG
      space: SPACE
      skip_cert_check: false

  - name: foobar-cf-service
    type: cf-service-resource
    source:
      api: https://api.foo.bar.cfapps.io
      username: USERNAME
      password: PASSWORD
      organization: ORG
      space: SPACE
      skip_cert_check: false

jobs:
- name: job-deploy-app
  public: true
  serial: true
  plan:
  - get: resource-web-app
    task: build
    file: resource-web-app/build.yml
  - put: foobar-cf
    params:
      manifest: build-output/manifest.yml
    on_success:
      put: foobar-cf-service
      params:
        manifest: build-output/manifest.yml
        service: p-mysql
        plan: 512mb
        instance_name: mysql-test

```

## Manifest Example

```yaml
---
applications:
- name: foobar-app
```

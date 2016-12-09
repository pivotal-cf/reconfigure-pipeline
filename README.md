# reconfigure-pipeline

This is a [concourse](https://concourse.ci) fly wrapper that fetches secrets for your pipeline without ever storing them on disk. 

## Example Usage

Let's say that you combine three types of secret notes in your pipeline: `Server` storing your AWS key pair (`my-aws-keys`), `SSH Key` with private key fetching your git repo (`repo-deploy-key`) and a freeform `Generic` with flat YAML for miscellaneous credentials (`misc-ci-creds`).

`reconfigure-pipeline -t ci -c mypipeline.yml` will understand the `lpass://` notation, fetch credentials and produce a YAML consumable by `fly`.

```
# mypipeline.yml
resources:
- name: golang
  type: docker-image
  source:
    repository: golang
    tag: latest

resource_types:
- name: terraform
  type: docker-image
  source:
    repository: ljfranklin/terraform-resource

resources:
  - name: terraform
    type: terraform
    source:
      storage:
        bucket: mybucket
        bucket_path: terraform-ci/
        access_key_id: lpass:///my-aws-keys/Username
        secret_access_key: lpass:///my-aws-keys/Password

  - name: my-ci-repo
    type: git
    source:
      uri: git@github.com:oozie/private-repo
      branch: master
      private_key: lpass:///repo-deploy-key/Private Key

jobs:
- name: do-my-thing
  public: true
  serial: true
  plan:
  - get: my-ci-repo
    trigger: true
  - task: do-my-thing
    params:
      datadog_api_key: lpass:///misc-ci-creds/Notes#datadog-api-key
      pivnet_api_key: lpass:///misc-ci-creds/Notes#pivnet-api-key
```

## Installation

```
go get github.com/pivotal-cf/reconfigure-pipeline
```

## Features & Limitations:
* At the moment `reconfigure-pipeline` can fetch credentials from any store as long as it's LastPass.
* The wrapper infers pipeline name from name of the YAML file. Pipeline name can be overriden with the `-p` option.

# Intro
fkit is a tool that will enable you to create templates to pass in to `eksctl`. This allows you to parameterize your declarative templates and use them for multiple environments, regions, or accounts. In addition to variables files, parameters can be looked up in your AWS account dynamically using a Lookup function in your template. This is useful if you want to look up your subnets according to a tagging standard.

It is written in Go and uses Go templating for a familiar interface. 

# Guide
## Creating a template file
This is an example of a simple `eksctl` declaration file that has been parameterized:
```yaml
#eksctl.tpl

apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: {{ .metadata.name }}
  region: {{ .metadata.region }}

vpc:
  subnets:
    private:
      {{ call .Lookup.ec2.Vpcs "Tags:environment" "Value:dev" }}
      
managedNodeGroups:
  - name: managed-ng-1
    minSize: 2
    maxSize: 4
    desiredCapacity: 3
    volumeSize: 20

```

Then, create a simple YAML file to describe your variables:
```yaml
#dev.yaml
metadata:
  name: my-cluster
  region: us-east-1
```

And run fkit:
```shell
./fkit --file eksclt.tpl --vars dev.yaml
```
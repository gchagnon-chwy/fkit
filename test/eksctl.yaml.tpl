apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: {{ .name }}
  region: {{ .region }}

vpc:
  subnets:
    private:


managedNodeGroups:
  - name: managed-ng-1
    minSize: 2
    maxSize: 4
    desiredCapacity: 3
    volumeSize: 20
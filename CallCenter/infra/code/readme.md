## Create Pods

1. kubectl apply -f c:\MyPrograms\learn-go\CallCenter\backend\src\services\teams\deployments\pods.yaml -n cloud-explorer-d

2. kubectl port-forward -n cloud-explorer-d pod/premteam 8080:8080

## create replicasets

1. kubectl create -f replicasets.yaml -n cloud-explorer-d --validate=false

2. kubectl delete pods premteam-replicaset-8qmcc -n cloud-explorer-d

## Run Mongo-Shell

- in a separate pod shell: `kubectl run -it -n cloud-explorer-d mongo-shell --image=mongo:4.0.17 --rm -- /bin/bash`
- in existing pod shell: `kubectl exec -it -n cloud-explorer-d mongodb-0  -- /bin/bash`

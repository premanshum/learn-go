## Create Pods

1. kubectl apply -f c:\MyPrograms\learn-go\CallCenter\backend\src\services\teams\deployments\pods.yaml -n cloud-explorer-d

2. kubectl port-forward -n cloud-explorer-d pod/premteam 8080:8080

## create replicasets

1. kubectl create -f replicasets.yaml -n cloud-explorer-d --validate=false

2. kubectl delete pods premteam-replicaset-8qmcc -n cloud-explorer-d

## Run Mongo-Shell

- in a separate pod shell: `kubectl run -it -n cloud-explorer-d mongo-shell --image=mongo:4.0.17 --rm -- /bin/bash`
- in existing pod shell: `kubectl exec -it -n cloud-explorer-d mongodb-0  -- /bin/bash`

Inside shell, connect to db using:

`mongosh mongodb-0.mongo-svc`

initiate the mongodb instances:
`rs.initiate({ _id: "rs0", members: [ { _id: 0, host:"mongodb-0.mongo-svc:27017"}, { _id: 1, host:"mongodb-1.mongo-svc:27017"}, { _id: 2, host:"mongodb-2.mongo-svc:27017"}]});`

in secondary, use below command:
`db.getMongo().setReadPref("primaryPreferred")`

To run the service locally  : `kubectl port-forward -n cloud-explorer-d svc/premteams-service 8080:5002`
To run the database locally : `kubectl port-forward -n cloud-explorer-d svc/mongo-svc 27017:27017`

kubectl -n cloud-explorer-d get pods
kubectl -n cloud-explorer-d get svc
kubectl -n cloud-explorer-d get pvc 
kubectl -n cloud-explorer-d get statefulset
kubectl -n cloud-explorer-d get all
kubectl -n cloud-explorer-d get pods --watch

kubectl -n cloud-explorer-d apply -f C:\MyPrograms\learn-go\CallCenter\infra\mongo\statefulset.yaml
kubectl -n cloud-explorer-d apply -f C:\MyPrograms\learn-go\CallCenter\infra\mongo\service.yaml

kubectl -n cloud-explorer-d delete pvc -l app=mongodb
kubectl -n cloud-explorer-d delete statefulset mongodb
kubectl -n cloud-explorer-d delete svc mongo-headless
kubectl -n cloud-explorer-d delete svc mongo-svc
kubectl -n cloud-explorer-d get replicaset premteam-replicaset -o yaml

# Run a mongoshell
kubectl run -it -n cloud-explorer-d mongo-shell --image=mongo --rm -- /bin/bash

# Port forwarding
kubectl port-forward -n cloud-explorer-d svc/premteams-service 8080:8080
kubectl port-forward -n cloud-explorer-d pod/ccteams-587bbccdb9-6zmvt 8080:8080
kubectl port-forward -n cloud-explorer-d svc/mongo-svc 27017:27017

kubectl exec -it -n cloud-explorer-d ccteams-5c9c5db54-5b7tp  -- /bin/sh

docker run --rm -it -p 7070:7070/tcp premanshu/team:2
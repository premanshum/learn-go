kubectl -n cloud-explorer-d get pods 
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

# Run a mongoshell
kubectl run -it -n cloud-explorer-d mongo-shell --image=mongo --rm -- /bin/bash



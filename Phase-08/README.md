## Prepare the secrets

### Copy the example secrets
```bash
cp api-secret.yaml.example api-secret.yaml
cp redis-secret.yaml.example redis-secret.yaml
cp repo-secret.yaml.example repo-secret.yaml
```

### Set your ghcr.io credentials
```bash
vim repo-secret.yaml
```

### Apply the whole directory
```bash
kubectl apply -f .
```

### Monitor
```bash
kubectl get all -n traceroute
```
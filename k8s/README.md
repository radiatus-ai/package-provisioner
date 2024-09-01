```
helm install external-secrets \
   external-secrets/external-secrets \
    -n external-secrets \
    --create-namespace \
    --set installCRDs=true
```

```
cat sa.yaml | kubectl apply -f -
cat service.yaml | k apply -f -
cat secret-store.yaml | kubectl apply -f -
cat secrets.yaml | kubectl apply -f -
cat deployment.yaml | kubectl apply -f -
```

# reptile

# Run with Docker
```
docker build . --tag chess
docker run --publish 8080:8080 chess
```

# Run without Docker
```
go run cmd/main.go
```
# Deploy 

Update frontend from `/app` with:
```
npm run build
```

Deploy from root directory with: 
```
fly deploy
```


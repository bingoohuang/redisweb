# go-redis-web
redis web admin based on go lang

# build
1. `go get -u github.com/jteeuwen/go-bindata/...`
2. `go get golang.org/x/tools/cmd/goimports`
3. `/gobin.sh & go build` 
5. build for linux :`env GOOS=linux GOARCH=amd64 go build -o go-redis-web.linux.bin`

# startup
```
bingoo@bingodeMacBook-Pro ~/G/go-redis-web> ./go-redis-web -h
Usage of ./go-redis-web:
  -contextPath string
    	context path
  -devMode
    	devMode(disable js/css minify)
  -port int
    	Port to serve. (default 8269)
  -servers string
    	servers list, eg: Server1=localhost:6379,Server2=password2/localhost:6388/0 (default "default=localhost:6379")

```



# snapshots
![image](https://user-images.githubusercontent.com/1940588/30140520-d5e9c8da-93a7-11e7-8b79-09cc3c24ed26.png)
![image](https://user-images.githubusercontent.com/1940588/30140593-45752924-93a8-11e7-8afc-033198aa13c1.png)
![image](https://user-images.githubusercontent.com/1940588/30140608-67b17132-93a8-11e7-8034-085e6f1ded26.png)
![image](https://user-images.githubusercontent.com/1940588/30140617-7977a8b4-93a8-11e7-955a-fe639d86b41b.png)
![image](https://user-images.githubusercontent.com/1940588/30140624-8b8e3b30-93a8-11e7-98fe-e09e79b91498.png)
![image](https://user-images.githubusercontent.com/1940588/30140641-a8b0c386-93a8-11e7-8d30-77a99eda6bfb.png)

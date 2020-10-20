# redisweb

redis web ui based on golang.

## build

1. `go get github.com/markbates/pkger/cmd/pkger && pkger`
1. `go build` or `go install ./...`
1. linux:`env GOOS=linux GOARCH=amd64 go build -o redisweb`
1. windows: `env GOOS=windows GOARCH=amd64 go build`

## startup

```bash
$ redisweb -h
Usage of redisweb:
  -config string
        config file path (default "redisweb.toml")

```

## snapshots
![image](https://user-images.githubusercontent.com/1940588/30140520-d5e9c8da-93a7-11e7-8b79-09cc3c24ed26.png)
![image](https://user-images.githubusercontent.com/1940588/30140593-45752924-93a8-11e7-8afc-033198aa13c1.png)
![image](https://user-images.githubusercontent.com/1940588/30140608-67b17132-93a8-11e7-8034-085e6f1ded26.png)
![image](https://user-images.githubusercontent.com/1940588/30140617-7977a8b4-93a8-11e7-955a-fe639d86b41b.png)
![image](https://user-images.githubusercontent.com/1940588/30140624-8b8e3b30-93a8-11e7-98fe-e09e79b91498.png)
![image](https://user-images.githubusercontent.com/1940588/30140641-a8b0c386-93a8-11e7-8d30-77a99eda6bfb.png)
![image](https://user-images.githubusercontent.com/1940588/30145525-68e4e82e-93c4-11e7-902b-18911786b05f.png)
![image](https://user-images.githubusercontent.com/1940588/30526969-cfb90608-9be8-11e7-8c78-e346a5a7c949.png)

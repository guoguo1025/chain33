[![pipeline status](https://api.travis-ci.org/bityuan/bityuan.svg?branch=master)](https://travis-ci.org/bityuan/bityuan/)
[![Go Report Card](https://goreportcard.com/badge/github.com/bityuan/bityuan)](https://goreportcard.com/report/github.com/bityuan/bityuan)

# 基于 chain33 区块链开发 框架 开发的 TIX公有链系统

#### 编译

```
git clone https://github.com/guoguo1025/tixchain.git $GOPATH/src/github.com/guoguo1025/tixchain
cd $GOPATH/src/github.com/guoguo1025/tixchain
go build -i -o tix
go build -i -o tix-cli github.com/guoguo1025/tixchain/cli
```

#### 运行

拷贝编译好的tix, tix-cli, tixChain.toml这三个文件置于同一个文件夹下，执行：

```
./tix -f tixChain.toml
```



# Mirror

同步 github issues 内容，并输出为 hugo 可构建的博客内容


## mirror 配置文件

mirror 使用 yaml 格式的配置文件， 参考 [conf.yaml](./conf.yaml)
mirror 查找运行目录下的mirror.yaml 作为配置文件

具体的配置内容有

- token: github token, 需要有 issue 读取权限
- login: github username
- repo:  要抓取的issue 所在的 repo 名称
- output:  转换成 hugo 博客内容后的本地输出目录，建议指定为 hugo 博客的 content 目录

## mirror 规则

### draft

打开的Issue 视为博客草稿，不抓取。只抓取已经关闭的 Issue 内容。


### content

issue 本身有 body 和 comment。 新建issue时的内容是body， 之后所有项目本人和其它 githuber 的追加的内容都是评论。

Mirror 支持将 Issue body 和 comment 都抓取下来余与body部分拼接合并作为完整博客内容。

由于 Issue 都是公开的，可能会被别人评论，作者还可能回复提问者。
为了区分这种情况，mirror 只拼接 comment 作者是项目作者本人， 且在 comment 开头有一句 html `<!--->` 注释语句的comment。


## Run

```shell
cp conf.yaml mirror.yaml
```
然后修改 mirror.yaml 中的内容

```shell 
go run app/main.go
```
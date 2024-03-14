保存到文件 ~/.ssh/config, 1080 为 socks 代理的端口
会自动代理配置了 ssh 格式的 git 地址.

```shell
Host github.com
ProxyCommand nc -X 5 -x localhost:1080 %%h %%p
``` 

如果需要将 https 格式转为 ssh 格式.
将下列代码追加到  ~/.gitconfig 中

```shell 
[url "ssh://git@github.com"]
    insteadOf = https://github.com
``` 

命令行代理
```shell
export http_proxy="%s" && export https_proxy="%s"
```


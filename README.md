# Auto Cert
> 自动获取、更新证书并部署至CDN
## 命令行
### 获取证书

```
cert obtain -r, --renewal: 自动更新
            -t, --renewal_time: 指定自动更新时间
            --cdn: 同时部署至cdn
            -h, --host: 获取证书的域名
            --config, -c: 指定配置文件路径
```

### 手动更新证书
```
cert renew -t, --renewal_time: 指定更新时间
           -h, --host: 更新证书的域名
           -r, --renewal: 自动更新
           --config, -c: 指定配置文件路径
```

### 查看自动更新
```
cert renew list
```

### 开启服务器
```
cert server start
```

### 关闭服务器
```
cert server stop
```

### 重启服务器
```
cert server restart
```
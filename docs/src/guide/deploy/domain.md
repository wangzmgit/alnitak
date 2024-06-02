# 域名配置

:::tip 前置条件
域名配置前，请确保已经成功安装了Nginx
:::

### 配置Nginx

在`/etc/nginx/conf.d/`目录新建`alnitak.conf`文件，内容如下：
```
server {
    listen       9030;
	server_name  localhost; 
	client_max_body_size 1024M;

    # 转发用户端
    location / {
		proxy_pass http://127.0.0.1:9010;
		proxy_set_header   Host             $host;
     	proxy_set_header   X-Real-IP        $remote_addr;						
        proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
		proxy_http_version 1.1;
    	proxy_set_header Upgrade $http_upgrade;
    	proxy_set_header Connection "upgrade";
	}

    #后台管理
    location /admin/ {
        root /usr/share/nginx/html;
        index index.html index.htm;
        try_files $uri $uri/ @admin;
    }

    # 解决后台管理history路由问题
    location @admin {
        rewrite ^.*$ /admin/index.html;
    }

    # 转发后端
    location /api/ {
		proxy_pass http://127.0.0.1:9000;
		proxy_set_header   Host             $host;
     	proxy_set_header   X-Real-IP        $remote_addr;						
        proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
		proxy_http_version 1.1;
    	proxy_set_header Upgrade $http_upgrade;
    	proxy_set_header Connection "upgrade";
	}
}
```

### 重启Nginx
使用以下命令重启nginx

```sh
nginx -s reload
```
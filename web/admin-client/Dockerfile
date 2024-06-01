FROM node:18-alpine AS builder

WORKDIR /web
COPY . .

# 设置npm镜像
RUN npm config set registry https://registry.npmmirror.com 
# 打包管理端网页
RUN yarn install
RUN yarn build 

# nginx镜像
FROM nginx:1.23.3-alpine as production-stage

# 移除nginx容器的default.conf文件、nginx配置文件
RUN rm /etc/nginx/conf.d/default.conf \
    && rm /etc/nginx/nginx.conf \
    && mkdir /etc/nginx/logs \
    && touch /etc/nginx/logs/error.log

# 把主机的nginx.conf文件复制到nginx容器的/etc/nginx文件夹下
COPY ./nginx.conf /etc/nginx/

# 拷贝前端vue项目打包后生成的文件到nginx下运行
COPY --from=builder /web/admin /usr/share/nginx/html/admin

# 暴露端口
EXPOSE 9030

# 使用daemon off的方式将nginx运行在前台保证镜像不至于退出
CMD ["nginx", "-g", "daemon off;"]
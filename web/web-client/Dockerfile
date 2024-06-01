FROM node:18-alpine AS builder

WORKDIR /web
COPY . .

# 设置npm镜像
RUN npm config set registry https://registry.npmmirror.com 
# 打包用户端网页
RUN yarn install
RUN yarn build 

FROM node:18-alpine

COPY --from=builder /web/.output /app/web/.output

WORKDIR /app/web

# 设置环境变量
ENV NITRO_PORT=9010

EXPOSE 9010

ENV ENV=docker
CMD ["node", "./.output/server/index.mjs"]
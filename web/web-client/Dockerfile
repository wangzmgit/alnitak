FROM node:18-alpine

WORKDIR /web
COPY . .

# 设置环境变量
ENV NITRO_PORT=9010

EXPOSE 9010

ENV ENV=docker
CMD ["node", "./.output/server/index.mjs"]
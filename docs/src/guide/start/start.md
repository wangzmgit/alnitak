# 项目启动

## 启动后端项目
::: tip 相关信息
在启动之前，需先创建一个数据库，其字符集设置为 utf8mb4。
系统在启动时会自动进行表的创建，不需要导入完整的表结构。
:::

### 运行项目
1. 进入到`server`目录
2. 使命令行输入`go run cmd\main.go -env dev`启动项目
3. 访问`localhost:9000`，如果页面中出现`404 not found`则说明项目运行成功


## 启动前端项目
::: tip 相关信息
前端需要使用yarn，如果尚未安装yarn，可以使用命令 `npm install -g yarn` 来安装yarn。
:::

### 运行用户端
1. 进入到`web/web-client`目录
2. 首次运行需要执行`yarn install` 安装项目依赖
3. 使命令行输入`yarn dev`启动项目

### 运行管理端
1. 进入到`web/admin-client`目录
2. 首次运行需要执行`yarn install` 安装项目依赖
3. 使命令行输入`yarn dev`启动项目
# 视频资源接口

## 修改资源标题

#### 请求URL
- `/api/v1/resource/modifyTitle `
  
#### 请求方式
- PUT 

####  请求头
- `Authorization': token`
- `"content-type": "application/json",`

#### 参数

| 参数名 | 必选 | 类型   | 说明     |
| :----- | :--- | :----- | -------- |
| id     | 是   | int    | 资源ID   |
| title  | 是   | string | 资源标题 |

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```

#### 备注
无

<!-- ************************ 分隔符 ************************ -->

## 删除资源

#### 请求URL
- `/api/v1/resource/deleteResource/资源ID `
  
#### 请求方式
- DELETE 

####  请求头
- `Authorization': token`

#### 返回示例 

``` json
{
  "code": 200,
  "data": null,
  "msg":"ok"
}
```
    
#### 备注
无
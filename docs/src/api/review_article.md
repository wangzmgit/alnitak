# 文章审核相关接口

## 文章审核通过

#### 请求URL
- `/api/v1/review/reviewArticleApproved `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`
  
#### 参数

| 参数名 | 必选 | 类型 | 说明   |
| :----- | :--- | :--- | ------ |
| aid    | 是   | int  | 文章ID |

### 返回示例 

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

## 文章审核不通过

#### 请求URL
- `/api/v1/review/reviewArticleFailed `
  
#### 请求方式
- POST 

####  请求头
- `Authorization': token`
- `"content-type": "application/json"`
  
#### 参数

| 参数名 | 必选 | 类型   | 说明       |
| :----- | :--- | :----- | ---------- |
| aid    | 是   | int    | 文章ID     |
| status | 是   | int    | 审核状态码 |
| remark | 否   | string | 问题备注   |

### 返回示例 

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

## 获取文章审核记录

#### 请求URL
- `/api/v1/review/getArticleReviewRecord?aid=文章ID `
  
#### 请求方式
- GET 

####  请求头
- `Authorization': token`

### 返回示例 

``` json
{
  "code": 200,
  "data": {
    "review": {
      "status": 2000,
      "remark": "",
      "createdAt": "",
    }
  },
  "msg":"ok"
}
```

#### 返回参数说明 
| 参数名   | 类型   | 说明                                |
| :------- | :----- | ----------------------------------- |
| status | int    | 审核状态 |
| remark     | string | 审核备注                            |
| createdAt  | time   | 审核时间                       |

#### 备注
无
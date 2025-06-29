const title = "弹幕网站标题";
const https = false;
const domain = "localhost:9000";
const mobile = "/mobile/";
const icp = "icp备案信息";
const security = "公网安备信息";
const keywords = "视频,弹幕";
const description = "这里是介绍";
const article = true; // 是否开启专栏模块

//上传文件大小限制，需要先修改后端大小限制
const maxImgSize = 5;//上传图片最大大小(单位M)
const maxVideoSize = 1024;//上传视频最大大小(单位M)

export const globalConfig = {
    title: title,
    keywords: keywords,
    description: description,
    https: https,
    domain: domain,
    mobile: mobile,
    icp: icp,
    security: security,
    maxImgSize: maxImgSize,
    maxVideoSize: maxVideoSize,
    article: article,
}
-- 授权 root 用户可以远程链接
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'AlnitakRoot@123';
flush privileges;
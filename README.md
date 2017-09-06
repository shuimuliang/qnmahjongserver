# 测试环境搭建


## mysql
- 安装 docker pull mysql
- 启动 docker run --name first-mysql -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 -d mysql
- 数据初始化 执行db/sql下的脚本

## qnmahjong
- 编译 GOOS=linux GOARCH=amd64 go build -o qnmahjong -ldflags "-X main.branch=`git rev-parse --abbrev-ref HEAD` -X main.commit=`git rev-parse HEAD`"
- 配置 修改conf_pro.toml配置文件
- 安装 docker build -t qnmahjong . 
- 启动 login : docker run -it --rm --name qnmahjong_login --link first-mysql:mysql -p 5001:5001 qnmahjong login
- 启动 logic : docker run -it --rm --name qnmahjong_logic --link first-mysql:mysql -p 5002:5002 qnmahjong logic

## 客户端登录
-- 示例登录地址 http://100.100.33.99:5001

---
## 联系方式/商业咨询:
- ![二维码](http://oupthc6v2.bkt.clouddn.com/qr.jpg?imageView2/2/w/200/h/200/format/jpg/q/75|imageslim)
- QQ技术支持群: 99124448

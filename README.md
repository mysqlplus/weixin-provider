# weixin-provider

发送实例：

tos=UserId1,UserId2,UserId3
	UserId为微信企业公众号通讯录员工账号
	消息接收者，多个接收者用','分隔，最多支持1000个
	指定为@all，则向关注该企业应用的全部成员发送

curl http://*.*.*.*:9000/weixin -d "tos=ops,dev&content=weixin Test"
### Golang Web 开发 脚手架封装


#### 一：业务分层设计

1：业务分层：controller、logic、dao 三层

2：用户请求参数是基于请求参数结构体映射的，使用了gin中的validator对参数进行了简单的校验

3：查询和保存用户都用到了User表的结构体，目的是在写入User数据或查询User数据的时候，可以用到

4：定义了返回状态码，和通用的返回接口

5：jwt模块认证# jiaoshoujia

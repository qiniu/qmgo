# mongox

基于 MongoDB 官方 mongo-go-driver（https://github.com/mongodb/mongo-go-driver）封装的一套数据库操作API。

Go 版本要求：Go 1.10及以上。
    
MongoDB 版本要求： MongoDB 2.6及以上。


## base.go

通用方法定义。


## client.go

数据库客户端连接方法定义。

提供 Client 客户端定义，Config 配置定义，建立数据库连接的接口等。


## collection.go

集合相关方法定义。

提供对集合记录进行增删改查、集合索引设置等。


## cursor.go

游标相关操作方法定义。

提供对游标遍历的相关操作。


## query.go

查询相关方法定义。

提供查询条件设置，如 sort、limit、skip、count、distinct 等，并返回一条或多条结果，亦或一个Cursor。


## interface.go

接口定义。

CursorI、CollectionI、QueryI。


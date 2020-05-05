#mq

1. 安装环境
	go 1.13 及以上
	docs/docker/docker-compose.yml  依赖环境

2. 配置环境
	conf/conf.yaml
		rabbitmq:
		  addr: 192.168.56.101
		  port: 5672
		  username: guest
		  password: guest0725
		grpc:
		  port: :50051
		mongodb:
		  addr: 192.168.56.101
		  port: 27017
		  username: admin
		  password: admin0725
		  database: news
		cron:
		  notifySpec: "*/1 * * * * ?"
		  sendSpec: "*/1 * * * * ?"
		  lockSpec: "*/5 * * * * ?"
		etcd: //暂不实现
		
3. 系统流程图：
	docs/消息事务系统流程图.png
	
4. 系统测试 

	1.插入mongodb conf 表
		{
			app_id:        "123",
			app_secret:    "test",
			servce_name:   "kkk",
			queue_name:    "demo",
			routing_key:   "demo",
			exchange_name: "demo",
			exchange_type: "topic",
			status:       1,
			updated_at:    1588665472,
			created_at:    1588665472,
		}
		
	2. 开启mq消费
	  go run consumer.go -conf  配置文件绝对路径  ;  例如： D:/item/go/mq/conf/conf.yaml
	
	3. 开启业务生产数据
	  http: go run test/http/http.go
	  grpc: 
			go run grpc/example_answer/server.go
			go run grpc/example_receve/client.go
	
	4. lock 表数据
		{
    		key: "send",
			value: 0,
			updated_at: 1588665926,
		}

		{
    		key: "notify",
			value: 0,
			updated_at: 1588665926,
		}

5. 运行项目
 	   go run cmd/main.go
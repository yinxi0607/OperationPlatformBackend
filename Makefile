
build:
	docker build -t registry.cn-hangzhou.aliyuncs.com/nhb/operation_platform_backend:latest .

	docker push registry.cn-hangzhou.aliyuncs.com/nhb/operation_platform_backend:latest

run:
	docker run -d -p 58180:58180 registry.cn-hangzhou.aliyuncs.com/nhb/operation_platform_backend:latest
## GET STARTED

**With Golang (require version 1.25+)**

```bash
make run
```

**With Docker**

```bash
docker build -t instay-be .
docker run -d -p 8080:8080 --name instay instay-be
```

### Project Structure

```
├── 📁 cmd
│   └── 📁 api
│       └── 🐹 main.go
├── 📁 configs
│   └── ⚙️ example.yaml
│   └── ⚙️ main.yaml
├── 📁 docs
│   ├── 🐹 docs.go
│   ├── ⚙️ swagger.json
│   └── ⚙️ swagger.yaml
├── 📁 internal
│   ├── 📁 common
│   │   ├── 🐹 constants.go
│   │   ├── 🐹 errors.go
│   │   ├── 🐹 mapper.go
│   │   └── 🐹 utils.go
│   ├── 📁 config
│   │   └── 🐹 main_config.go
│   ├── 📁 container
│   │   ├── 🐹 auth_container.go
│   │   ├── 🐹 booking_container.go
│   │   ├── 🐹 chat_container.go
│   │   ├── 🐹 dashboard_container.go
│   │   ├── 🐹 department_container.go
│   │   ├── 🐹 file_container.go
│   │   ├── 🐹 main_container.go
│   │   ├── 🐹 notification_container.go
│   │   ├── 🐹 order_container.go
│   │   ├── 🐹 request_container.go
│   │   ├── 🐹 review_container.go
│   │   ├── 🐹 room_container.go
│   │   ├── 🐹 service_container.go
│   │   ├── 🐹 sse_container.go
│   │   ├── 🐹 user_container.go
│   │   └── 🐹 ws_container.go
│   ├── 📁 handler
│   │   ├── 🐹 auth_handler.go
│   │   ├── 🐹 booking_handler.go
│   │   ├── 🐹 chat_handler.go
│   │   ├── 🐹 dashboard_handler.go
│   │   ├── 🐹 department_handler.go
│   │   ├── 🐹 file_handler.go
│   │   ├── 🐹 notification_handler.go
│   │   ├── 🐹 order_handler.go
│   │   ├── 🐹 request_handler.go
│   │   ├── 🐹 review_handler.go
│   │   ├── 🐹 room_handler.go
│   │   ├── 🐹 service_handler.go
│   │   ├── 🐹 sse_handler.go
│   │   ├── 🐹 user_handler.go
│   │   └── 🐹 ws_handler.go
│   ├── 📁 hub
│   │   ├── 🐹 sse_hub.go
│   │   └── 🐹 ws_hub.go
│   ├── 📁 initialization
│   │   ├── 🐹 logger.go
│   │   ├── 🐹 postgresql.go
│   │   ├── 🐹 rabbitmq.go
│   │   ├── 🐹 redis.go
│   │   ├── 🐹 s3.go
│   │   └── 🐹 snowflake.go
│   ├── 📁 middleware
│   │   ├── 🐹 authentication.go
│   │   └── 🐹 request.go
│   ├── 📁 model
│   │   ├── 🐹 booking_model.go
│   │   ├── 🐹 chat_model.go
│   │   ├── 🐹 department_model.go
│   │   ├── 🐹 notification_model.go
│   │   ├── 🐹 order_model.go
│   │   ├── 🐹 request_model.go
│   │   ├── 🐹 review_model.go
│   │   ├── 🐹 room_model.go
│   │   ├── 🐹 service_model.go
│   │   └── 🐹 user_model.go
│   ├── 📁 provider
│   │   ├── 📁 jwt
│   │   │   └── 🐹 jwt.go
│   │   ├── 📁 mq
│   │   │   └── 🐹 message_queue.go
│   │   └── 📁 smtp
│   │       ├── 📁 templates
│   │       │   └── 🌐 auth.html
│   │       └── 🐹 smtp.go
│   ├── 📁 repository
│   │   ├── 📁 implement
│   │   │   ├── 🐹 booking_repo_impl.go
│   │   │   ├── 🐹 chat_repo_impl.go
│   │   │   ├── 🐹 department_repo_impl.go
│   │   │   ├── 🐹 notification_repo_impl.go
│   │   │   ├── 🐹 order_repo_impl.go
│   │   │   ├── 🐹 request_repo_impl.go
│   │   │   ├── 🐹 review_repo_impl.go
│   │   │   ├── 🐹 room_repo_impl.go
│   │   │   ├── 🐹 service_repo_impl.go
│   │   │   └── 🐹 user_repo_impl.go
│   │   ├── 🐹 booking_repository.go
│   │   ├── 🐹 chat_repository.go
│   │   ├── 🐹 department_repository.go
│   │   ├── 🐹 notification_repository.go
│   │   ├── 🐹 order_repository.go
│   │   ├── 🐹 request_repository.go
│   │   ├── 🐹 review_repository.go
│   │   ├── 🐹 room_repository.go
│   │   ├── 🐹 service_repository.go
│   │   └── 🐹 user_repository.go
│   ├── 📁 router
│   │   ├── 🐹 auth_router.go
│   │   ├── 🐹 booking_router.go
│   │   ├── 🐹 chat_router.go
│   │   ├── 🐹 dashboard_router.go
│   │   ├── 🐹 department.go
│   │   ├── 🐹 file_router.go
│   │   ├── 🐹 notification_router.go
│   │   ├── 🐹 order_router.go
│   │   ├── 🐹 request_router.go
│   │   ├── 🐹 review_router.go
│   │   ├── 🐹 room_router.go
│   │   ├── 🐹 service_router.go
│   │   ├── 🐹 sse_router.go
│   │   ├── 🐹 user_router.go
│   │   └── 🐹 ws_router.go
│   ├── 📁 server
│   │   └── 🐹 server.go
│   ├── 📁 service
│   │   ├── 📁 implement
│   │   │   ├── 🐹 auth_svc_impl.go
│   │   │   ├── 🐹 booking_svc_impl.go
│   │   │   ├── 🐹 chat_svc_impl.go
│   │   │   ├── 🐹 dashboard_svc_impl.go
│   │   │   ├── 🐹 department_svc_impl.go
│   │   │   ├── 🐹 file_svc_impl.go
│   │   │   ├── 🐹 notification_svc_impl.go
│   │   │   ├── 🐹 order_svc_impl.go
│   │   │   ├── 🐹 request_svc_impl.go
│   │   │   ├── 🐹 review_svc_impl.go
│   │   │   ├── 🐹 room_svc_impl.go
│   │   │   ├── 🐹 service_svc_impl.go
│   │   │   └── 🐹 user_svc_impl.go
│   │   ├── 🐹 auth_service.go
│   │   ├── 🐹 booking_service.go
│   │   ├── 🐹 chat_service.go
│   │   ├── 🐹 dashboard_service.go
│   │   ├── 🐹 department_service.go
│   │   ├── 🐹 file_service.go
│   │   ├── 🐹 notification_service.go
│   │   ├── 🐹 order_service.go
│   │   ├── 🐹 request_service.go
│   │   ├── 🐹 review_service.go
│   │   ├── 🐹 room_service.go
│   │   ├── 🐹 service_service.go
│   │   └── 🐹 user_service.go
│   ├── 📁 types
│   │   ├── 🐹 data.go
│   │   ├── 🐹 request.go
│   │   └── 🐹 response.go
│   └── 📁 worker
│       ├── 🐹 listen_worker.go
│       └── 🐹 mq_worker.go
├── 📁 logs
│   └── 📄 app.log
├── 📁 pkg
│   ├── 📁 bcrypt
│   │   └── 🐹 bcrypt.go
│   └── 📁 snowflake
│       └── 🐹 snowflake.go
├── ⚙️ .dockerignore
├── ⚙️ .gitignore
├── 🐳 Dockerfile
├── 📄 Makefile
├── 📝 README.md
├── 📄 go.mod
└── 📄 go.sum
```

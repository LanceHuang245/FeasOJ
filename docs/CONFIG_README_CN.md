# 配置文件说明
当你第一次运行程序时，会自动创建一个 `config.toml` 文件，或者你可以手动在项目目录下创建一个 `config.toml` 文件。

## Config文件样例

### 后端配置文件样例
```toml
[server]
host = "127.0.0.1:37882"
enable_https = false
cert_path = "./certificate/fullchain.pem"
key_path = "./certificate/privkey.key"

[rabbitmq]
host = "amqp://username:password@127.0.0.1:5672/"

[consul]
host = "localhost:8500"

[features]
image_guard_enabled = false
profanity_detector_enabled = false

[database]
type = "postgresql"     # mysql or postgresql（大小写严格）
host = "localhost"
port = 5432
name = "feasoj"
user = "username"
password = "password"
ssl_mode = "disable"
max_open_conns = 240
max_idle_conns = 100
max_life_time = 32

[redis]
host = "localhost:6379"
password = ""

[mail]
host = "smtp.qq.com"
port = 465
user = "example@qq.com"
password = "password"

[jwt]
signing_method = "HS256"
token_expire_hours = 720
secret_key = "jwtsecretkey12313"
```

### JudgeCore Configuration
```toml
[consul]
host = "127.0.0.1:8500"
service_name = "JudgeCore"
service_id = "JudgeCore-1"

[rabbitmq]
host = "amqp://username:password@127.0.0.1:5672/"

[server]
host = "127.0.0.1"
port = 37885
enable_https = false
cert_path = "./certificate/fullchain.pem"
key_path = "./certificate/privkey.key"

[sandbox]
memory = 2147483648
nano_cpus = 0.5
cpu_shares = 1024
max_concurrent = 5

[database]
type = "postgresql"     # mysql 或 postgresql（大小写严格）
host = "localhost"
port = 5432
name = "feasoj"
user = "postgres"
password = "Sing5200"
ssl_mode = "disable"
```
# Configuration
When first running the program, it will automatically create a `config.toml` file, or you can manually create a `config.toml` file in the project directory.

## Configuration Example

### Backend Configuration
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
type = "postgresql"     # mysql or postgresql (case-sensitive)
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
type = "postgresql"     # postgresql or mysql (case-sensitive)
host = "localhost"
port = 5432
name = "feasoj"
user = "postgres"
password = "Sing5200"
ssl_mode = "disable"
```
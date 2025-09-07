## Package `config`

YAML configuration loader and strongly-typed configuration struct.

### Types

```go
type Config struct {
    Server struct {
        Port string `yaml:"port"`
    } `yaml:"server"`
    DB struct {
        Host     string `yaml:"host"`
        User     string `yaml:"user"`
        Password string `yaml:"password"`
        Database string `yaml:"database"`
    } `yaml:"db"`
    AI struct {
        Moonshot struct {
            BaseURL string `yaml:"baseURL"`
            APIKey  string `yaml:"api_key"`
            Model   string `yaml:"model"`
        } `yaml:"moonshot"`
        Deepseek struct {
            BaseURL       string `yaml:"baseURL"`
            APIKey        string `yaml:"api_key"`
            Model         string `yaml:"model"`
            SupportModels string `yaml:"supportModels"`
        } `yaml:"deepseek"`
    } `yaml:"ai"`
}
```

### Public Functions

```go
func LoadConfig(filename string) (*Config, error)
```
Reads YAML from `filename` into `Config`.

### Example: Load Config

```go
cfg, err := config.LoadConfig("config.yml")
if err != nil { panic(err) }
fmt.Println("DB host:", cfg.DB.Host)
```

### Example `config.yml`

```yaml
server:
  port: "30003"

db:
  host: "localhost"
  user: "root"
  password: "root"
  database: "zhuanz_java"

ai:
  moonshot:
    baseURL: "https://api.moonshot.cn/v1"
    api_key: "sk-xxxxxxxx"
    model: "moonshot-v1-32k"
  deepseek:
    baseURL: "https://api.deepseek.com"
    api_key: "sk-xxxxxxxx"
    model: "deepseek-chat"
    supportModels: "deepseek-chat,deepseek-reasoner"
```


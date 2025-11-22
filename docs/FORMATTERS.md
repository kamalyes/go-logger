# Go Logger 格式化器指南

## 目录

- [🎨 格式化器概述](#-格式化器概述)
- [📄 内置格式化器](#-内置格式化器)
- [🔧 自定义格式化器](#-自定义格式化器)
- [🌈 高级格式化](#-高级格式化)
- [📋 格式化模板](#-格式化模板)
- [🎯 最佳实践](#-最佳实践)

## 🎨 格式化器概述

go-logger 提供了强大而灵活的格式化系统，支持多种内置格式和完全自定义的格式化器。格式化器负责将日志记录转换为最终的输出格式。

### 格式化器架构

```
┌─────────────────────────────────────────────┐
│            格式化器接口 (Formatter)           │
└─────────────────┬───────────────────────────┘
                  │
┌─────────────────┼───────────────────────────┐
│                 │                           │
├─ 内置格式化器 ──┼─ 自定义格式化器 ──────────┤
│  • JSON         │  • Template               │
│  • Text         │  • Custom                 │
│  • Structured   │  • Composite              │
│  • CSV          │  • Conditional            │
│  • XML          │                           │
└─────────────────┴───────────────────────────┘
```

### 格式化器类型

| 类型 | 用途 | 性能 | 可读性 | 适用场景 |
|------|------|------|--------|----------|
| JSON | 结构化数据 | ⭐⭐⭐⭐ | ⭐⭐⭐ | API、日志分析 |
| Text | 人类可读 | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐⭐ | 开发、调试 |
| Structured | 键值对 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ | 系统监控 |
| CSV | 表格数据 | ⭐⭐⭐ | ⭐⭐⭐ | 数据分析 |
| XML | 企业集成 | ⭐⭐ | ⭐⭐ | 企业系统 |

## 📄 内置格式化器

### JSON 格式化器

最常用的结构化日志格式，适合日志收集和分析：

```go
import "github.com/kamalyes/go-logger/formatter"

// 创建JSON格式化器
jsonFormatter := formatter.NewJSONFormatter()

// 基础配置
jsonFormatter.SetTimestampFormat(time.RFC3339Nano)
jsonFormatter.SetDisableTimestamp(false)
jsonFormatter.SetDisableHTMLEscape(true)
jsonFormatter.SetPrettyPrint(false) // 紧凑格式

// 字段配置
jsonFormatter.SetFieldMap(formatter.FieldMap{
    formatter.FieldKeyTime:  "timestamp",
    formatter.FieldKeyLevel: "level",
    formatter.FieldKeyMsg:   "message",
    formatter.FieldKeyFunc:  "function",
    formatter.FieldKeyFile:  "source",
})

// 自定义字段
jsonFormatter.SetCallerPrettyfier(func(f *runtime.Frame) (string, string) {
    filename := path.Base(f.File)
    return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
})

// 应用到日志器
logger := logrus.New()
logger.SetFormatter(jsonFormatter)

// 示例输出
logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "action":  "login",
    "ip":      "192.168.1.100",
}).Info("用户登录成功")

// 输出:
// {"timestamp":"2024-01-01T10:30:45.123456789Z","level":"info","message":"用户登录成功","function":"main.handleLogin()","source":"main.go:42","user_id":12345,"action":"login","ip":"192.168.1.100"}
```

### JSON 高级配置

```go
// 创建高级JSON格式化器
advancedJSON := formatter.NewAdvancedJSONFormatter(&formatter.JSONConfig{
    // 时间配置
    TimestampFormat:   time.RFC3339Nano,
    DisableTimestamp:  false,
    TimestampKey:      "time",
    
    // 输出配置
    PrettyPrint:       true,  // 美化输出
    DisableHTMLEscape: true,  // 禁用HTML转义
    SortKeys:          true,  // 排序键名
    
    // 字段配置
    LevelKey:    "level",
    MessageKey:  "msg",
    ErrorKey:    "error",
    CallerKey:   "caller",
    StackKey:    "stack",
    
    // 数据类型配置
    DataKey:         "fields",      // 自定义字段的父键
    NestedFieldSeparator: ".",      // 嵌套字段分隔符
    
    // 错误处理
    ErrorFieldName:   "format_error",
    SkipErrorFields:  true,
    
    // Hook配置
    EnableHooks:      true,
    MaxFieldLength:   1024,     // 最大字段长度
    TruncateLongFields: true,   // 截断长字段
})

// 美化输出示例
logger.SetFormatter(advancedJSON)
logger.WithFields(logrus.Fields{
    "user": map[string]interface{}{
        "id":   12345,
        "name": "张三",
        "email": "zhangsan@example.com",
    },
    "request": map[string]interface{}{
        "method": "POST",
        "path":   "/api/users",
        "params": map[string]string{
            "action": "create",
            "source": "web",
        },
    },
}).Info("用户创建请求")

// 美化输出:
// {
//   "time": "2024-01-01T10:30:45.123456789Z",
//   "level": "info",
//   "msg": "用户创建请求",
//   "caller": "main.go:42",
//   "fields": {
//     "request": {
//       "method": "POST",
//       "path": "/api/users",
//       "params": {
//         "action": "create",
//         "source": "web"
//       }
//     },
//     "user": {
//       "email": "zhangsan@example.com",
//       "id": 12345,
//       "name": "张三"
//     }
//   }
// }
```

### Text 格式化器

人类友好的文本格式，适合开发和调试：

```go
// 创建文本格式化器
textFormatter := formatter.NewTextFormatter()

// 基础配置
textFormatter.SetTimestampFormat("2006-01-02 15:04:05")
textFormatter.SetFullTimestamp(true)
textFormatter.SetDisableColors(false)   // 启用颜色
textFormatter.SetDisableTimestamp(false)
textFormatter.SetDisableSorting(false)  // 启用字段排序

// 字段配置
textFormatter.SetFieldMap(formatter.FieldMap{
    formatter.FieldKeyTime:  "time",
    formatter.FieldKeyLevel: "level",
    formatter.FieldKeyMsg:   "msg",
})

// 调用者信息
textFormatter.SetCallerPrettyfier(func(f *runtime.Frame) (string, string) {
    filename := path.Base(f.File)
    return fmt.Sprintf("[%s]", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
})

// 应用到日志器
logger.SetFormatter(textFormatter)

// 示例输出
logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "action":  "login",
}).Info("用户登录成功")

// 输出 (带颜色):
// 2024-01-01 10:30:45 [INFO] 用户登录成功 action=login user_id=12345
```

### Text 高级配置

```go
// 创建高级文本格式化器
advancedText := formatter.NewAdvancedTextFormatter(&formatter.TextConfig{
    // 时间配置
    TimestampFormat:      "2006-01-02 15:04:05.000",
    FullTimestamp:        true,
    DisableTimestamp:     false,
    
    // 颜色配置
    ForceColors:          true,
    DisableColors:        false,
    EnvironmentOverrideColors: false,
    
    // 字段配置
    DisableSorting:       false,
    SortingFunc:         nil,  // 使用默认排序
    DisableLevelTruncation: false,
    PadLevelText:        true,   // 填充级别文本对齐
    
    // 引用配置
    QuoteEmptyFields:    true,   // 引用空字段
    QuoteCharacter:      "\"",   // 引用字符
    
    // 调用者配置
    CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
        dir := filepath.Dir(f.File)
        filename := filepath.Base(f.File)
        return fmt.Sprintf("%s()", filepath.Base(f.Function)), 
               fmt.Sprintf("%s/%s:%d", filepath.Base(dir), filename, f.Line)
    },
    
    // 自定义格式
    CustomFormat: func(entry *logrus.Entry) string {
        timestamp := entry.Time.Format("2006-01-02 15:04:05.000")
        level := strings.ToUpper(entry.Level.String())
        
        // 添加颜色
        var levelColor string
        switch entry.Level {
        case logrus.ErrorLevel:
            levelColor = "\033[31m" // 红色
        case logrus.WarnLevel:
            levelColor = "\033[33m" // 黄色
        case logrus.InfoLevel:
            levelColor = "\033[36m" // 青色
        case logrus.DebugLevel:
            levelColor = "\033[37m" // 白色
        default:
            levelColor = "\033[0m"  // 默认色
        }
        
        var fields string
        for k, v := range entry.Data {
            fields += fmt.Sprintf(" %s=%v", k, v)
        }
        
        return fmt.Sprintf("%s %s[%s]\033[0m %s%s", 
            timestamp, levelColor, level, entry.Message, fields)
    },
})

logger.SetFormatter(advancedText)

// 自定义颜色方案
colorScheme := formatter.ColorScheme{
    InfoLevelColor:  formatter.ColorCyan,
    WarnLevelColor:  formatter.ColorYellow,
    ErrorLevelColor: formatter.ColorRed,
    FatalLevelColor: formatter.ColorMagenta,
    PanicLevelColor: formatter.ColorRed,
    DebugLevelColor: formatter.ColorWhite,
    TraceColor:      formatter.ColorGray,
}
advancedText.SetColorScheme(colorScheme)
```

### Structured 格式化器

键值对格式，平衡了可读性和结构化：

```go
// 创建结构化格式化器
structFormatter := formatter.NewStructuredFormatter()

// 配置选项
structFormatter.SetOptions(&formatter.StructuredOptions{
    TimestampFormat:     "2006-01-02T15:04:05.000Z07:00",
    KeyValueSeparator:   "=",
    FieldSeparator:      " ",
    QuoteValues:         true,
    QuoteKeys:           false,
    EscapeQuotes:        true,
    ShowCaller:          true,
    ShowTimestamp:       true,
    ShowLevel:           true,
    UppercaseLevel:      true,
    
    // 字段顺序
    FieldOrder: []string{"timestamp", "level", "caller", "message"},
    
    // 字段映射
    FieldNames: map[string]string{
        "timestamp": "time",
        "level":     "lvl",
        "message":   "msg",
        "caller":    "src",
    },
})

logger.SetFormatter(structFormatter)

logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "action":  "login",
    "duration": "150ms",
}).Info("用户操作完成")

// 输出:
// time="2024-01-01T10:30:45.123Z" lvl="INFO" src="main.go:42" msg="用户操作完成" action="login" duration="150ms" user_id="12345"
```

### CSV 格式化器

表格格式，适合数据分析：

```go
// 创建CSV格式化器
csvFormatter := formatter.NewCSVFormatter()

// 配置CSV选项
csvFormatter.SetOptions(&formatter.CSVOptions{
    Separator:           ",",
    Quote:               "\"",
    Header:              true,    // 输出表头
    TimestampFormat:     time.RFC3339,
    EscapeQuotes:        true,
    
    // 定义列
    Columns: []formatter.CSVColumn{
        {Name: "timestamp", Field: "time", Type: "datetime"},
        {Name: "level", Field: "level", Type: "string"},
        {Name: "message", Field: "message", Type: "string"},
        {Name: "user_id", Field: "user_id", Type: "integer"},
        {Name: "action", Field: "action", Type: "string"},
        {Name: "duration_ms", Field: "duration", Type: "integer", 
         Converter: func(v interface{}) interface{} {
             if s, ok := v.(string); ok {
                 if d, err := time.ParseDuration(s); err == nil {
                     return int64(d / time.Millisecond)
                 }
             }
             return 0
         }},
    },
    
    // 缺失值处理
    DefaultValues: map[string]interface{}{
        "user_id": 0,
        "action":  "unknown",
        "duration_ms": 0,
    },
})

logger.SetFormatter(csvFormatter)

// 第一次调用时输出表头:
// timestamp,level,message,user_id,action,duration_ms
logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "action":  "login",
    "duration": "150ms",
}).Info("用户登录")

// 输出:
// "2024-01-01T10:30:45Z","info","用户登录",12345,"login",150
```

### XML 格式化器

XML格式，适合企业系统集成：

```go
// 创建XML格式化器
xmlFormatter := formatter.NewXMLFormatter()

// 配置XML选项
xmlFormatter.SetOptions(&formatter.XMLOptions{
    RootElement:     "LogEntry",
    TimestampFormat: time.RFC3339,
    Indent:          "  ",      // 缩进
    PrettyPrint:     true,      // 美化输出
    
    // 元素映射
    ElementNames: map[string]string{
        "timestamp": "Timestamp",
        "level":     "Level",
        "message":   "Message",
        "caller":    "Source",
        "data":      "Fields",
    },
    
    // 属性配置
    UseAttributes: true,
    AttributeMapping: map[string]string{
        "level": "level",
        "time":  "timestamp",
    },
    
    // 命名空间
    Namespace: "http://your-company.com/logging/v1",
    
    // CDATA 包装的字段
    CDATAFields: []string{"message", "error"},
})

logger.SetFormatter(xmlFormatter)

logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "action":  "login",
}).Error("登录失败")

// 输出:
// <?xml version="1.0" encoding="UTF-8"?>
// <LogEntry xmlns="http://your-company.com/logging/v1" level="error" timestamp="2024-01-01T10:30:45Z">
//   <Timestamp>2024-01-01T10:30:45Z</Timestamp>
//   <Level>error</Level>
//   <Message><![CDATA[登录失败]]></Message>
//   <Source>main.go:42</Source>
//   <Fields>
//     <Field name="action" type="string">login</Field>
//     <Field name="user_id" type="integer">12345</Field>
//   </Fields>
// </LogEntry>
```

## 🔧 自定义格式化器

### 基础自定义格式化器

```go
// 实现 Formatter 接口
type CustomFormatter struct {
    TimestampFormat string
    ShowColors      bool
    Prefix          string
}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    // 时间戳
    timestamp := entry.Time.Format(f.TimestampFormat)
    
    // 级别
    level := strings.ToUpper(entry.Level.String())
    if f.ShowColors {
        level = f.colorizeLevel(level)
    }
    
    // 构建消息
    var message strings.Builder
    
    // 前缀
    if f.Prefix != "" {
        message.WriteString(fmt.Sprintf("[%s] ", f.Prefix))
    }
    
    // 基础格式
    message.WriteString(fmt.Sprintf("%s [%s] %s", timestamp, level, entry.Message))
    
    // 添加字段
    if len(entry.Data) > 0 {
        message.WriteString(" |")
        for key, value := range entry.Data {
            message.WriteString(fmt.Sprintf(" %s:%v", key, value))
        }
    }
    
    message.WriteString("\n")
    return []byte(message.String()), nil
}

func (f *CustomFormatter) colorizeLevel(level string) string {
    switch level {
    case "ERROR":
        return "\033[31m" + level + "\033[0m" // 红色
    case "WARN":
        return "\033[33m" + level + "\033[0m" // 黄色
    case "INFO":
        return "\033[32m" + level + "\033[0m" // 绿色
    case "DEBUG":
        return "\033[36m" + level + "\033[0m" // 青色
    default:
        return level
    }
}

// 使用自定义格式化器
customFormatter := &CustomFormatter{
    TimestampFormat: "2006-01-02 15:04:05",
    ShowColors:      true,
    Prefix:          "APP",
}

logger.SetFormatter(customFormatter)

logger.WithFields(logrus.Fields{
    "user": "张三",
    "ip":   "192.168.1.100",
}).Info("用户访问")

// 输出:
// [APP] 2024-01-01 10:30:45 [INFO] 用户访问 | user:张三 ip:192.168.1.100
```

### 模板格式化器

使用 Go 模板引擎的强大功能：

```go
// 创建模板格式化器
templateFormatter := formatter.NewTemplateFormatter()

// 设置模板
template := `{{.Timestamp.Format "2006-01-02 15:04:05"}} [{{.Level | upper}}] {{.Message}}
{{- if .Fields}}
  Fields:
  {{- range $key, $value := .Fields}}
    {{$key}}: {{$value}}
  {{- end}}
{{- end}}
{{- if .Error}}
  Error: {{.Error}}
{{- end}}
{{- if .Caller}}
  Source: {{.Caller.Function}} ({{.Caller.File}}:{{.Caller.Line}})
{{- end}}
`

err := templateFormatter.SetTemplate(template)
if err != nil {
    log.Fatal("设置模板失败:", err)
}

// 注册自定义函数
templateFormatter.RegisterFunc("upper", strings.ToUpper)
templateFormatter.RegisterFunc("lower", strings.ToLower)
templateFormatter.RegisterFunc("formatDuration", func(d time.Duration) string {
    return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
})

logger.SetFormatter(templateFormatter)

// 使用示例
logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "request_id": "req-abc123",
    "duration": time.Millisecond * 150,
}).Info("请求处理完成")

// 输出:
// 2024-01-01 10:30:45 [INFO] 请求处理完成
//   Fields:
//     duration: 150ms
//     request_id: req-abc123
//     user_id: 12345
//   Source: main.handleRequest (handler.go:42)
```

### 条件格式化器

根据条件选择不同的格式化器：

```go
// 创建条件格式化器
conditionalFormatter := formatter.NewConditionalFormatter()

// 配置条件和对应的格式化器
conditionalFormatter.AddCondition(
    func(entry *logrus.Entry) bool {
        // 错误级别使用JSON格式
        return entry.Level == logrus.ErrorLevel
    },
    formatter.NewJSONFormatter(),
)

conditionalFormatter.AddCondition(
    func(entry *logrus.Entry) bool {
        // 生产环境使用结构化格式
        return os.Getenv("ENV") == "production"
    },
    formatter.NewStructuredFormatter(),
)

// 默认格式化器 (文本格式)
conditionalFormatter.SetDefaultFormatter(formatter.NewTextFormatter())

logger.SetFormatter(conditionalFormatter)

// 不同条件下的输出
logger.Info("普通信息")        // 使用文本格式
logger.Error("发生错误")       // 使用JSON格式
```

### 复合格式化器

组合多个格式化器：

```go
// 创建复合格式化器
compositeFormatter := formatter.NewCompositeFormatter()

// 添加多个输出格式
compositeFormatter.AddFormatter("console", formatter.NewTextFormatter())
compositeFormatter.AddFormatter("file", formatter.NewJSONFormatter())
compositeFormatter.AddFormatter("audit", &AuditFormatter{})

// 配置路由规则
compositeFormatter.SetRoutingRules(map[string]func(*logrus.Entry) bool{
    "console": func(entry *logrus.Entry) bool {
        return entry.Level >= logrus.InfoLevel
    },
    "file": func(entry *logrus.Entry) bool {
        return true // 所有日志写入文件
    },
    "audit": func(entry *logrus.Entry) bool {
        // 只有审计日志
        return entry.Data["audit"] == true
    },
})

logger.SetFormatter(compositeFormatter)

// 审计日志示例
logger.WithFields(logrus.Fields{
    "audit": true,
    "user_id": 12345,
    "action": "sensitive_operation",
}).Warn("执行敏感操作")
```

## 🌈 高级格式化

### 字段过滤器

```go
// 创建字段过滤器格式化器
filterFormatter := formatter.NewFieldFilterFormatter(
    formatter.NewJSONFormatter(), // 基础格式化器
)

// 配置过滤规则
filterFormatter.SetFieldFilter(&formatter.FieldFilter{
    // 包含的字段
    IncludeFields: []string{"user_id", "action", "timestamp", "message"},
    
    // 排除的字段
    ExcludeFields: []string{"password", "secret", "token"},
    
    // 敏感字段脱敏
    SensitiveFields: map[string]formatter.MaskFunc{
        "email": formatter.MaskEmail,
        "phone": formatter.MaskPhone,
        "ip":    formatter.MaskIP,
        "id_card": func(value interface{}) interface{} {
            if s, ok := value.(string); ok && len(s) > 6 {
                return s[:3] + "***" + s[len(s)-3:]
            }
            return "***"
        },
    },
    
    // 字段重命名
    FieldMapping: map[string]string{
        "user_id": "uid",
        "request_id": "req_id",
    },
    
    // 字段验证
    FieldValidators: map[string]func(interface{}) bool{
        "user_id": func(v interface{}) bool {
            if id, ok := v.(int); ok {
                return id > 0
            }
            return false
        },
    },
})

logger.SetFormatter(filterFormatter)

logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "email": "user@example.com",
    "password": "secret123",
    "action": "login",
}).Info("用户登录")

// 输出 (敏感信息被脱敏):
// {"timestamp":"2024-01-01T10:30:45Z","level":"info","message":"用户登录","uid":12345,"email":"u***@example.com","action":"login"}
```

### 动态格式化器

根据运行时条件动态调整格式：

```go
// 创建动态格式化器
dynamicFormatter := formatter.NewDynamicFormatter()

// 注册格式化器
dynamicFormatter.RegisterFormatter("json", formatter.NewJSONFormatter())
dynamicFormatter.RegisterFormatter("text", formatter.NewTextFormatter())
dynamicFormatter.RegisterFormatter("csv", formatter.NewCSVFormatter())

// 设置选择策略
dynamicFormatter.SetSelectionStrategy(func(entry *logrus.Entry) string {
    // 根据字段选择格式化器
    if format, exists := entry.Data["format"]; exists {
        if f, ok := format.(string); ok {
            return f
        }
    }
    
    // 根据级别选择
    if entry.Level >= logrus.ErrorLevel {
        return "json"  // 错误使用JSON便于分析
    }
    
    // 根据环境选择
    if os.Getenv("LOG_FORMAT") != "" {
        return os.Getenv("LOG_FORMAT")
    }
    
    return "text" // 默认文本格式
})

logger.SetFormatter(dynamicFormatter)

// 使用示例
logger.Info("普通消息")  // 使用文本格式

logger.WithField("format", "json").Info("JSON格式消息")  // 强制JSON格式

logger.Error("错误消息")  // 自动使用JSON格式
```

### 分级格式化器

不同级别使用不同格式：

```go
// 创建分级格式化器
levelFormatter := formatter.NewLevelBasedFormatter()

// 配置不同级别的格式化器
levelFormatter.SetFormatterForLevel(logrus.DebugLevel, formatter.NewTextFormatter())
levelFormatter.SetFormatterForLevel(logrus.InfoLevel, formatter.NewTextFormatter())
levelFormatter.SetFormatterForLevel(logrus.WarnLevel, formatter.NewStructuredFormatter())
levelFormatter.SetFormatterForLevel(logrus.ErrorLevel, formatter.NewJSONFormatter())
levelFormatter.SetFormatterForLevel(logrus.FatalLevel, formatter.NewJSONFormatter())

// 设置范围格式化器
levelFormatter.SetFormatterForLevelRange(
    logrus.DebugLevel, logrus.InfoLevel,
    &formatter.SimpleFormatter{Format: "[{level}] {message}\n"},
)

logger.SetFormatter(levelFormatter)
```

## 📋 格式化模板

### 内置模板

```go
// 使用内置模板
templates := formatter.GetBuiltinTemplates()

// 简单模板
simpleTemplate := templates["simple"]
// 格式: "2006-01-02 15:04:05 [LEVEL] message"

// 详细模板
detailedTemplate := templates["detailed"]
// 格式: "2006-01-02 15:04:05.000 [LEVEL] source:line - message {fields}"

// 紧凑模板
compactTemplate := templates["compact"]
// 格式: "15:04:05 LVL msg fields"

// 调试模板
debugTemplate := templates["debug"]
// 包含完整的调用栈和所有调试信息

// 应用模板
templateFormatter := formatter.NewTemplateFormatter()
templateFormatter.SetTemplate(detailedTemplate)
logger.SetFormatter(templateFormatter)
```

### 自定义模板

```go
// 创建自定义模板
customTemplates := map[string]string{
    "api_log": `{{.Timestamp.Format "2006-01-02T15:04:05.000Z07:00"}} [{{.Level | upper}}] {{.Message}}
{{- if .Fields.method}} Method: {{.Fields.method}}{{end}}
{{- if .Fields.path}} Path: {{.Fields.path}}{{end}}
{{- if .Fields.status}} Status: {{.Fields.status}}{{end}}
{{- if .Fields.duration}} Duration: {{.Fields.duration}}{{end}}
{{- if .Fields.user_id}} User: {{.Fields.user_id}}{{end}}`,

    "error_log": `🚨 ERROR REPORT 🚨
Time: {{.Timestamp.Format "2006-01-02 15:04:05"}}
Level: {{.Level | upper}}
Message: {{.Message}}
{{- if .Error}}
Error Details: {{.Error}}
{{- end}}
{{- if .Caller}}
Source: {{.Caller.Function}}
File: {{.Caller.File}}:{{.Caller.Line}}
{{- end}}
{{- if .Fields}}
Additional Info:
{{- range $key, $value := .Fields}}
  {{$key}}: {{$value}}
{{- end}}
{{- end}}
=================================`,

    "audit_log": `[AUDIT] {{.Timestamp.Format "2006-01-02 15:04:05"}}
User: {{.Fields.user_id | default "anonymous"}}
Action: {{.Fields.action}}
Resource: {{.Fields.resource | default "unknown"}}
Result: {{.Fields.result | default "success"}}
IP: {{.Fields.ip | default "unknown"}}
{{- if .Fields.details}}
Details: {{.Fields.details}}
{{- end}}`,
}

// 注册模板
for name, template := range customTemplates {
    formatter.RegisterTemplate(name, template)
}

// 使用模板
apiFormatter := formatter.NewTemplateFormatter()
apiFormatter.SetTemplate(customTemplates["api_log"])

errorFormatter := formatter.NewTemplateFormatter()
errorFormatter.SetTemplate(customTemplates["error_log"])

// 在不同场景中使用
logger.SetFormatter(apiFormatter)
logger.WithFields(logrus.Fields{
    "method": "POST",
    "path": "/api/users",
    "status": 201,
    "duration": "150ms",
    "user_id": 12345,
}).Info("API请求完成")

logger.SetFormatter(errorFormatter)
logger.WithFields(logrus.Fields{
    "user_id": 12345,
    "operation": "database_query",
}).Error("数据库连接失败")
```

### 模板函数库

```go
// 注册自定义模板函数
templateFuncs := map[string]interface{}{
    // 字符串操作
    "upper":    strings.ToUpper,
    "lower":    strings.ToLower,
    "title":    strings.Title,
    "trim":     strings.TrimSpace,
    "truncate": func(s string, length int) string {
        if len(s) <= length {
            return s
        }
        return s[:length] + "..."
    },
    
    // 数值操作
    "add": func(a, b int) int { return a + b },
    "sub": func(a, b int) int { return a - b },
    "mul": func(a, b int) int { return a * b },
    "div": func(a, b int) int { 
        if b != 0 { return a / b }
        return 0 
    },
    
    // 时间操作
    "formatDuration": func(d interface{}) string {
        if duration, ok := d.(time.Duration); ok {
            if duration < time.Millisecond {
                return fmt.Sprintf("%.2fμs", float64(duration.Nanoseconds())/1000)
            } else if duration < time.Second {
                return fmt.Sprintf("%.2fms", float64(duration.Nanoseconds())/1e6)
            }
            return duration.String()
        }
        return fmt.Sprintf("%v", d)
    },
    
    "timeAgo": func(t time.Time) string {
        duration := time.Since(t)
        if duration < time.Minute {
            return "刚才"
        } else if duration < time.Hour {
            return fmt.Sprintf("%d分钟前", int(duration.Minutes()))
        } else if duration < time.Hour*24 {
            return fmt.Sprintf("%d小时前", int(duration.Hours()))
        }
        return fmt.Sprintf("%d天前", int(duration.Hours()/24))
    },
    
    // 条件操作
    "default": func(defaultVal, value interface{}) interface{} {
        if value == nil || value == "" {
            return defaultVal
        }
        return value
    },
    
    "isEmpty": func(value interface{}) bool {
        return value == nil || value == ""
    },
    
    // JSON操作
    "toJSON": func(v interface{}) string {
        if data, err := json.Marshal(v); err == nil {
            return string(data)
        }
        return fmt.Sprintf("%v", v)
    },
    
    "prettyJSON": func(v interface{}) string {
        if data, err := json.MarshalIndent(v, "", "  "); err == nil {
            return string(data)
        }
        return fmt.Sprintf("%v", v)
    },
    
    // 数据脱敏
    "maskEmail": func(email string) string {
        parts := strings.Split(email, "@")
        if len(parts) != 2 {
            return "***@***.***"
        }
        username := parts[0]
        domain := parts[1]
        if len(username) <= 2 {
            return "***@" + domain
        }
        return username[:1] + "***" + username[len(username)-1:] + "@" + domain
    },
    
    "maskPhone": func(phone string) string {
        if len(phone) < 7 {
            return "***"
        }
        return phone[:3] + "****" + phone[len(phone)-4:]
    },
}

// 注册函数到模板引擎
templateFormatter.RegisterFuncs(templateFuncs)

// 使用模板函数的复杂模板
complexTemplate := `{{.Timestamp.Format "2006-01-02 15:04:05"}} [{{.Level | upper}}] {{.Message | truncate 100}}
{{- if .Fields.email}} Email: {{.Fields.email | maskEmail}}{{end}}
{{- if .Fields.phone}} Phone: {{.Fields.phone | maskPhone}}{{end}}
{{- if .Fields.duration}} Duration: {{.Fields.duration | formatDuration}}{{end}}
{{- if .Fields.count}} Count: {{.Fields.count | default 0}}{{end}}
{{- if .Fields.data}} Data: {{.Fields.data | prettyJSON}}{{end}}`

templateFormatter.SetTemplate(complexTemplate)
```

## 🎯 最佳实践

### 性能优化

```go
// 1. 使用对象池减少内存分配
var formatterPool = sync.Pool{
    New: func() interface{} {
        return &CustomFormatter{
            buffer: make([]byte, 0, 1024), // 预分配缓冲区
        }
    },
}

type PooledFormatter struct {
    buffer []byte
}

func (f *PooledFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    formatter := formatterPool.Get().(*CustomFormatter)
    defer formatterPool.Put(formatter)
    
    formatter.buffer = formatter.buffer[:0] // 重置缓冲区
    
    // 格式化逻辑...
    result := make([]byte, len(formatter.buffer))
    copy(result, formatter.buffer)
    
    return result, nil
}

// 2. 使用预编译的模板
var compiledTemplates = map[string]*template.Template{}

func init() {
    // 预编译常用模板
    for name, tmpl := range builtinTemplates {
        compiled, err := template.New(name).Parse(tmpl)
        if err == nil {
            compiledTemplates[name] = compiled
        }
    }
}

// 3. 避免反射和类型断言
func optimizedFormat(entry *logrus.Entry) string {
    var buf strings.Builder
    buf.Grow(256) // 预分配容量
    
    // 直接字符串操作，避免反射
    buf.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
    buf.WriteByte(' ')
    buf.WriteByte('[')
    buf.WriteString(strings.ToUpper(entry.Level.String()))
    buf.WriteByte(']')
    buf.WriteByte(' ')
    buf.WriteString(entry.Message)
    
    return buf.String()
}
```

### 错误处理

```go
// 安全的格式化器
type SafeFormatter struct {
    fallbackFormatter logrus.Formatter
    maxRetries        int
}

func (f *SafeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
    defer func() {
        if r := recover(); r != nil {
            log.Printf("格式化器发生panic: %v", r)
        }
    }()
    
    // 尝试格式化
    for i := 0; i < f.maxRetries; i++ {
        result, err := f.tryFormat(entry)
        if err == nil {
            return result, nil
        }
        
        log.Printf("格式化失败，重试 %d/%d: %v", i+1, f.maxRetries, err)
    }
    
    // 使用备用格式化器
    if f.fallbackFormatter != nil {
        return f.fallbackFormatter.Format(entry)
    }
    
    // 最后的备用方案
    return []byte(fmt.Sprintf("%s [%s] %s\n", 
        entry.Time.Format(time.RFC3339), 
        strings.ToUpper(entry.Level.String()), 
        entry.Message)), nil
}

func (f *SafeFormatter) tryFormat(entry *logrus.Entry) ([]byte, error) {
    // 实际格式化逻辑
    return nil, nil
}
```

### 配置管理

```go
// 格式化器配置管理
type FormatterConfig struct {
    Type     string                 `yaml:"type" json:"type"`
    Options  map[string]interface{} `yaml:"options" json:"options"`
    Template string                 `yaml:"template" json:"template"`
}

func CreateFormatterFromConfig(config FormatterConfig) (logrus.Formatter, error) {
    switch config.Type {
    case "json":
        formatter := formatter.NewJSONFormatter()
        if config.Options != nil {
            applyJSONOptions(formatter, config.Options)
        }
        return formatter, nil
        
    case "text":
        formatter := formatter.NewTextFormatter()
        if config.Options != nil {
            applyTextOptions(formatter, config.Options)
        }
        return formatter, nil
        
    case "template":
        formatter := formatter.NewTemplateFormatter()
        if config.Template != "" {
            if err := formatter.SetTemplate(config.Template); err != nil {
                return nil, err
            }
        }
        return formatter, nil
        
    default:
        return nil, fmt.Errorf("未知的格式化器类型: %s", config.Type)
    }
}

// 配置文件示例
var formatterConfigs = map[string]FormatterConfig{
    "development": {
        Type: "text",
        Options: map[string]interface{}{
            "timestamp_format": "2006-01-02 15:04:05",
            "disable_colors": false,
            "full_timestamp": true,
        },
    },
    
    "production": {
        Type: "json",
        Options: map[string]interface{}{
            "timestamp_format": time.RFC3339Nano,
            "pretty_print": false,
            "disable_html_escape": true,
        },
    },
    
    "audit": {
        Type: "template",
        Template: `[AUDIT] {{.Timestamp.Format "2006-01-02 15:04:05"}} {{.Level | upper}} {{.Message}}{{range $k, $v := .Fields}} {{$k}}={{$v}}{{end}}`,
    },
}
```

### 测试策略

```go
// 格式化器测试
func TestCustomFormatter(t *testing.T) {
    formatter := &CustomFormatter{
        TimestampFormat: "2006-01-02 15:04:05",
        ShowColors: false,
    }
    
    entry := &logrus.Entry{
        Time:    time.Date(2024, 1, 1, 10, 30, 45, 0, time.UTC),
        Level:   logrus.InfoLevel,
        Message: "测试消息",
        Data: logrus.Fields{
            "user_id": 12345,
            "action":  "test",
        },
    }
    
    result, err := formatter.Format(entry)
    assert.NoError(t, err)
    
    expected := "2024-01-01 10:30:45 [INFO] 测试消息 | action:test user_id:12345\n"
    assert.Equal(t, expected, string(result))
}

// 基准测试
func BenchmarkFormatter(b *testing.B) {
    formatter := &CustomFormatter{}
    entry := &logrus.Entry{
        Time:    time.Now(),
        Level:   logrus.InfoLevel,
        Message: "基准测试消息",
        Data:    logrus.Fields{"key": "value"},
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        formatter.Format(entry)
    }
}
```

---

更多格式化器相关信息请参考：

- [📊 性能详解](PERFORMANCE.md) - 格式化器性能优化
- [🔧 配置指南](CONFIGURATION.md) - 详细配置说明
- [📚 使用指南](USAGE.md) - 完整使用示例
- [🧩 适配器系统](ADAPTERS.md) - 适配器与格式化器集成
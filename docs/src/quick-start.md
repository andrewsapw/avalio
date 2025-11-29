# avalio

`avalio` - это простая программа для отслеживания доступности разных ресурсов. Цели проекта:

- предоставить простой инструмент мониторинга
- конфигурацию через файл, без использования GUI
- минимальное потребление ресурсов - идеально для запуска на Raspberry PI или дешевом VPS
- сборка в один испольняемый файл, отсутствие внешних зависимостей


## Сборка

Проект собирается как обычный Go-проект:

```bash
$ go build .
```

Это создаст исполняемый файл `./avalio`

## Создание конфигурации

На данный момеент, `avalio` конфигурируется с помощью .toml файла. Путь до которого передается аргуемнтом `-config`. Например:

```
$ avalio -config ./config.toml
```

Создадим пример конфигурационного файла:

```toml
log_level = "debug"

[[resources.http]]
name = 'example'
url = 'https://example.com'
expected_status = 200

[[resources.http]]
name = 'google'
url = 'https://google.com'
expected_status = 200

[[notificators.telegram]]
token = "..."
chat_id = "..."
name = 'telegram'

[[monitors.cron]]
name = 'everyday'
resources = ['example', 'google']
notificators = ['telegram']
cron = '* * * * *'
```

А теперь опишем все по порядку. Первая строка задает уровень логгирования:

```
log_level = "debug"
```

Далее инициализируются два HTTP-ресурса:

```toml
[[resources.http]]
name = 'example'
url = 'https://example.com'
expected_status = 200

[[resources.http]]
name = 'google'
url = 'https://google.com'
expected_status = 200
```

Эти ресурсы описывают _что будет проверятся_. В нашем случае, мы будем проверять доступность двух ресурсов: `https://example.com` и `https://google.com`.

После ресурсов, описываем нотификаторы, то есть каналы, по которым будут приходить уведомления. В нашем файле конфигурации задается лишь один канал, Telegram:

```toml
[[notificators.telegram]]
name = 'bot'
token = "..."
chat_id = "..."
```

И наконец, задаем настройки мониторов, они соденяют в себе ресурсы и нотификаторы, то есть описываю _как проверять_ ресурсы и _куда отправлять_ уведомления. 

```toml
[[monitors.cron]]
name = 'every-minute'
resources = ['example', 'google']
notificators = ['bot']
cron = '* * * * *'
```

Мы задали монитор типа `cron`, который будет проверять ресурсы `example` и `google`, и отправлять уведомления через Telegram-нотификатор `bot`.

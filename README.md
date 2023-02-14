Cервер представляет собой простую бэкенд-систему для приложения, реализующего взаимодействие с базой данных сотрудников некоторой фирмы. Сервер предоставляет клиентам следующий REST API:

* POST   /employee/           :  создаёт нового сотрудника и возвращает его ID
* GET    /employee/<id>       :  возвращает одного сотрудника по его ID
* GET    /employee/           :  возвращает всех сотрудников
* GET    /employee/<lastName> :  возвращает список сотрудников с указанной фамилией
* DELETE /employee/<id>       :  удаляет сотрудника по ID
* PUT    /employee/<id>       :  обновляет информацию о сотруднике по его ID

Сервер можно запустить командой:

```bash
    $ SERVERPORT=4112 go run ./cmd/main.go
```

В качестве SERVERPORT можно использовать любой порт, который будет прослушивать локальный сервер в ожидании подключений.

Сгенерировать сертификат и ключ можно командой:

```bash
    $ go run /usr/local/go/src/crypto/tls/generate_cert.go --ecdsa-curve P256 --host localhost
```

Протестировать сервер можно с помощью curl:

```bash
    $ curl --cacert cert.pem https://localhost:SERVERPORT/
```

Запустить клиентский код можно командой:

```bash
    $ go run pkg/client/auth-client.go -user NAME -pass PASSWORD -addr localhost:4112/employee/
```
где NAME, PASSWORD - имя и пароль пользователя из authdb (тестовые значения: {joe, 1234}, {mary, 5678})
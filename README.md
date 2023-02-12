Cервер представляет собой простую бэкенд-систему для приложения, реализующего взаимодействие с базой данных сотрудников некоторой фирмы. Сервер предоставляет клиентам следующий REST API:

* POST   /employee/           :  создаёт нового сотрудника и возвращает его ID
* GET    /employee/<id>       :  возвращает одного сотрудника по его ID
* GET    /employee/           :  возвращает всех сотрудников
* GET    /employee/<lastName> :  возвращает список сотрудников с указанной фамилией
* DELETE /employee/           :  удаляет всех сотрудников
* DELETE /employee/<id>       :  удаляет сотрудника по ID
* PUT    /employee/<id>       :  обновляет информацию о сотруднике по его ID

Сервер можно запустить командой:

$ SERVERPORT=4112 go run .

В качестве SERVERPORT можно использовать любой порт, который будет прослушивать локальный сервер в ожидании подключений.

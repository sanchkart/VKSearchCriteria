Требуется разработать систему поиска пользователей ВК по критериям. Система получает запросы с критериями, запускает сбор id, сохраняет результаты в БД. 
Далее результаты возвращаются по HTTP по запросу.

Язык Go, БД Postgres. JSON HTTP API.


Конфиг (в формате JSON).
	- параметры подключения к БД;
	- количество параллельно работюащих горутин по сбору пользователей ВК;
	- набор токенов для работы с API ВК.


Внешнее API.


	1. Запрос создания задачи пересечения подписчиков сообществ.

	/tsa/members_intersect
	Вход:
		auth - токен авторизации
		groups - массив сообществ
		members_min - минимальное число сообществ, в которых должен состоять пользователь;

	Выход:
		request_id - идентификатор заявки (UUID).


	Пример запроса:
		{
			"auth": "46259032-320e-4a42-b610-270661feb008",
			"groups": [
				"https://vk.com/vseokurske",
				"https://vk.com/like_kursk",
				"https://vk.com/prolashes_ru"
			],
			"member_min": 2
		}

	2. Запрос получения результатов. Вызывается пока finished в ответе равен false. offset увеличивается, чтобы получить все результаты частями.

	/tsa/get_result
	Вход:
		auth - токен авторизации
		request_id - идентификатор заявки
		offset - смещение результатов

	Выход:
		ids - массив идентификаторов.
		finished - признак того, что обработка задачи закончилась.

	Пример запроса:
	{
		"auth": "46259032-320e-4a42-b610-270661feb008",
		"request_id": "9df7bf1c-0de5-4824-b6dd-f78bc1e2b4ea",
		"offset": 0
	}

	Авторизация к обоим методам по ключу, который нужно искать в таблице user.


Дополнительные требования.
	- после перезагрузки машины работа над незавершенными задачами должна начаться автоматически;

	- работа с задачами по получению списков пользователей не должна блокировать HTTP сервис;

	- работа над новыми заданиями должна начинаться как можно быстрее;

	- обработка заданий должна идти параллельно;


Схема БД.

	1. таблица user:
		user_uuid - уникальный идентификатор пользователя;
		key - длинный строковый ключ для авторизации (приходит в запросе);
		name - имя пользователя, чтобы понимать кто это;

	2. таблица request:
		request_uuid - уникальный идентификатор заявки;
		user_uuid - идентификатор пользователя, создавшего заявку;
		type - тип заявки. пока что тип один - пересечение сообществ, позже будут новые.
		created_at - время создания заявки;
		status - статус заявки (PROCESSING - в работе, DONE - успешно завершена, CANCELLED - отменена);
		params - JSON с параметрами заявки. структура зависит от типа заявки (поле type). для пересечения групп выглядит примерно так:
		{
			groups: ['http://vk.com/g1', 'http://vk.com/g2', 'http://vk.com/g3'],
			members_min: 2
		}

	3. таблица result:
		result_id - идентификатор строки результата;
		request_uuid - идентификатор заявки;
		id - идентификатор пользователя;
		added_at - время, когда результат был добавлен;

API vk.com.
	Для получения списка пользователей группы можно использовать метод groups.getMembers (https://vk.com/dev/groups.getMembers). Для этого метода не нужна авторизация по токену ВК. Для пересечения пользователей в сообществах не стоит скачивать всех пользователей всех сообществ и пересекать эти множества. Так как это приведет к большому потреблению памяти. Нужно делать как в merge sort (https://ru.wikipedia.org/wiki/Сортировка_слиянием). То есть для каждой группы получаем порцию пользователей (задаем параметры ВК sort=id_asc, offset=0, count=1000). Далее для каждой группы заводим переменную которая указывает на текущий элемент. Изначально все переменные равны 0. Далее сравниваем id пользователей по текущим позициям в каждом из списков. Если нашли пользователя, который состоит в нужном нам количестве групп, сохраняем его в таблице результатов. Увеличиваем позицию в списке пользователей для группы где текущий id минимален. И т.д.
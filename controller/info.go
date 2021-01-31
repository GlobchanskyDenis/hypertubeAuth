package controller

import (
	"HypertubeAuth/logger"
	"net/http"
)

func info(w http.ResponseWriter, r *http.Request) {
	response := `<html><header><title>Hypertube endpoint info</title></header><body>
	<h2><b>Endpoints:</b></h2>

	<b>GET /api/auth/basic</b> - базовая авторизация пользователя</br>
	<b>GET /api/auth/oauth42</b> - делегированная авторизация пользователя через api 42</br>
	<b>GET /api/info</b> - информация об энпоинтах, детали реализации</br>
	<b>GET /api/profile/get</b> - возвращает твоего или чужого юзера (по его id)</br>

	<b>PUT /api/profile/create</b> - базовая регистрация пользователя пользователя</br>

	<b>PATCH /api/email/patch</b> - изменение почты пользователя</br>
	<b>PATCH /api/passwd/patch</b> - изменение пароля пользователя</br>
	<b>PATCH /api/profile/patch</b> - изменение остальных полей пользователя пользователя</br>

	<b>POST /api/email/confirm</b> - подтверждение почты пользователя</br>
	<b>POST /api/email/resend</b> - повторная отправка регистрационного письма на почту пользователя</br>
	<b>POST /api/passwd/repair</b> - отправка письма на почту пользователя с целью восстановления пароля</br>
	<b>POST /api/auth/check</b> - проверка, авторизирован ли пользователь. Для использования нужно знать пароль к серверу. (Это для Гриши)</br>

	<b>DELETE /api/profile/delete</b> - удаление пользователя</br>
	</br>
	<b>Поля пользователя</b> Отражены в структуре. Встречаются в запросах в разном составе</br>
	UserBasic {</br>
		&emsp;&emsp; userId	&emsp;&emsp;       integer</br>
		&emsp;&emsp; email      &emsp;&emsp;&emsp; string</br>
		&emsp;&emsp; passwd     &emsp;&emsp;       string</br>
		&emsp;&emsp; username   &emsp; string</br>
		&emsp;&emsp; firstName &emsp; string</br>
		&emsp;&emsp; lastName  &emsp; string</br>
	}</br>
	</br>

	<b>Обработка ошибок</b></br>
	В случае проваленного запроса сервер отвечает json в request body имеющим следующие поля:</br>
	description_ru - описание ошибки на русском. Можно показывать напрямую пользователю</br>
	description_eng - описание ошибки на английском. Можно показывать напрямую пользователю</br>
	code - код ошибки (тебе он врядли понадобится, это для моего тестирования)</br>
	Также сервер отвечает кодами, на которые тебе нужно реагировать:</br>
	400 BadRequest - в случае если ТЫ накосячил с запросом. Пользователь по идее не может привести к подобной ситуации</br>
	401 Unauthorized - сервер авторизации так ответить не может, но другие серверы так сообщают что нужно перебросить юзера на авторизацию</br>
	406 NotAcceptable - бизнес логика. Такое встречается если действие юзера невозможно из-за состояния БД (нет такого юзера и т.д)</br>
	422 UnprocessableEntity - в случае если юзер накосячил с аргументом запроса. Например мыло невалидное или пароль короткий</br>
	500 InternalServerError - может возникать в случае если например в конфигах неверная инфа (пороль БД, пароль почты...)</br>
	&emsp;&emsp;также возможны случаи, когда 500 возвращается в случае косяка в моем коде. Просьба сообщать мне о таких случаях</br>
	&emsp;&emsp;Отследить такое поведение просто - видишь 500 - смотришь что при этом напечатал в консоль мой логгер</br>
	&emsp;&emsp;и шлешь мне что ожидал, что получил и последние НЕСКОЛЬКО строк из логгера</br></br>

	<b>GET /api/auth/basic Базовая авторизация</b></br>
	Оформлена по стандарту RFC 2617, раздел 2. Идентификатор пользователя (его почта) и пароль передаются</br>
	одним заголовком, закодированным в base64. Схема:</br>
	Authorization: Basic base64( url_encoded_email <b>:</b> url_encoded_password ) </br>
	Если мыло = "Aladdin" а пароль = "open sesame", то финальный вид заголовка примет вид:</br>
	Authorization: Basic QWxhZGRpbjpvcGVuIHNlc2FtZQ==</br>
	В моем примере (папка client) есть рабочий вариант кода</br>
	Успешная авторизация - код 200 и json в request body содержащий одно поле - accessToken. Токен бессрочный</br>
	Провал авторизации: описанная выше ошибка</br></br>

	<b>PUT /api/profile/create Базовая регистрация пользователя</b></br>
	Обязательные поля для заполнения: email, passwd, username. Завернуть их в json и в тело запроса</br>
	Поля firstName, lastName, image_body задаются в эндпоинте PATCH /api/profile/patch</br>
	Успешная регистрация - код 200. В ответе есть тело полей пользователя, но это в основном для возможности тестирования</br>
	Провал регистрации: описанная выше ошибка</br></br>
	
	<b>GET /api/auth/oauth42 Делегированная авторизация через api 42</b></br>
	Этот эндпоинт не для фронта. Чтобы авторизароваться, нужно отправить GET запрос ИЗ ФОРМЫ по ссылочке, которую я дам отдельно</br>
	Эта ссылочка будет редиректить api 42 на данный эндпоинт, который будет отвечать тебе токеном в хидере</br>
	(в случае успеха) либо ошибкой в случае провала или отказа пользователя предоставлять права</br>
	Успех и провал отвечают редиректом. Ошибки ты не увидишь. Поэтому ты мне должен предоставить два эндпоинта</br>
	для редиректа. Успешный и провельных. Успешный пусть считывает хидер и сразу запрашивает поля юзера из </br>
	эндпоинта GET /api/profile/get (дожидаешься полной загрузки страницы, считываешь хидер и получаешь поля)</br>
	В случае провела пусть редирект ведет обратно на страницу авторизации.</br></br>

	<b>GET /api/profile/get - возвращение полей юзера</b></br>
	Эндпоинт доступен только авторизованным пользователям. Для авторизации нужно в заголовок accessToken вставить</br>
	авторизационный токен. Если аргумента в url нет - возвращает поля твоего юзера (идентифицирует по токену)</br>
	Если в url вставить аргумент userId=42 то вернутся поля юзера 42 (приватное поле email затрется)</br></br>

	<b>POST /api/email/recend - повторная отправка письма на почту для подтверждения почты</b></br>
	В теле заголовка в json должно быть запаковано одно поле email. Никаких больше деталей</br></br>

	<b>POST /api/email/confirm - подтверждение почты пользователя</b></br>
	В теле заголовка должно быть поле 

	<b>Эндпоинт проверки, авторизации пользователя</b></br>
	Предназначен только для внутренних запросов от других сервисов. Требуется знать пароль от сервера (безопасность)</br>
	В теле запроса нужно в json завернуть поля accessToken для моего авторизационного токена и server_passwd для пароля сервера</br>
	Три варианта ответа сервера. 200 - токен валиден, пользователь авторизован</br>
	401 - проверка подписи провалена, пользователь не авторизован</br>
	Все остальное - ошибки, которые следует обрабатывать отдельно (пароль сервера не тот, и т д)</br></br>

	</body></html>`

	successResponse(w, []byte(response))
	logger.Log(r, "programmer wants to know endpoints")
}

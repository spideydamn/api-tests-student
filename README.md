# API-тесты TestBoard

Позитивные и негативные тесты создания объявлений и работы с избранным

## Подготовка

Нужны Git, [Go 1.25.1 или выше](https://go.dev/doc/install) и [golangci-lint версии 2](https://golangci-lint.run/docs/welcome/install/local/), совместимый с установленным Go

```bash
git clone https://github.com/spideydamn/api-tests-student.git
cd api-tests-student
go mod download
```

Создайте два разных тестовых аккаунта на [TestBoard](https://testboard.avito.com/)
В корне проекта создайте `.env.override` и укажите их данные:

```dotenv
TEST_LOGIN=first_user_login
TEST_PASSWORD=first_user_password
FAVORITES_LOGIN=second_user_login
FAVORITES_PASSWORD=second_user_password
DEBUG=false
```

Первый аккаунт создаёт объявления, второй добавляет их в избранное

## Запуск

Из корня проекта:

```bash
go test -p 1 -count=1 -v ./...
golangci-lint run
```

`-p 1` запускает пакеты по очереди, чтобы ограничить нагрузку на стенд
Тесты сами создают и удаляют свои данные

[Swagger API](https://testboard.avito.com/swagger/index.html)

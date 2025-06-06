# Go Shop

## Текущее состояние:
Проект запущен на сервере, можно посомотреть по ссылке [go-shop.byBirukov.ru](http://go-shop.bybirukov.ru)
<br> Подключен Github Actions, автодеплой

**Монорепозиторий интернет-магазина с базовым функционалом.**  
Frontend, Backend и БД запускаются одной командой через `docker-compose`.

## Технологии

- **Backend:** Go 1.23, PostgreSQL
- **Frontend:** Vue 3, Vite, Nginx
- **DevOps:** Docker, Docker Compose


## Быстрый старт

### Предварительные требования

- Go 1.23+
- Docker & Docker Compose

### Запуск проекта

```bash
git clone https://github.com/alexbirbirdev/go-shop.git
cd go-shop
docker-compose up --build
```


### Просмотр
Frontend: http://localhost:3000
Backend API: http://localhost:8080


Встроенная поддержка CORS включена в backend

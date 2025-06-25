# Тестовое DeNet


## Для запуска использовать ```make run```. Сервер запускается в докер контейнерах, доступен по адресу ```http://localhost:8088```

## Доступные эндпоинты
### -POST /register - регистрация нового пользователя 
### -POST /login - аутентификация пользователя
### -GET /users/{id}/status - вся доступная информация о пользователе
### -GET /users/leaderboard - топ пользователей с самым большим балансом
### -POST /users/{id}/task/complete - выполнение задания 
### -POST /users/{id}/referrer - ввод реферального кода 

## Тестирование
### -POST /register
![image](https://github.com/user-attachments/assets/909fb895-93ea-4762-b68b-cac47c278480)

### -POST /login
![image](https://github.com/user-attachments/assets/d414cc84-ff92-4f6d-bc88-d05792f9fef1)

### -POST /users/6/task/complete
![image](https://github.com/user-attachments/assets/3d67adee-10b1-4107-8747-a0542683532e)

### -POST /users/6/referrer
![image](https://github.com/user-attachments/assets/f28ba806-5a2e-4a20-b062-d41ebe29a926)

### -GET /users/6/status
![image](https://github.com/user-attachments/assets/084174fc-9dd4-4b26-b4e8-f2fad48fdd26)

### -GET /users/leaderboard 
![image](https://github.com/user-attachments/assets/e2d7c5ba-7540-4b60-b817-380b6bd985aa)


## Обработка ошибок
### Попытка выполнить задачу повторно
![image](https://github.com/user-attachments/assets/70259c9d-40ed-4e73-b7fc-3cad3c42c470)

### Попытка зарегестировать пользователя с почтой, которая уже занята
![image](https://github.com/user-attachments/assets/f69b28b8-e792-430f-97bb-0bc0bddc48e5)

### Попытка аутентификации пользователя с неверным паролем
![image](https://github.com/user-attachments/assets/8a7356cc-24e7-4de0-893b-6466af99c1bb)

### Попытка получить доступ к не своему ресурсу
![image](https://github.com/user-attachments/assets/2af09fa9-3ab8-4c88-b952-24a2e4b58634)



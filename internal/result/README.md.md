# draw results service

## endpoints
примеры эндпоинтов для тестов
1. `GET /api/draws/{id}/results`: Получение выигрышной комбинации тиража.
```bash
curl --request GET \
  --url http://localhost:8080/api/draws/3/results \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.0.2'
```
пример ответа:
```JSON
{
	"result": [
		12,
		6,
		23,
		33,
		34
	]
}
```
2. `GET /api/tickets/{id}/check-result` (только USER): Проверка результата билета.
```bash
curl --request GET \
  --url http://localhost:8080/api/tickets/107/check-result \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.0.2'
```

пример ответа:
win_combination - выигрышная комбинация для данной лотереи
combination - комбинация данного билета
win_count - кол-во совпавших чисел комбинации билета и выигрышной комбинации лотереи.
```json
{
	"result": {
		"win_combination": [
			12,
			6,
			23,
			33,
			34
		],
		"combination": [
			22,
			6,
			30,
			3,
			9
		],
		"win_count": 1
	}
}
```

3.  `PUT /api/draws/{id}/results/generate (только ADMIN): генерация выигрышной комбинации для выбранной лотереи, идемпотентный. После него можно получить выигрышную комбинацию по `GET /api/draws/{id}/results`
```bash
curl --request PUT \
  --url http://localhost:8080/api/draws/3/results/generate \
  --header 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/11.0.2'
```
пример ответа:
```json
{
	"result": [
		12,
		6,
		23,
		33,
		34
	]
}
```
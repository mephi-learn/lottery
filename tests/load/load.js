import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Создаем метрику для отслеживания процента ошибок
const errorRate = new Rate('errors');

export const options = {
  // Настройка для тестирования 1000 запросов в секунду
  scenarios: {
    // Тест с постепенным увеличением нагрузки
    ramp_up_test: {
      executor: 'ramping-arrival-rate',
      startRate: 50,      // Начинаем с 50 RPS
      timeUnit: '1s',     // Единица времени - 1 секунда
      preAllocatedVUs: 100, // Предварительно выделенные VUs для начала
      maxVUs: 2000,       // Максимальное количество VUs (с запасом)
      stages: [
        { duration: '30s', target: 200 },  // Постепенное увеличение до 200 RPS
        { duration: '30s', target: 500 },  // Увеличение до 500 RPS
        { duration: '30s', target: 1000 }, // Достигаем целевых 1000 RPS
        { duration: '1m', target: 1000 },  // Поддерживаем 1000 RPS в течение 1 минуты
        { duration: '30s', target: 0 },    // Постепенное снижение нагрузки
      ],
    },
  },
  // Пороговые значения (thresholds) для проверки требований
  thresholds: {
    // Проверка, что 95% запросов выполняются быстрее 500 мс (0.5 секунды)
    http_req_duration: ['p(95)<500'],
    // Проверка, что коэффициент ошибок менее 1%
    errors: ['rate<0.01'],
  },
};

export default function () {
  // Задаём URL и заголовки (headers) из curl-запроса
  const url = 'http://localhost:8080/api/draws/active?=';
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': 'k6',
    },
  };

  // Выполняем GET-запрос
  const response = http.get(url, params);

  // Проверяем ответ и регистрируем ошибки
  const checkResult = check(response, {
    'статус 200': (r) => r.status === 200,
    'время ответа менее 500ms': (r) => r.timings.duration < 500,
  });

  // Увеличиваем счетчик ошибок, если хотя бы одна проверка не прошла
  errorRate.add(!checkResult);

  // Выводим в консоль детальную информацию при ошибках (опционально)
  if (!checkResult) {
    console.log(`Ошибка: Статус ${response.status}, Время: ${response.timings.duration}ms`);
  }

  // Без паузы между запросами для достижения максимальной нагрузки
}

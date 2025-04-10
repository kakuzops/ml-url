import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

const errorRate = new Rate('errors');

export const options = {
  scenarios: {
    high_load: {
      executor: 'ramping-vus',
      startVUs: 1000,
      stages: [
        { duration: '1m', target: 5000 },
        { duration: '2m', target: 10000 },
        { duration: '5m', target: 11000 },
        { duration: '10m', target: 11000 },
        { duration: '2m', target: 0 },
      ],
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500'],
    errors: ['rate<0.1'],
    http_reqs: ['rate>33000'],
  },
};

const BASE_URL = 'http://host.docker.internal:8080';

function generateRandomUrl() {
  const domains = ['example.com', 'test.com', 'demo.org'];
  const paths = ['page', 'article', 'product', 'service'];
  const domain = domains[Math.floor(Math.random() * domains.length)];
  const path = paths[Math.floor(Math.random() * paths.length)];
  const id = Math.floor(Math.random() * 1000);
  return `https://${domain}/${path}/${id}`;
}

export default function () {
  const urlToShorten = generateRandomUrl();
  const shortenResponse = http.post(`${BASE_URL}/shorten`, JSON.stringify({
    url: urlToShorten
  }), {
    headers: { 'Content-Type': 'application/json' },
  });

  check(shortenResponse, {
    'status is 201': (r) => r.status === 201,
    'has short_url': (r) => r.json('short_url') !== undefined,
  }) || errorRate.add(1);

  const shortUrl = shortenResponse.json('short_url');
  const shortCode = shortUrl.split('/').pop();

  const infoResponse = http.get(`${BASE_URL}/info/${shortCode}`);
  check(infoResponse, {
    'status is 200': (r) => r.status === 200,
    'has original_url': (r) => r.json('original_url') !== undefined,
  }) || errorRate.add(1);

  const redirectResponse = http.get(`${BASE_URL}/${shortCode}`, {
    redirects: 0,
  });
  check(redirectResponse, {
    'status is 302': (r) => r.status === 302,
  }) || errorRate.add(1);

  sleep(1);
} 
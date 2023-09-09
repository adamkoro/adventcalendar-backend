import http from "k6/http";

export const options = {
  scenarios: {
    constant_request_rate: {
      executor: 'constant-arrival-rate',
      rate: 200,
      timeUnit: '5s', // 1000 iterations per second, i.e. 1000 RPS
      duration: '10s',
      preAllocatedVUs: 10, // how large the initial pool of VUs would be
      maxVUs: 100, // if the preAllocatedVUs are not enough, we can initialize more
    },
  },
};

export default function () {
  const response = http.post("http://localhost:8080/api/login",'{"username":"admin","password":"admin"}');
  console.log(response.body);
}
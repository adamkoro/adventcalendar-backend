import http from "k6/http";

export const options = {
  iterations: 100,
};

export default function () {
  const response = http.post("http://localhost:8080/api/login",'{"username":"admin","password":"admin"}');
  console.log(response.body);
}
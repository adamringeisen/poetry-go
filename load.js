import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { duration: "20s", target: 350 },
    { duration: "20s", target: 500 },
    { duration: "20s", target: 1000 },
    { duration: "10s", target: 2000 },
  ],
};

export default function () {
  const res = http.get("http://localhost:8080/mc");
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}
// k6 run load.js

import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  stages: [
    { duration: "30s", target: 50 },
    { duration: "30s", target: 100 },
    { duration: "30s", target: 350 },
  ],
};

export default function () {
  const res = http.get("http://localhost:8080/mc");
  check(res, { "status was 200": (r) => r.status == 200 });
  sleep(1);
}
// k6 run load.js

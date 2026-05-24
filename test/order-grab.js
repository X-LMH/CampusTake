import http from 'k6/http';
import { Counter } from 'k6/metrics';

import { BASE_URL, RIDER } from './config.js';

const successGrab = new Counter('success_grab');
const failedGrab = new Counter('failed_grab');

export const options = {
    vus: 100,
    duration: '1s',
};

export function setup() {

    const payload = JSON.stringify({
        phone: RIDER.phone,
        password: RIDER.password
    });

    const res = http.post(
        `${BASE_URL}/api/login/password`,
        payload,
        {
            headers: {
                'Content-Type': 'application/json',
            },
        }
    );

    const token = JSON.parse(res.body).data.token;

    return { token };
}

export default function (data) {

    const payload = JSON.stringify({
        order_id: 18
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${data.token}`,
        },
    };

    const res = http.post(
        `${BASE_URL}/api/rider/order/grab`,
        payload,
        params
    );

    const body = JSON.parse(res.body);

    if (body.code === 0) {
        successGrab.add(1);
    } else {
        failedGrab.add(1);
    }
}
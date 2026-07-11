import http from 'k6/http';

import { USER } from './config.js';
import { BASE_URL } from './config.js';
import { login } from './login.js';

export const options = {
    vus: 100,
    duration: '10s',
};

export function setup() {

    const token = login(
        USER.phone,
        USER.password
    );

    return { token };
}

export default function (data) {

    const params = {
        headers: {
            Authorization: `Bearer ${data.token}`,
        },
    };

    http.get(
        `${BASE_URL}/api/order/list?page=1&size=10`,
        params
    );
}
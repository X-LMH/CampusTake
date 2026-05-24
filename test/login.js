import http from 'k6/http';
import { BASE_URL } from './config.js';

export function login(phone, password) {

    const payload = JSON.stringify({
        phone: phone,
        password: password
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
        },
    };

    const res = http.post(
        `${BASE_URL}/api/login/password`,
        payload,
        params
    );

    console.log('login response:', res.body);

    const body = JSON.parse(res.body);

    return body.data.token;
}
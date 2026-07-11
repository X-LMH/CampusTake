import http from 'k6/http';

import { BASE_URL, USER } from './config.js';

export const options = {
    vus: 100,
    duration: '10s',
};

export function setup() {

    // 登录获取 token
    const loginPayload = JSON.stringify({
        phone: USER.phone,
        password: USER.password
    });

    const loginRes = http.post(
        `${BASE_URL}/api/login/password`,
        loginPayload,
        {
            headers: {
                'Content-Type': 'application/json',
            },
        }
    );

    const token = JSON.parse(loginRes.body).data.token;

    return { token };
}

export default function (data) {

    const payload = JSON.stringify({
        order_type: 1,
        pickup_address_id: 2,
        delivery_address_id: 1,
        reward_amount: 19.9,
        remark: "无"
    });

    const params = {
        headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${data.token}`,
        },
    };

    const res = http.post(
        `${BASE_URL}/api/order/create`,
        payload,
        params
    );

    // 如果不是业务成功，打印出来
    const body = JSON.parse(res.body);

    if (body.code !== 0) {
        console.log(res.body);
    }
}
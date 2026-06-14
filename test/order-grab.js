import http from 'k6/http';
import { Counter } from 'k6/metrics';
import { BASE_URL, RIDER, USER } from './config.js';

const successGrab = new Counter('success_grab');
const failedGrab = new Counter('failed_grab');

export const options = {
    vus: 100,
    duration: '3s',
};

// 前置：登录用户 + 骑手，创建100个【已支付】可抢订单
export function setup() {
    // 1. 用户登录（用来创建、支付订单）
    const userRes = http.post(
        `${BASE_URL}/api/login/password`,
        JSON.stringify({
            phone: USER.phone,
            password: USER.password
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    const userToken = JSON.parse(userRes.body).data.token;

    // 2. 骑手登录（用来抢单）
    const riderRes = http.post(
        `${BASE_URL}/api/login/password`,
        JSON.stringify({
            phone: RIDER.phone,
            password: RIDER.password
        }),
        { headers: { 'Content-Type': 'application/json' } }
    );
    const riderToken = JSON.parse(riderRes.body).data.token;

    // 3. 批量创建 100 个已支付订单（完全按你给的API）
    const orderIds = [];
    for (let i = 0; i < 100; i++) {
        // ====================== 创建订单（完全匹配你的结构体）======================
        const createRes = http.post(
            `${BASE_URL}/api/order/create`,
            JSON.stringify({
                order_type: 1,
                pickup_address_id: 2,
                delivery_address_id: 1,
                reward_amount: 19.9,
                remark: "无"
            }),
            {
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${userToken}`
                }
            }
        );

        // 你返回的是 order_id，直接取
        const createData = JSON.parse(createRes.body);
        const orderId = createData.data.order_id;

        // ====================== 支付订单 ======================
        http.post(
            `${BASE_URL}/api/order/pay`,
            JSON.stringify({ order_id: orderId }),
            {
                headers: {
                    'Content-Type': 'application/json',
                    Authorization: `Bearer ${userToken}`
                }
            }
        );

        orderIds.push(orderId);
    }

    return { riderToken, orderIds };
}

// ====================== 高并发抢单（每个VU抢独立订单）======================
export default function (data) {
    // 每个VU抢自己的订单，永不冲突
    const orderId = data.orderIds[(__VU - 1) % data.orderIds.length];

    const res = http.post(
        `${BASE_URL}/api/rider/order/grab`,
        JSON.stringify({ order_id: orderId }),
        {
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${data.riderToken}`
            }
        }
    );

    const body = JSON.parse(res.body);
    body.code === 0 ? successGrab.add(1) : failedGrab.add(1);
}
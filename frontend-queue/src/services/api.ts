import type { Queue } from "../types/queue"
import axios from "axios";

const API_BASE = `http://localhost:3000/api`

export async function getQueues(): Promise<Queue[]> {
  const res = await axios.get(`${API_BASE}/queue/`);
  
  if (res.status !== 200) {
    throw new Error("Failed to fetch queues");
  }

  return res.data as Queue[]; // ไม่ต้องใส่ () เพราะเป็น property
}

export async function createQueueWithQrcode(name:string,phone:string):Promise<Queue> {
  const res = await axios.post(`${API_BASE}/queue/qr`, { name,phone})  
  return res.data.queue as Queue;
}
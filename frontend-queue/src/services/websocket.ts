import { ref } from "vue";

const ws = ref<WebSocket | null>(null);

const currentQueue = ref(0);
const totalQueue = ref(0);

export function connectQueueWebSocket(queueNumber: string , onUpdate: () => void) {
  ws.value = new WebSocket(
    // `wss://localhost:3000/queue?queue_no=${queueNumber}`
    `ws://localhost:3000/queue?queue_no=${queueNumber}`
  );

  ws.value.onopen = () => console.log("WebSocket connected");

    ws.value.onmessage = (event) => {
    console.log("Raw event:", event.data);
    try {
        const data = JSON.parse(event.data);
        currentQueue.value = Number(data.current_queue) || 0;
        totalQueue.value = Number(data.total_queue) || 0;
        onUpdate();
    } catch (e) {
        console.error("JSON parse error:", e);
    }
};

  ws.value.onerror = (err) => console.error("WebSocket error:", err);

  ws.value.onclose = () => console.log("WebSocket disconnected");
}

export function disconnectQueueWebSocket(){
    if (ws.value){
        ws.value.close();
        ws.value = null;
    }
}

export {currentQueue,totalQueue}

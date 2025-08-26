<script setup lang="ts">
import { ref, onUnmounted, watch } from "vue";
import { createQueueWithQrcode } from "../services/api";
import {
  disconnectQueueWebSocket,
  connectQueueWebSocket,
  totalQueue,
  currentQueue,
} from "../services/websocket";

import QueueForm from "../components/QueueForm.vue";
import Queuedetail from "../components/Queuedetail.vue";

const queueNumber = ref<string | null>(null);
const qrImage = ref<string | null>(null);
const error = ref<string | null>(null);

const queueProgress = ref(0);
const remainingQueue = ref(0);

// Subscribe currentQueue/totalQueue reactive จาก service
watch([currentQueue, totalQueue, queueNumber], () => {
  if (!queueNumber.value) return;

  const queueNo = Number(queueNumber.value) || 0;
  const current = Number(currentQueue.value) || 0;
  const total = Number(totalQueue.value) || 0;

  remainingQueue.value = Math.max(queueNo - current, 0);
  queueProgress.value = total > 0 ? (current / total) * 100 : 0;

  if (remainingQueue.value === 0 && queueNo > 0) {
    alert("ถึงคิวของคุณแล้ว!");
  }
});

async function handleFormSubmit({ name , phone }: { name: string , phone:string}) {
  error.value = null;
  queueNumber.value = null;
  try {
    const queue = await createQueueWithQrcode(name,phone);
    queueNumber.value = queue.queue_no;
    qrImage.value = queue.qr_code ? `data:image/png;base64,${queue.qr_code}` : null;
    if (queueNumber.value) connectQueueWebSocket(queueNumber.value, () => {});
  } catch {
    error.value = "เกิดข้อผิดพลาด";
  }
}

// Cleanup on unmount
onUnmounted(() => {
  disconnectQueueWebSocket();
});
</script>
<template>
  <div
    class="flex flex-col items-center justify-center min-h-screen bg-gray-100"
  >
    <div class="bg-white shadow-lg rounded-2xl p-6 w-full max-w-md">
      <h1 class="text-2xl font-bold text-center text-blue-600 mb-4">
        Queue App
      </h1>

      <!--Queue Form-->
      <QueueForm @submit="handleFormSubmit"/>
  
      <!-- Queue Info -->
      <Queuedetail
        :queueNumber="queueNumber"
        :queueProgress="queueProgress"
        :remainingQueue="remainingQueue"
        :qrImage="qrImage"
        :error="error"
      />
    </div>
  </div>
</template>

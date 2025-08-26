<script setup lang="ts">
import { defineProps } from "vue";

const props = defineProps<{
  queueNumber: string | null;
  queueProgress: number;
  remainingQueue: number;
  qrImage: string | null;
  error: string | null;
}>();
</script>

<template>
  <div v-if="props.queueNumber" class="mt-6 text-center">
    <p class="text-xl font-bold text-green-600 mb-2">
      🎉 คุณได้รับคิวที่ {{ props.queueNumber }}
    </p>

    <!-- Progress Tracker -->
    <div class="w-full bg-gray-200 h-4 rounded-full overflow-hidden mb-2">
      <div
        class="bg-blue-500 h-4 rounded-full transition-all"
        :style="{ width: props.queueProgress + '%' }"
      ></div>
    </div>
    <p class="text-sm text-gray-600">
      อีก {{ props.remainingQueue }} คิวถึงคิวของคุณ
    </p>

    <!-- QR Code -->
    <img
      v-if="props.qrImage"
      :src="props.qrImage"
      alt="QR Code"
      class="mx-auto mt-4 w-40 h-40"
    />
  </div>

  <p v-if="props.error" class="mt-4 text-center text-red-600">{{ props.error }}</p>
</template>
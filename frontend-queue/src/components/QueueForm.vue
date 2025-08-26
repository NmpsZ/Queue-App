<template>
  <form @submit.prevent="handleSubmit" class="space-y-4">
    <input
      v-model="name"
      type="text"
      placeholder="Name"
      class="w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-400"
    />
    <input
      v-model="phone"
      type="tel"
      placeholder="Phone"
      class="w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-400"
    />
    <button
      type="submit"
      class="w-full bg-blue-500 text-white font-semibold p-3 rounded-lg hover:bg-blue-600 transition"
      :disabled="loading"
    >
      {{ loading ? "กำลังจอง..." : "จองคิว" }}
    </button>
  </form>
</template>

<script setup lang="ts">
import { ref } from "vue";

const name = ref<string>("");
const phone = ref<string>("");
const loading = ref(false);

const emit = defineEmits<{
  (e: 'submit', payload: { name: string; phone: string }): void
}>();

function handleSubmit() {
  if (!name.value.trim()) return;
  loading.value = true;
  emit("submit", { name: name.value, phone: phone.value });
  loading.value = false;
}
</script>
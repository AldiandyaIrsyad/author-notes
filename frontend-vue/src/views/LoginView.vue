<template>
  <div class="bg-gray-100 flex items-center justify-center min-h-screen w-full">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md md:max-w-lg">
      <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">
        Author Notes
      </h1>
      <h2 class="text-xl font-semibold text-center text-gray-700 mb-8">
        Login to your account
      </h2>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label
            for="username"
            class="block text-gray-700 text-sm font-medium mb-2"
            >Username</label
          >
          <input
            type="text"
            id="username"
            v-model="username"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-gray-900"
            placeholder="Enter your username"
          />
        </div>

        <div class="mb-6">
          <label
            for="password"
            class="block text-gray-700 text-sm font-medium mb-2"
            >Password</label
          >
          <input
            type="password"
            id="password"
            v-model="password"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-gray-900"
            placeholder="Enter your password"
          />
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full bg-indigo-600 text-white py-2 px-4 rounded-md font-semibold hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition duration-200 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isLoading">Logging in...</span>
          <span v-else>Login</span>
        </button>
      </form>
      <p class="text-center text-gray-600 text-sm mt-6">
        Don't have an account?
        <router-link to="/register" class="text-indigo-600 hover:underline"
          >Register here</router-link
        >
      </p>
    </div>
    <!-- Basic message display -->
    <div
      v-if="message"
      :class="
        messageType === 'success'
          ? 'bg-green-100 text-green-800'
          : 'bg-red-100 text-red-800'
      "
      class="fixed top-5 left-1/2 transform -translate-x-1/2 p-4 rounded-md shadow-md"
    >
      {{ message }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { authService } from "@/services/auth"; 
import axios from "axios"; 

const username = ref(""); 
const password = ref("");
const message = ref("");
const messageType = ref<"success" | "error" | "">("");
const isLoading = ref(false); // Add loading state
const router = useRouter();

const showMessage = (msg: string, type: "success" | "error") => {
  message.value = msg;
  messageType.value = type;
  setTimeout(() => {
    message.value = "";
    messageType.value = "";
  }, 3000);
};

const handleLogin = async () => {
  if (!username.value || !password.value) {
    showMessage("Please enter both username and password.", "error");
    return;
  }

  isLoading.value = true;
  message.value = "";
  messageType.value = "";

  try {
    const payload = {
      username: username.value,
      password: password.value,
    };
    const response = await authService.login(payload);

    // Handle success - Store token (e.g., localStorage) and redirect
    localStorage.setItem("authToken", response.token); 
    showMessage("Login successful! Redirecting...", "success");
    isLoading.value = false;

    // Redirect after a short delay
    setTimeout(() => {
      router.push("/dashboard"); 
    }, 1500);
  } catch (error) {
    isLoading.value = false;
    console.error("Login failed:", error);
    let errorMessage = "Login failed. Please try again.";

    if (axios.isAxiosError(error) && error.response) {
      // Use error message from backend if available
      errorMessage = error.response.data?.error || errorMessage;
      if (error.response.status === 401) {
        errorMessage = "Invalid username or password.";
      }
    }
    showMessage(errorMessage, "error");
  }
};
</script>

<style scoped>
/* Add any component-specific styles here if needed */
/* Using Tailwind classes directly in the template is preferred */
</style>

<template>
  <div class="bg-gray-100 flex items-center justify-center min-h-screen">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md md:max-w-lg">
      <h1 class="text-2xl font-bold text-center text-gray-800 mb-6">
        Author Notes
      </h1>
      <h2 class="text-xl font-semibold text-center text-gray-700 mb-8">
        Create your account
      </h2>

      <form @submit.prevent="handleRegister">
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
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            placeholder="Choose a username"
          />
        </div>

        <div class="mb-4">
          <label
            for="email"
            class="block text-gray-700 text-sm font-medium mb-2"
            >Email</label
          >
          <input
            type="email"
            id="email"
            v-model="email"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            placeholder="Enter your email address"
          />
        </div>

        <div class="mb-4">
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
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            placeholder="Create a password"
          />
        </div>

        <div class="mb-6">
          <label
            for="confirmPassword"
            class="block text-gray-700 text-sm font-medium mb-2"
            >Confirm Password</label
          >
          <input
            type="password"
            id="confirmPassword"
            v-model="confirmPassword"
            required
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
            placeholder="Confirm your password"
          />
        </div>

        <button
          type="submit"
          :disabled="isLoading"
          class="w-full bg-indigo-600 text-white py-2 px-4 rounded-md font-semibold hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition duration-200 ease-in-out disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="isLoading">Registering...</span>
          <span v-else>Register</span>
        </button>
      </form>
      <p class="text-center text-gray-600 text-sm mt-6">
        Already have an account?
        <router-link to="/login" class="text-indigo-600 hover:underline"
          >Login here</router-link
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
import { authService } from "@/services/auth"; // Import the auth service
import axios from "axios"; // Import axios to check for AxiosError

const username = ref("");
const email = ref("");
const password = ref("");
const confirmPassword = ref("");
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

const handleRegister = async () => {
  // Basic validation
  if (
    !username.value ||
    !email.value ||
    !password.value ||
    !confirmPassword.value
  ) {
    showMessage("Please fill in all fields.", "error");
    return;
  }
  if (password.value !== confirmPassword.value) {
    showMessage("Passwords do not match.", "error");
    return;
  }

  isLoading.value = true;
  message.value = "";
  messageType.value = "";

  try {
    const payload = {
      username: username.value,
      email: email.value,
      password: password.value,
    };
    await authService.register(payload);

    // Handle success
    showMessage("Registration successful! Redirecting to login...", "success");
    isLoading.value = false;

    // Redirect to login after a delay
    setTimeout(() => {
      router.push("/login");
    }, 1500);
  } catch (error) {
    isLoading.value = false;
    console.error("Registration failed:", error);
    let errorMessage = "Registration failed. Please try again.";
    // Check if it's an axios error with a response
    if (axios.isAxiosError(error) && error.response) {
      // Use error message from backend if available
      errorMessage = error.response.data?.error || errorMessage;
      if (error.response.status === 409) {
        // Conflict - User already exists
        errorMessage = "Username or email already exists.";
      }
    }
    showMessage(errorMessage, "error");
  }
};
</script>

<style scoped>
/* Add any component-specific styles here if needed */
</style>

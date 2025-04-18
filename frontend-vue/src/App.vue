<template>
  <div id="app-container" class="font-sans">
    <router-view /> <!-- This will render the matched route component (LoginView or RegisterView) -->
  </div>
</template>

<script setup lang="ts">
// No script needed for this basic setup, unless you add global logic/layout elements
</script>

<style>
/* Global styles can remain in style.css or be added here if needed */
/* Ensure body/html takes full height if needed for layout */
html, body, #app, #app-container {
  height: 100%;
  margin: 0;
  padding: 0;
}

/* Apply Inter font globally if not already handled by Tailwind config or style.css */
body {
  font-family: 'Inter', sans-serif;
}
</style>

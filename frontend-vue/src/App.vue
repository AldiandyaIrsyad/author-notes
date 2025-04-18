<template>
  <div id="app-container" class="font-sans">
    <router-view />
  </div>
</template>

<script setup lang="ts"></script>

<style>
html,
body,
#app,
#app-container {
  height: 100%;
  margin: 0;
  padding: 0;
}

body {
  font-family: "Inter", sans-serif;
}
</style>

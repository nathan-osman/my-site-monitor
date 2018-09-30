<style>
body {
  font-family: sans-serif;
  margin: 0;
}

.container {
  margin: 0 16px;
}

@media (min-width: 992px) {
  .container {
    margin: auto;
    max-width: 970px;
  }
}
</style>

<template>
  <div>
    <AppHeader />

    <div class="container">
      <router-view />
    </div>
  </div>
</template>

<script>
import AppHeader from './AppHeader'

export default {
  name: 'App',
  components: {
    AppHeader
  }
}
</script>

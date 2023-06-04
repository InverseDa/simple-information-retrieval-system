<template>
  <div class="hello">
    <a-input-search class="my-input"
                    v-model:value="search_value" 
                    placeholder="请输入要搜索的东西~"
                    size="large"
                    :loading="is_loading"
                    enter-button/>
  </div>
</template>

<script setup>
import {ref} from 'vue'
import axios from 'axios'
const search_value = ref('')
const is_loading = false
axios.defaults.baseURL = 'http://' + window.location.hostname + ':8000'

</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.my-input {
  width: 500px;
  margin: 0 auto;
}
</style>

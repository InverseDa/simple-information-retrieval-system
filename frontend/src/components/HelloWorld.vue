<template>
  <div class="hello">
    <a-input-search class="my-input"
                    v-model:value="search_value" 
                    placeholder="请输入要搜索的东西~"
                    size="large"
                    :loading="is_loading"
                    enter-button
                    @search="searching"/>
  </div>
</template>

<script setup>
import {ref} from 'vue'
import axios from 'axios'
const search_value = ref('')
const is_loading = ref(false)

const searching = () => {
  is_loading.value = true
  axios.post('http://localhost:8888/api/query', {
    query: search_value.value
  })
    .then(res => {
      console.log(res)
      is_loading.value = false
    })
    .catch(err => {
      console.log(err)
      is_loading.value = false
    })
}

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

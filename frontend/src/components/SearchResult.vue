<template>
  <Suspense>
    <div>
      <a-input-search class="my-input" v-model:value="search_value" placeholder="请输入要搜索的东西~" size="large"
        :loading="is_loading" enter-button @search="searching" />
      {{ listData }}
      <!-- <a-list item-layout="vertical" size="large" :pagination="pagination" :data-source="listData">
        <template #footer>
          <div>
            <b>Powered by miaokeda.</b>
          </div>
        </template>
        <template #renderItem="{ item }">
          <a-list-item key="item.title">
            <template #actions>
              <span v-for="{ type, text } in actions" :key="type">
                <component :is="type" style="margin-right: 8px;" />
                {{ text }}
              </span>
            </template>
            <template #extra>
              <img width="272" alt="logo" src="https://gw.alipayobjects.com/zos/rmsportal/mqaQswcyDLcXyDKnZfES.png" />
            </template>
            <a-list-item-meta :description="item.description">
              <template #title>
                <a :href="item.href">{{ item.title }}</a>
              </template>
            </a-list-item-meta>
            {{ item.content }}
          </a-list-item>
        </template>
      </a-list> -->
    </div>
  </Suspense>
</template>

<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router';
// import { StarOutlined, LikeOutlined, MessageOutlined } from '@ant-design/icons-vue';
import axios from 'axios'

const search_value = ref('')
const is_loading = ref(false)
const route = useRoute()
const listData =ref([])
if (route.params.query)
  search_value.value = JSON.parse(route.params.query)

// const pagination = {
//   onChange: page => {
//     console.log(page);
//   },
//   pageSize: 3,
// }

// const actions = [{
//   type: 'StarOutlined',
//   text: '156',
// }, {
//   type: 'LikeOutlined',
//   text: '156',
// }, {
//   type: 'MessageOutlined',
//   text: '2',
// }]

const searching = () => {
  is_loading.value = true
  deal_axios()
}

const deal_axios = async () => {
  try {
    const res = await axios.post('http://localhost:8888/api/query', {
      query: search_value.value
    })
    listData.value=[] 
    for (let i = 0; i < res.data.pagesString.length; i++) {
      listData.value.push({
        href: 'https://www.antdv.com/',
        title: res.data.pagesString[i].title,
        description: 'Ant Design, a design language for background applications, is refined by Ant UED Team.',
        content: res.data.pagesString[i].content,
      })
    }
    console.log(listData.value)
    is_loading.value = false
  } catch (err) {
    console.log(err)
    is_loading.value = false;
  }
}

if (search_value.value !== "") {
  await deal_axios()
  console.log("YES")
}

</script>

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
  text-align: left;
}

.my-input {
  width: 500px;
  margin: 0 auto;
}

a-list {
  background-color: #f0f0f0;
  padding: 20px;
}

a-list-item {
  border: 1px solid #ccc;
  margin-bottom: 10px;
  padding: 10px;
}

a-list-item-meta-title {
  font-size: 18px;
  font-weight: bold;
  text-align: left;
}

a-list-item-meta-description {
  font-size: 14px;
}
</style>
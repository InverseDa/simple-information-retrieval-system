<template>
  <div>
    <a-input-search
      class="my-input"
      v-model:value="search_value"
      placeholder="请输入要搜索的东西~"
      size="large"
      :loading="is_loading"
      enter-button
      @search="searching"
    />
    <div class="my-list">
      <a-list
        item-layout="vertical"
        size="large"
        :pagination="pagination"
        :data-source="listData"
      >
        <template #footer>
          <div>
            <b>Powered by miaokeda.</b>
          </div>
        </template>
        <template #renderItem="{ item }">
          <a-list-item key="item.title">
            <template #extra>
              <img
                width="272"
                alt="logo"
                src="https://gw.alipayobjects.com/zos/rmsportal/mqaQswcyDLcXyDKnZfES.png"
              />
            </template>
            <a-list-item-meta :description="item.description">
              <template #title>
                <a :href="item.href">{{ item.title }}</a>
              </template>
            </a-list-item-meta>
            {{ item.content }}
          </a-list-item>
        </template>
      </a-list>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import { useRoute } from "vue-router";
// import { StarOutlined, LikeOutlined, MessageOutlined } from '@ant-design/icons-vue';
import axios from "axios";

const search_value = ref("");
const is_loading = ref(false);
const route = useRoute();
const listData = ref([]);

if (route.params.query) {
  search_value.value = route.params.query;
  searching();
}

const pagination = {
  onChange: (page) => {
    console.log(page);
  },
  pageSize: 3,
};

async function searching() {
  is_loading.value = true;
  await deal_axios();
  is_loading.value = false;
}

async function deal_axios() {
  try {
    const res = await axios.post("http://localhost:8888/api/query", {
      query: search_value.value,
    });
    listData.value = [];
    for (let i = 0; i < res.data.pagesString.length; i++) {
      listData.value.push({
        href: "https://www.antdv.com/",
        title: res.data.pagesString[i].title,
        description:
          "Ant Design, a design language for background applications, is refined by Ant UED Team.",
        content: res.data.pagesString[i].content,
      });
    }
    console.log(listData.value);
  } catch (err) {
    console.log(err);
  }
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

.my-list {
  width: 800px;
  margin: 0 auto;
}
</style>

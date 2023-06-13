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
    <div>
      <a></a>
    </div>
    <div class="my-list">
      <br />
      <br />
      <template v-for="(pages, index) in listData" :key="index">
        <a-descriptions :title="`Result ${index + 1}`" bordered>
          <a-descriptions-item
            label="Title"
            style="text-align: center"
            :span="3"
          >
            <a :href="pages.href" style="color: #42b983; text-align: left">{{
              pages.title
            }}</a>
          </a-descriptions-item>
          <a-descriptions-item label="Content" style="text-align: center">
            <pre class="content" style="text-align: left">
              {{ pages.content }}
            </pre>
          </a-descriptions-item>
        </a-descriptions>
        <br />
        <br />
      </template>
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
const fuzzyData = ref([]);

if (route.params.query) {
  search_value.value = route.params.query;
  searching();
}

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
    console.log(res.data);
    if (res.data.status === "success") {
      listData.value = [];
      for (let i = 0; i < res.data.pagesString.length; i++) {
        listData.value.push({
          href: res.data.pagesString[i].url,
          title: res.data.pagesString[i].title.replace(" ", ""),
          // description:
          //   "Ant Design, a design language for background applications, is refined by Ant UED Team.",
          content: res.data.pagesString[i].content
            .replace(
              /(\d+|[一二三四五六七八九十]+|（\d+|[一二三四五六七八九十]+）|\d+\.)(?=([^0-9]|$))(◆|▲|●)/g,
              "$1\n$2"
            )
            .replace(/ {4,}/g, "\n"),
        });
      }
    } else {
      fuzzyData.value = [];
      listData.value = [];
      for (let i = 0; i < res.data.fuzzySearchString.length; i++) {
        fuzzyData.value.push(res.data.fuzzySearchString[i]);
        console.log(fuzzyData.value[i]);
      }
    }
  } catch (err) {
    console.log(err);
  }
}
</script>

<style scoped>
.my-input {
  width: 500px;
  margin: 0 auto;
}
.content {
  white-space: pre-line;
}
.my-list {
  width: 1000px;
  margin: 0 auto;
}
</style>

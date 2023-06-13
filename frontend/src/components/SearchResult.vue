<template>
  <div>
    <div class="my-input">
      <a-input-search
        v-model:value="search_value"
        placeholder="请输入要搜索的东西~"
        size="large"
        :loading="is_loading"
        enter-button
        @search="press"
      />
      <p v-if="listData.length !== 0" style="text-align: left">
        搜索耗时：{{ parseInt(searchTime.valueOf()) + "ms" }}
      </p>
      <p v-if="fuzzyData.length !== 0" style="text-align: left">
        您是不是在找：
        <a
          v-for="(item, index) in fuzzyData"
          :key="index"
          :href="`/result/${item}`"
          style="color: red"
        >
          {{ item + "\t" }}
        </a>
      </p>
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
import { useRouter } from "vue-router";

const search_value = ref("");
const is_loading = ref(false);
const route = useRoute();
const router = useRouter();
const listData = ref([]);
const fuzzyData = ref([]);
const searchTime = ref(0);

if (route.params.query) {
  search_value.value = route.params.query;
  searching();
}

async function press() {
  fuzzyData.value = [];
  const keyword = search_value.value;

  if (keyword) {
    // 使用$router.push()方法实现路由跳转
    router.replace({ name: "Result", params: { query: keyword } });
  }
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
      searchTime.value = res.data.time * 1000;
      for (let i = 0; i < res.data.pagesString.length; i++) {
        let str = res.data.pagesString[i].content
          .replace(
            /(\d+|[一二三四五六七八九十]+|（\d+|[一二三四五六七八九十]+）|\d+\.)(?=([^0-9]|$))(◆|▲|●)/g,
            "$1\n$2"
          )
          .replace(/ {4,}/g, "\n");
        listData.value.push({
          href: res.data.pagesString[i].url,
          title: res.data.pagesString[i].title.replace(" ", ""),
          content: str,
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
.highlight {
  background-color: yellow;
}
</style>

<template>
  <a :href="getUrl" target="_blank" rel="noopener noreferrer">
    <main>
      <h4>{{ getTitle }}</h4>
      <p>{{ getDescription }}</p>
    </main>
    <img src="https://avatars3.githubusercontent.com/u/34158802?s=400&v=4">
  </a>
</template>

<script>
export default {
  props: {
    url: {
      type: String,
      required: true,
    },
    title: {
      type: String
    },
    description: {
      type: String
    }
  },
  computed: {
    getTitle () {
      return this.title || this.$page.title || this.$page.frontmatter.title || ""
    },
    getDescription () {
      return this.description || this.$page.description || this.$page.frontmatter.description || ""
    },
    getUrl () {
      return this.url || this.$page.url || this.$page.frontmatter.url || ""
    }
  }
}
</script>

<style lang="stylus" scoped>
a {
  border: solid 1px #eaecef;
  display: flex;
  justify-content: space-between;
  height: 150px;
  color: #2c3e50;
  box-shadow: 0 0 4px 0 #eaecef;
}
a:hover {
  border-color: #3eaf7c;
  text-decoration: none!important;
  box-shadow: 0 0 10px 0 #eaecef;
}
a:hover h4 {
  color: #3eaf7c;
}
main {
  padding: 1em;
}
h4 {
  margin-top: 0;
}
p {
  color: lighten(#2c3e50, 40);
  font-size: .9em;
}
img {
  height: 150px;
  width: 150px;
}
</style>

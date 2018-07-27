<template>
  <div class="tabs">
    <nav>
      <a v-for="(tab, i) in tabs" :key="i"
        :class="{ 'active': tab.isActive }"
        @click="selectTab(tab.id)">
        {{ tab.title }}
      </a>
    </nav>
    <div class="tabs-content" ref="tabs">
      <slot></slot>
    </div>
  </div>
</template>

<script>
export default {
  data: () => ({
    tabs: [],
    activeId: ''
  }),
  mounted() {
    this.tabs = this.$children
    if (!this.tabs.length) return
    this.selectTab(this.tabs[0].id)
  },
  methods: {
    findTab(id) {
      return this.tabs.find(tab => tab.id === id)
    },
    selectTab(selectedTabId) {
      const selectedTab = this.findTab(selectedTabId)
      if (!selectedTab) return
      this.tabs.forEach(tab => {
        tab.isActive = (tab.id === selectedTab.id)
      })
      this.activeId = selectedTab.id
    }
  }
}
</script>

<style lang="styl" scoped>
.tabs {
  margin: 1em 0;
}
nav {
  display: flex;
  justify-content: flex-start;
  padding: 0;
  margin: 0;
  transform: translateY(1px);
}
a {
  border: solid 1px #eaecef;
  z-index: 2;
  border-radius: 3px 3px 0 0;
  list-style: none;
  margin-right: -1px;
  background: rgba(0,0,0,0.015);
  cursor: pointer;
  padding: .5em 1em;
}
a.active {
  z-index: 2;
  border-bottom-color: white;
  background: white;
}
.tabs-content {
  border: solid 1px #eaecef;
  padding: 1em;
  border-radius: 0 3px 3px 3px;
}
</style>

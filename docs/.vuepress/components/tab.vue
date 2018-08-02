<template>
  <section v-show="isActive">
    <slot></slot>
  </section>
</template>


<script>
export default {
  props: {
    title: {
      type: String,
      required: true
    }
  },
  data: () => ({
    isActive: false
  }),
  computed: {
    id() {
      return this.title.toLowerCase().replace(/ /g, '-')
    }
  }
}
</script>

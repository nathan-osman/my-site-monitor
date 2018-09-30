import Vue from 'vue';
import Router from 'vue-router'
import Home from './components/Home'
import Site from './components/Site'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Home',
      component: Home
    },
    {
      path: '/sites/:id',
      name: 'Site',
      component: Site
    }
  ]
})

import { createRouter, createWebHistory } from 'vue-router'

export const router = createRouter({
  history: createWebHistory(),
  scrollBehavior: (_to, _from, saved) => saved ?? { top: 0 },
  routes: [
    { path: '/', name: 'home', component: () => import('@/views/HomeView.vue') },
    { path: '/people', name: 'people', component: () => import('@/views/PeopleView.vue') },
    { path: '/people/:slug', name: 'person', component: () => import('@/views/PersonDetailView.vue'), props: true },
    { path: '/works', name: 'works', component: () => import('@/views/WorksView.vue') },
    { path: '/works/:slug', name: 'work', component: () => import('@/views/WorkDetailView.vue'), props: true },
    { path: '/timeline', name: 'timeline', component: () => import('@/views/TimelineView.vue') },
    { path: '/glossary', name: 'glossary', component: () => import('@/views/GlossaryView.vue') },
    { path: '/glossary/:slug', name: 'term', component: () => import('@/views/TermDetailView.vue'), props: true },
    { path: '/graph', name: 'graph', component: () => import('@/views/GraphView.vue') },
    { path: '/search', name: 'search', component: () => import('@/views/SearchView.vue') },
    { path: '/about', name: 'about', component: () => import('@/views/AboutView.vue') },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('@/views/NotFoundView.vue') },
  ],
})

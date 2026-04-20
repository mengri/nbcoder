import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import DefaultLayout from '@/layouts/DefaultLayout.vue'
import EmptyLayout from '@/layouts/EmptyLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    component: EmptyLayout,
    children: [
      {
        path: '',
        name: 'Login',
        component: () => import('@/views/Login/index.vue')
      }
    ]
  },
  {
    path: '/',
    component: DefaultLayout,
    redirect: '/projects',
    children: [
      {
        path: 'projects',
        name: 'Projects',
        component: () => import('@/views/Project/index.vue')
      },
      {
        path: 'projects/:id/cards',
        name: 'Cards',
        component: () => import('@/views/Card/Pool.vue')
      },
      {
        path: 'projects/:id/cards/:cardId',
        name: 'CardDetail',
        component: () => import('@/views/Card/Detail.vue')
      },
      {
        path: 'projects/:id/pipelines',
        name: 'Pipelines',
        component: () => import('@/views/Pipeline/List.vue')
      },
      {
        path: 'pipelines/:id',
        name: 'PipelineDetail',
        component: () => import('@/views/Pipeline/Detail.vue')
      },
      {
        path: 'projects/:id/agents',
        name: 'Agents',
        component: () => import('@/views/Agent/Tasks.vue')
      },
      {
        path: 'projects/:id/monitor',
        name: 'AgentMonitor',
        component: () => import('@/views/Agent/Monitor.vue')
      },
      {
        path: 'projects/:id/knowledge',
        name: 'Knowledge',
        component: () => import('@/views/Knowledge/Documents.vue')
      },
      {
        path: 'ai',
        name: 'AI',
        component: () => import('@/views/AI/Providers.vue')
      },
      {
        path: 'ai/statistics',
        name: 'AIStatistics',
        component: () => import('@/views/AI/Statistics.vue')
      }
    ]
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router

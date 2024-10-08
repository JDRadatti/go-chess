import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView
        },
        {
            path: '/game/:id',
            name: 'game',
            component: () => import('../views/GameView.vue')
        },
        {
            path: '/play/',
            name: 'play',
            component: () => import('../views/PlayView.vue')
        },
        {
            path: '/:pathMatch(.*)*',
            name: 'NotFound',
            component: () => import('../views/NotFound.vue')
        },
    ]
})

export default router

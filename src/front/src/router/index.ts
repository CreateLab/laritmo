import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import CourseView from '../views/CourseView.vue'
import LectureDetailView from '../views/LectureDetailView.vue'
import LabDetailView from '../views/LabDetailView.vue'
import AuthView from '../views/AuthView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'home',
            component: HomeView,
        },
        {
            path: '/auth',
            name: 'auth',
            component: AuthView,
        },
        {
            path: '/courses/:id',
            name: 'course',
            component: CourseView,
        },
        {
            path: '/courses/:courseId/lectures/:id',
            name: 'lecture-detail',
            component: LectureDetailView,
        },
        {
            path: '/courses/:courseId/labs/:id',
            name: 'lab-detail',
            component: LabDetailView,
        },
    ],
})

export default router
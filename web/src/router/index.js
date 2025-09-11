import { createRouter, createWebHashHistory } from 'vue-router';

const routes = [
    {
        path: '/',
        component: () => import('../pages/Main.vue'),
        meta: {
            titleKey: 'message.mainpage'
        }
    },
    {
        path: '/rank',
        component: () => import('../pages/Rank.vue'),
        meta: {
            titleKey: 'message.rank'
        }
    },
    {
        path: '/problemset',
        component: () => import('../pages/Problem/Main.vue'),
        meta: {
            titleKey: 'message.problemset'
        }
    },
    {
        path: '/login',
        component: () => import('../pages/Account/Login.vue'),
        meta: {
            titleKey: 'message.login'
        }
    },
    {
        path: '/register',
        component: () => import('../pages/Account/Register.vue'),
        meta: {
            titleKey: 'message.register'
        }
    },
    {
        path: '/profile/:Username',
        component: () => import('../pages/Profile.vue'),
        meta: {
            titleKey: 'message.profile'
        },
        beforeEnter: (to, from, next) => {
            next();
        }
    },
    {
        path: '/problemset/:Pid',
        component: () => import('../pages/Problem/Details.vue'),
        meta: {
            titleKey: 'message.problem'
        },
        beforeEnter: (to, from, next) => {
            next();
        }
    },
    {
        path: '/reset', component: () => import('../pages/Account/Reset.vue'),
        meta: {
            titleKey: 'message.resetpwd'
        }
    },
    {
        path: '/status', component: () => import('../pages/Status.vue'),
        meta: {
            titleKey: 'message.status'
        }
    },
    {
        path: '/competitions', component: () => import('../pages/Competition/Main.vue'),
        meta: {
            titleKey: 'message.competition'
        }
    },
    {
        path: '/competitions/:cid',
        component: () => import('../pages/Competition/Details.vue'),
        meta: {
            titleKey: 'message.competition'
        },
        beforeEnter: (to, from, next) => {
            next();
        }
    },
    {
        path: '/discussion', component: () => import('../pages/Discuss/Main.vue'),
        meta: {
            titleKey: 'message.discussion'
        }
    },
    {
        path: '/discussion/create', component: () => import('../pages/Discuss/New.vue'),
        meta: {
            titleKey: 'message.createDiscussion'
        }
    },
    {
        path: '/settings', 
        component: () => import('../pages/Settings.vue'),
        meta: {
            titleKey: 'message.settings'
        }
    },
    {
        path: '/discussion/:Did',
        component: () => import('../pages/Discuss/Details.vue'),
        meta: {
            titleKey: 'message.discussion'
        },
    },
    {
        path: '/admin', component: () => import('../pages/Admin/Main.vue'),
        meta: {
            titleKey: 'message.admin'
        }
    },
    {
        path: '/admin/problem', component: () => import('../pages/Admin/Problemset.vue'),
        meta: {
            titleKey: 'message.problemmanagement'
        }
    },
    {
        path: '/admin/competition', component: () => import('../pages/Admin/Competitions.vue'),
        meta: {
            titleKey: 'message.competitionmanagement'
        }
    },
    {
        path: '/admin/user', component: () => import('../pages/Admin/Account.vue'),
        meta: {
            titleKey: 'message.usermanagement'
        }
    },
    {
        path: '/admin/ipstat', component: () => import('../pages/Admin/Ipstat.vue'),
        meta: {
            titleKey: 'message.ipstat'
        }
    },
    {
        path: '/403', component: () => import('../pages/403.vue'),
        meta: {
            titleKey: 'message.nopermission'
        }
    }
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});

export default router;

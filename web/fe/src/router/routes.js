const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [{ path: '', component: () => import('pages/home/Home.Page.vue') }],
  },
  {
    path: '/auth/login-web2',
    name: 'Login',
    component: () => import('src/components/auth/Login.Web2.Component.vue'),
    meta: { requiresAuth: false, guest: true },
  },
  {
    path: '/auth/connect2wallet',
    name: 'Conn2Wall',
    component: () => import('src/components/auth/Connect.Wallet.Component.vue'),
    meta: { requiresAuth: false, guest: true },
  },

  // Always leave this as last one,
  // but you can also remove it
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes

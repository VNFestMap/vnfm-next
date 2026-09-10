import tailwindcss from '@tailwindcss/vite'

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',

  devtools: { enabled: false },

  extends: ['@kungal/ui-nuxt'],

  css: ['~/assets/css/main.css'],

  app: {
    pageTransition: { name: 'kun-page', mode: 'out-in' },
    layoutTransition: { name: 'kun-page', mode: 'out-in' },
    head: {
      htmlAttrs: { lang: 'zh-Hans' },
      meta: [
        {
          name: 'theme-color',
          media: '(prefers-color-scheme: light)',
          content: '#f4f4f7'
        },
        {
          name: 'theme-color',
          media: '(prefers-color-scheme: dark)',
          content: '#0a0a0a'
        }
      ]
    }
  },

  modules: [
    '@nuxt/eslint',
    '@nuxtjs/color-mode',
    '@pinia/nuxt',
    'pinia-plugin-persistedstate/nuxt',
    'nuxt-schema-org'
  ],

  devServer: {
    host: '127.0.0.1',
    port: 3710
  },

  pinia: {
    storesDirs: ['./store/**']
  },

  colorMode: {
    preference: 'system',
    fallback: 'light',
    globalName: '__VNFM_COLOR_MODE__',
    componentName: 'ColorScheme',
    classPrefix: 'kun-',
    classSuffix: '-mode',
    storageKey: 'vnfm-color-mode'
  },

  vite: {
    plugins: [tailwindcss()]
  },

  imports: {
    dirs: ['./composables', './config'],
    presets: [
      {
        from: '@kungal/ui-core',
        imports: ['cn']
      }
    ]
  },

  runtimeConfig: {
    public: {
      apiBase: process.env.NUXT_PUBLIC_API_BASE || 'http://127.0.0.1:3711/api/v1',
      oidcIssuer:
        process.env.NUXT_PUBLIC_OIDC_ISSUER || 'http://127.0.0.1:9277',
      oidcServerUrl:
        process.env.NUXT_PUBLIC_OIDC_SERVER_URL ||
        'http://127.0.0.1:9277/api/v1',
      oidcFrontendUrl:
        process.env.NUXT_PUBLIC_OIDC_FRONTEND_URL || 'http://127.0.0.1:9420',
      oidcClientId: process.env.NUXT_PUBLIC_OIDC_CLIENT_ID || 'vnfm-dev',
      oidcRedirectUri:
        process.env.NUXT_PUBLIC_OIDC_REDIRECT_URI ||
        'http://127.0.0.1:3710/auth/callback',
      siteUrl: process.env.NUXT_PUBLIC_SITE_URL || 'http://127.0.0.1:3710',
      r2PublicBase: process.env.NUXT_PUBLIC_R2_PUBLIC_BASE || ''
    }
  }
})

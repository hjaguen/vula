import { defineConfig } from 'vitepress'

export default defineConfig({
  title: 'Vula Manual',
  description: 'Next-Gen AI & Voice Developer OS Engine for Ubuntu 24.04 LTS',
  base: '/vula/',
  head: [
    ['link', { rel: 'icon', href: '/vula/favicon.ico' }],
    ['meta', { name: 'theme-color', content: '#7AA2F7' }]
  ],

  themeConfig: {
    logo: '/logo.svg',
    siteTitle: 'Vula Manual',

    nav: [
      { text: 'Inicio', link: '/' },
      { text: 'Manual', link: '/guide/getting-started' },
      { text: 'CLI Reference', link: '/reference/cli-reference' },
      { text: 'GitHub', link: 'https://github.com/hjaguen/vula' }
    ],

    sidebar: [
      {
        text: '📖 El Manual de Vula',
        items: [
          { text: 'Primeros Pasos e Instalación', link: '/guide/getting-started' },
          { text: 'HUD Overlay Interactivo (Super + Space)', link: '/guide/hud-overlay' },
          { text: 'Control de Sistema por IA (vula do)', link: '/guide/vula-do-ai' },
          { text: 'Asistente y Dictado por Voz', link: '/guide/voice-assistant' }
        ]
      },
      {
        text: '🎨 Personalización & Entorno',
        items: [
          { text: '19 Temas & Fondos HD', link: '/guide/themes-and-wallpapers' },
          { text: 'Tiling Híbrido & Atajos', link: '/guide/tiling-and-shortcuts' },
          { text: 'Herramientas Dev (Fuentes, WebApps & Medios)', link: '/guide/developer-tools' }
        ]
      },
      {
        text: '📋 Referencia Completa',
        items: [
          { text: 'Comandos CLI (Cheat Sheet)', link: '/reference/cli-reference' }
        ]
      }
    ],

    socialLinks: [
      { icon: 'github', link: 'https://github.com/hjaguen/vula' }
    ],

    footer: {
      message: 'Licencia MIT | Vula AI Developer OS',
      copyright: 'Copyright © 2026 Mauricio & Vula OS Team'
    },

    search: {
      provider: 'local'
    }
  }
})

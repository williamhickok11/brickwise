import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import type { App } from 'vue'

const BrickwisePreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#e8f4fc',
      100: '#c5e3f6',
      200: '#9ecfef',
      300: '#6eb8e5',
      400: '#3d9ed9',
      500: '#1a6fb5',
      600: '#155d9a',
      700: '#114b7d',
      800: '#0d3a61',
      900: '#092a47',
      950: '#061c30',
    },
  },
})

export function setupPrimeVue(app: App) {
  app.use(PrimeVue, {
    theme: {
      preset: BrickwisePreset,
      options: {
        darkModeSelector: false,
        cssLayer: false,
      },
    },
  })
}

import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import type { App } from 'vue'

const BrickwisePreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#f1f9f1',
      100: '#dcefdc',
      200: '#bce0bc',
      300: '#93cc93',
      400: '#7ec87e',
      500: '#81c784',
      600: '#66bb6a',
      700: '#4c9a51',
      800: '#3a7840',
      900: '#2a572f',
      950: '#1a3820',
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

import PrimeVue from 'primevue/config'
import Aura from '@primevue/themes/aura'
import { definePreset } from '@primevue/themes'
import type { App } from 'vue'

/** Deep teal primary — strong accent on warm neutrals. */
const BrickwisePreset = definePreset(Aura, {
  semantic: {
    primary: {
      50: '#e6f5f4',
      100: '#c0e6e3',
      200: '#96d5d1',
      300: '#66c0ba',
      400: '#3aa39c',
      500: '#0f766e',
      600: '#0d635c',
      700: '#0b524c',
      800: '#08403c',
      900: '#062e2b',
      950: '#041f1d',
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

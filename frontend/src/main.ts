import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import 'primeicons/primeicons.css'
import router from './router'
import App from './App.vue'
import { setupPrimeVue } from './plugins/primevue'
import './assets/main.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
setupPrimeVue(app)
app.use(ToastService)
app.use(ConfirmationService)
app.component('Toast', Toast)
app.component('ConfirmDialog', ConfirmDialog)

app.mount('#app')

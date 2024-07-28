
import '@/assets/saas-icon/iconfont.css'
import saIcon from "@/assets/saas-icon/SaIcon.vue";

export default {
    install(app) {
        app.component('saIcon', saIcon);
    }
}
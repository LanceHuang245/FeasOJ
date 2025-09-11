import vuetify from './vuetify'
import router from '../router/index';
import vue3AceEditor from './vue3-ace-editor';
import AvatarCropper from 'vue-avatar-cropper';
import { i18n } from './vue-i18n';

export function registerPlugins(app) {
  app.use(vuetify)
  app.use(router)
  app.use(i18n)
  app.use(AvatarCropper)
}

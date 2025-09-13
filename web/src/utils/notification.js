import { ref } from 'vue';
import { apiUrl } from './axios';
import { i18n } from '../plugins/vue-i18n';

const events = ref([]);
const snackbar = ref(false);
const snackbarMessage = ref('');

export function showNotification(message) {
  events.value.push(message);
  snackbarMessage.value = message;
  snackbar.value = true;
}

export function closeNotification() {
  snackbar.value = false;
}

export function useNotificationState() {
  return {
    events,
    snackbar,
    snackbarMessage,
    closeNotification,
  };
}

//
// —— SSE 逻辑 ———————————————————————————————————————————————
//
let es = null;

export function initSSE(id) {
  function connect() {
    const language = i18n.global.locale;
    es = new EventSource(`${apiUrl}/notification/${id}?lang=${language.value}`);

    es.onmessage = (event) => {
      showNotification(event.data);
    };

    es.onerror = () => {
      console.error(i18n.global.t('message.sseError'));
      es.close();
      // 3 秒后重连
      setTimeout(connect, 3000);
    };
  }

  connect();
}

export function closeSSE() {
  if (es) {
    es.close();
    es = null;
  }
}

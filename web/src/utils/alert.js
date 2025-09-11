import { ref } from 'vue';

export const dialog = ref(false);
export const dialogMessage = ref('');

export const showAlert = (message, path) => {
  dialogMessage.value = message;
  dialog.value = true;
  if (path === "reload") {
    setTimeout(() => {
      location.reload();
      dialog.value = false;
    }, 1500);
  } else if (path === "") {

  } else if (path === "/") {
    setTimeout(() => {
      window.location = '#' + path;
      location.reload();
      dialog.value = false;
    }, 1000);
  } else {
    setTimeout(() => {
      window.location = '#' + path;
      location.reload();
      dialog.value = false;
    }, 1000);
  }
};

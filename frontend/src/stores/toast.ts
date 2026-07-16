import { writable } from 'svelte/store';

export type ToastType = 'success' | 'error' | 'info';

export interface ToastMessage {
  id: string;
  type: ToastType;
  message: string;
  duration?: number;
}

function createToastStore() {
  const { subscribe, update } = writable<ToastMessage[]>([]);

  function show(message: string, type: ToastType = 'info', duration: number = 3000) {
    const id = Math.random().toString(36).substring(2, 9);
    const toast = { id, message, type, duration };
    
    update(toasts => [...toasts, toast]);

    if (duration > 0) {
      setTimeout(() => {
        dismiss(id);
      }, duration);
    }
  }

  function dismiss(id: string) {
    update(toasts => toasts.filter(t => t.id !== id));
  }

  return {
    subscribe,
    show,
    dismiss,
    success: (msg: string, duration?: number) => show(msg, 'success', duration),
    error: (msg: string, duration: number = 5000) => show(msg, 'error', duration),
    info: (msg: string, duration?: number) => show(msg, 'info', duration),
  };
}

export const toast = createToastStore();

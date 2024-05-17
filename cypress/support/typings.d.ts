export {};

declare global {
  namespace Cypress {
    interface cy {
      state<T = unknown>(state: string): T;
    }
  }
}

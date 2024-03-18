export {};

declare global {
  namespace Cypress {
    interface cy {
      state(state: string): unknown;
    }
  }
}

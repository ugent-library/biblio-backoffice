export {};

declare global {
  namespace Cypress {
    type Alias = `@${string}`;

    interface cy {
      state<T = unknown>(state: string): T;
    }
  }
}

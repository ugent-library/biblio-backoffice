export {};

declare global {
  namespace Cypress {
    type Alias = `@${string}`;

    type State = {
      ctx: Mocha.Context & Ctx;
      current: Command;
    };

    type Ctx = {
      CSRFToken?: string;
    };

    interface cy {
      state<S extends keyof State>(key: S): State[S];
    }
  }
}

export {};

declare global {
  namespace Cypress {
    type Alias = `@${string}`;

    type State = {
      aliases: Record<string, unknown>;
      ctx: Mocha.Context & Ctx;
      current: Command;
    };

    type Ctx = {
      CSRFToken?: string;
    };

    interface cy {
      state<S extends keyof State, V extends State[S]>(
        key: S,
        value: V,
      ): State[S];
      state<S extends keyof State>(key: S): State[S];
    }
  }
}

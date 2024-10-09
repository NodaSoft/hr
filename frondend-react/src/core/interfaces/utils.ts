export type ModuleSelector<M, T extends unknown[], R> = (
  moduleState: M,
  ...args: T
) => R

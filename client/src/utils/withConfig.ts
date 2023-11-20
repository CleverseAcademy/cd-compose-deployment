import { IBaseRequestConfig } from "../entities/base.request";

type RequestFunc<P, T> = (p: P) => Promise<T>;

const withBaseConfig: <B extends IBaseRequestConfig, P extends B, T>(
  f: RequestFunc<P, T>
) => (baseConfig: B) => (p: Omit<P, keyof B>) => ReturnType<typeof f> =
  <B, P, T>(f) =>
  (baseConfig) =>
  (p) =>
    f({ ...p, ...baseConfig });

export default withBaseConfig;

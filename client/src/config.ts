const {
  CD_CLI_PRIVATE_KEY_PEM: ENV_CLI_PRIVATE_KEY_PEM,
  CD_CLI_MAX_RETRY: ENV_CLI_MAX_RETRY,
} = process.env;

if (!ENV_CLI_PRIVATE_KEY_PEM)
  throw new Error("CD_CLI_PRIVATE_KEY_PEM must be provided");

export const PRIVATE_KEY_PEM = ENV_CLI_PRIVATE_KEY_PEM;
export const MAX_RETRY = Number(ENV_CLI_MAX_RETRY) || 3;

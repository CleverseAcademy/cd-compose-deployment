import { AxiosError } from "axios";
import deploy from "../api/deploy";
import getJti from "../api/getJTI";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeployment } from "../entities/deployment.model";

const deployCommand = async (
  {
    host,
    port,
    service,
    image,
    priority,
    ref,
  }: IDeployment & Omit<IBaseRequestConfig, "privateKey">,
  retry: number = 0
) => {
  if (!process.env.CD_CLI_PRIVATE_KEY_PEM)
    throw new Error("CD_CLI_PRIVATE_KEY_PEM must be provided");
  if (retry > Number(process.env.CD_CLI_MAX_RETRY || 3))
    throw new Error("Max retry exceed");

  const baseConfig = {
    host,
    port,
    privateKey: process.env.CD_CLI_PRIVATE_KEY_PEM!,
  };
  const configuredJtiRequest = getJti(baseConfig);
  const configuredDeployRequest = deploy(baseConfig);

  try {
    return await configuredJtiRequest({}).then(async (jti) => {
      return configuredDeployRequest({
        jti,
        priority,
        ref,
        image,
        service,
      });
    });
  } catch (error) {
    if (error instanceof AxiosError) {
      switch (error.response?.status) {
        case 401:
        case 403:
          console.error(
            `Authorization failed:\n\tmessage: ${error.response.data}\nPlease check environment variable "CD_CLI_PRIVATE_KEY_PEM"`
          );
          break;

        case 424:
          console.error(
            `Expired JTI\n\tmessage: ${error.response.data}, retrying...`
          );
          return deployCommand(
            {
              host,
              image,
              port,
              priority,
              ref,
              service,
            },
            retry + 1
          );

        case 500:
          console.error(`Server error: ${error.response.data}`);
          break;

        default:
          console.error(error);
          break;
      }
    }

    throw error;
  }
};

export default deployCommand;

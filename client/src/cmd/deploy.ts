import { AxiosError } from "axios";
import deploy from "../api/deploy";
import getJti from "../api/getJTI";
import { MAX_RETRY, PRIVATE_KEY_PEM } from "../config";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeploymentPayload } from "../entities/deployment.model";

const deployCommand = async (
  {
    host,
    port,
    service,
    image,
    priority,
    ref,
  }: IDeploymentPayload & Omit<IBaseRequestConfig, "privateKey">,
  retry: number = 0
) => {
  if (retry > MAX_RETRY) throw new Error("Max retry exceed");

  const baseConfig = {
    host,
    port,
    privateKey: PRIVATE_KEY_PEM,
  };
  const configuredJtiRequest = getJti(baseConfig);
  const configuredDeployRequest = deploy(baseConfig);

  try {
    return await configuredJtiRequest({}).then((jti) =>
      configuredDeployRequest({
        jti,
        priority,
        ref,
        image,
        service,
      })
    );
  } catch (error) {
    if (error instanceof AxiosError) {
      switch (error.response?.status) {
        case 401:
        case 403:
          throw new Error(
            `Authorization failed:\n\tmessage: ${error.response.data}\nPlease check environment variable "CD_CLI_PRIVATE_KEY_PEM"`
          );

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
          throw new Error(`Server error: ${error.response.data}`);
      }

      if (error.code === "ECONNREFUSED")
        throw new Error(
          "Can't connect to the server, does the service properly configured?"
        );
    }

    throw error;
  }
};

export default deployCommand;

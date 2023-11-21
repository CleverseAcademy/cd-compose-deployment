import { AxiosError } from "axios";
import getJti from "../api/getJTI";
import nextDeployment from "../api/nexDeployment";
import { MAX_RETRY, PRIVATE_KEY_PEM } from "../config";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeploymentRequest } from "../entities/deployment.request";

const getNextDeploymentCommand = async (
  {
    host,
    port,
    service,
  }: Omit<IDeploymentRequest, "jti"> & Omit<IBaseRequestConfig, "privateKey">,
  retry: number = 0
) => {
  if (retry > MAX_RETRY) throw new Error("Max retry exceed");

  const baseConfig = {
    host,
    port,
    privateKey: PRIVATE_KEY_PEM,
  };
  const configuredJtiRequest = getJti(baseConfig);
  const configuredNextDeploymentRequest = nextDeployment(baseConfig);

  try {
    return await configuredJtiRequest({}).then((jti) =>
      configuredNextDeploymentRequest({
        jti,
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
          return getNextDeploymentCommand(
            {
              host,
              port,
              service,
            },
            retry + 1
          );

        case 500:
          throw new Error(`Server error: ${error.response.data}`);
      }

      if (error.code === "ECONNREFUSED")
        throw new Error(
          "Can't connect to the server, does the service is properly configured?"
        );
    }

    throw error;
  }
};

export default getNextDeploymentCommand;

import axios from "axios";
import DeploymentDto from "../dto/deployment.dto";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeploymentRequest } from "../entities/deployment.request";
import { getRequestSignature } from "../utils/getRequestSignature";
import withBaseConfig from "../utils/withConfig";

type IGetNextDeploymentArgs = IBaseRequestConfig & IDeploymentRequest;

const getNextDeployment = ({
  host,
  port,
  jti,
  privateKey,
  service,
}: IGetNextDeploymentArgs) => {
  return axios
    .request<DeploymentDto>({
      method: "GET",
      maxBodyLength: Infinity,
      url: `http://${host}:${port}/deploy/latest/${service}`,
      headers: {
        Authorization: getRequestSignature({
          jti,
          privateKey,
        }),
      },
    })
    .then(({ status, data }) => {
      switch (status) {
        case 204:
          return "No deployment found";

        case 200:
        default:
          return data;
      }
    });
};

export default withBaseConfig(getNextDeployment);

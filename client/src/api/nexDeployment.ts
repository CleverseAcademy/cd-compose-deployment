import axios from "axios";
import DeploymentDto from "../dto/deployment.dto";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeploymentRequest } from "../entities/deployment.request";
import withBaseConfig from "../utils/withConfig";
import { getRequestSignature } from "./getRequestSignature";

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
        case 200:
          return data;

        case 204:
          return "No deployment found";

        default:
          return data;
      }
    });
};

export default withBaseConfig(getNextDeployment);

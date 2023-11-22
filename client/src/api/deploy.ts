import axios from "axios";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeploymentPayload } from "../entities/deployment.model";
import { IDeploymentRequest } from "../entities/deployment.request";
import { getRequestSignature } from "../utils/getRequestSignature";
import withBaseConfig from "../utils/withConfig";

type IDeployArgs = IDeploymentPayload & IBaseRequestConfig & IDeploymentRequest;

const deploy = ({
  host,
  port,
  priority,
  ref,
  service,
  image,
  jti,
  privateKey,
}: IDeployArgs) => {
  const data = JSON.stringify({
    p: priority,
    r: `cli-${ref}`,
    s: service,
    i: image,
  });

  return axios
    .request<string>({
      method: "POST",
      url: `http://${host}:${port}/deploy`,
      headers: {
        "Content-Type": "application/json",
        Authorization: getRequestSignature({ jti, privateKey, data }),
      },
      data: data,
    })
    .then((response) => response.data);
};

export default withBaseConfig(deploy);

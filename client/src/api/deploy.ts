import axios from "axios";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeployment } from "../entities/deployment.model";
import withBaseConfig from "../utils/withConfig";
import { getRequestSignature } from "./getRequestSignature";

type IDeployArgs = IDeployment &
  IBaseRequestConfig & {
    jti: string;
  };

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
    r: `cli-${ref} ${new Date().toISOString()}`,
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

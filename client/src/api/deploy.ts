import axios from "axios";
import { sign } from "jsonwebtoken";
import { createHash } from "node:crypto";
import { IBaseRequestConfig } from "../entities/base.request";
import { IDeployment } from "../entities/deployment.model";
import withBaseConfig from "../utils/withConfig";

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
        Authorization: sign(
          {
            cs: createHash("sha256").update(data).digest("hex"),
          },
          privateKey,
          {
            algorithm: "ES256",
            jwtid: jti,
            notBefore: "0s",
            expiresIn: "55s",
          }
        ),
      },
      data: data,
    })
    .then((response) => response.data);
};

export default withBaseConfig(deploy);

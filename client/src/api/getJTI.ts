import axios from "axios";
import { IBaseRequestConfig } from "../entities/base.request";
import withBaseConfig from "../utils/withConfig";

const getJTI = ({ host, port }: IBaseRequestConfig) =>
  axios
    .request<string>({
      method: "GET",
      url: `http://${host}:${port}/nextJTI`,
    })
    .then((response) => response.data);

export default withBaseConfig(getJTI);

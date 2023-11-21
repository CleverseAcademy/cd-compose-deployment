import { sign } from "jsonwebtoken";
import { createHash } from "node:crypto";

interface IArgsGetRequestSignature {
  jti: string;
  privateKey: string;
  data?: string;
}

export function getRequestSignature({
  jti,
  privateKey,
  data,
}: IArgsGetRequestSignature): string {
  return sign(
    data
      ? {
          cs: createHash("sha256").update(data).digest("hex"),
        }
      : {},
    privateKey,
    {
      algorithm: "ES256",
      jwtid: jti,
      notBefore: "0s",
      expiresIn: "55s",
    }
  );
}

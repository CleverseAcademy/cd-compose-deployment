import { program } from "commander";
import deployCommand from "./cmd/deploy";

program
  .command("deploy")
  .description("Deploy a new docker image to an existing compose project")
  .argument(
    "<service name>",
    "Target service corresponding with service item key specified in compose.y(a)ml"
  )
  .requiredOption(
    "-i, --image <docker image name>",
    "docker image to deploy to the target"
  )
  .requiredOption(
    "-t, --target <server ip>",
    "Target host which running docker image: cloudiana/compose-deployment"
  )
  .option(
    "-p, --priority <number>",
    "Deployment priority, higher number get a higher chance to be deployed"
  )
  .option(
    "-r, --ref <string>",
    "Deployment priority, higher number get a higher chance to be deployed"
  )
  .option(
    "--git",
    "Determine deployment priority from git commit chronologically"
  )
  .action((service, { target: host, image, priority, ref }) => {
    deployCommand({
      host,
      port: 3000,
      priority: Number(priority),
      ref,
      image,
      service,
    }).then((v) => console.log(`Deployment accepted: #queue = ${v}`));
  });

program.parse();

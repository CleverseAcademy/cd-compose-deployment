import { program } from "commander";
import deployCommand from "./cmd/deploy";
import getNextDeploymentCommand from "./cmd/nextDeployment";

program
  .command("deploy")
  .description("Deploy a new docker image to a target server")
  .option(
    "--git",
    "Determine deployment priority from git commit chronologically"
  )
  .option(
    "-p, --priority <number>",
    "Deployment priority, higher number get a higher chance to be deployed"
  )
  .option(
    "-r, --ref <string>",
    "Deployment priority, higher number get a higher chance to be deployed"
  )
  .requiredOption(
    "-i, --image <docker image name>",
    "docker image to deploy to the target"
  )
  .requiredOption(
    "-t, --target <server ip>",
    "Target host which running docker image: cloudiana/compose-deployment"
  )
  .option("--port <number>", "service port", "3000")
  .argument(
    "<service name>",
    "Target service corresponding with service item key specified in compose.y(a)ml"
  )
  .action((service, { target: host, image, priority, ref, port }) => {
    deployCommand({
      host,
      port: port || 3000,
      priority: Number(priority),
      ref,
      image,
      service,
    }).then((v) => console.log(`Deployment accepted: #queue = ${v}`));
  });

program
  .command("next-deployment")
  .description("Get next deployment info from server's current queue")
  .requiredOption(
    "-t, --target <server ip>",
    "Target host which running docker image: cloudiana/compose-deployment"
  )
  .option("--port <number>", "service port", "3000")
  .argument(
    "<service name>",
    "Target service corresponding with service item key specified in compose.y(a)ml"
  )
  .action((service, { target: host, port }) => {
    getNextDeploymentCommand({
      host,
      port: port || 3000,
      service,
    }).then((v) => console.table(v));
  });

program.parse();

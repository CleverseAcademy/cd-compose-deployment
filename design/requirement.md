# User Stories
1. As a team leader, I will only deploy a stable version of reviewed code to a production, so that I will know where to start an investigation when issues happened at production website.
2. As a whole team, we want to deploy the latest version of pre-production code from the `staging` branch, so that we can ensure the quality of the production website before a production deployment.
3. As a whole team, we want to manage the version of code from `staging` branch deployed to the server, so that we can review overall progress from time-to-time.

# Functional Requirements
1. The system shall update the image version specified in `compose.yml` immediately, then create a container with the same configuration as prior, once received an authorized API call.
2. The system shall prohibit an unauthorized request from updating the deployment version.
3. The system shall provide a sufficient user authentication system, to allow user-level API access.
4. The system shall provide a deployment log to an authorized user, for further issue investigation.
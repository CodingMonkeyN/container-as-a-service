export interface ContainerDeployment {
  namespace: string;
  name: string;
  image: string;
  cpu: string;
  memory: string;
  replicas: number;
  environmentVariables: Record<string, string>;
}

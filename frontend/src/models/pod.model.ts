export interface Pod {
  namespace: string;
  name: string;
  cpu: string;
  memory: string;
  ready: boolean
  status: string;
}
